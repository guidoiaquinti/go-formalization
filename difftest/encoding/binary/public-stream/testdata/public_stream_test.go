package binary_test

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type streamFixture struct {
	codecFixture
	Chunk, ReaderEnd, WriterEnd int
}

var streamError = errors.New("scripted I/O error")
var streamWrappedEOF = fmt.Errorf("custom EOF: %w", io.EOF)
var streamPanic = new(int)
var streamReaderErrors = []error{io.EOF, streamError, streamWrappedEOF}

type fixtureReader struct {
	remaining []byte
	chunk     int
	stop      error
	requests  []uint64
}

func (r *fixtureReader) Read(p []byte) (int, error) {
	r.requests = append(r.requests, uint64(len(p)))
	n := copy(p[:min(r.chunk, len(p))], r.remaining)
	r.remaining = r.remaining[n:]
	if len(r.remaining) == 0 {
		return n, r.stop
	}
	return n, nil
}

type fixtureWriter struct {
	original []byte
	output   []byte
	calls    int
	response int
	alias    bool
}

func (w *fixtureWriter) Write(p []byte) (int, error) {
	w.calls++
	w.output = slices.Clone(p)
	w.alias = len(p) > 0 && len(w.original) > 0 && &p[0] == &w.original[0]
	switch w.response {
	case 1:
		return 0, streamError
	case 2:
		panic(streamPanic)
	default:
		// Deliberately report a short successful write. binary.Write ignores n.
		return 0, nil
	}
}

func streamValue(c codecFixture) (any, reflect.Value) {
	if c.Form == 3 {
		return "unsupported", reflect.Value{}
	}
	return c.materialize()
}

func nativeStream(c streamFixture) string {
	var order binary.ByteOrder = binary.LittleEndian
	if c.Order != 0 {
		order = binary.BigEndian
	}
	target, destination := streamValue(c.codecFixture)
	reader := fixtureReader{remaining: c.Backing[:c.Length], chunk: c.Chunk, stop: streamReaderErrors[c.ReaderEnd]}
	read := streamStatus(binary.Read(&reader, order, target))
	value, _ := streamValue(c.codecFixture)
	original, _ := value.([]byte)
	writer := fixtureWriter{original: original, response: c.WriterEnd}
	write := func() (status string) {
		defer func() {
			if failure := recover(); failure != nil {
				if failure != streamPanic {
					panic(failure)
				}
				status = "panic"
			}
		}()
		return streamStatus(binary.Write(&writer, order, value))
	}()
	var bytes []uint64
	for _, b := range writer.output {
		bytes = append(bytes, uint64(b))
	}
	return fmt.Sprintf("%s:%s:%s|%s:%d:%t:%s", read, streamWire(observedWords(c.codecFixture, destination)),
		streamWire(reader.requests), write, writer.calls, writer.alias, streamWire(bytes))
}

func streamStatus(err error) string {
	switch err {
	case nil:
		return "ok"
	case io.EOF:
		return "eof"
	case io.ErrUnexpectedEOF:
		return "unexpected"
	case streamError:
		return "other"
	case streamWrappedEOF:
		return "wrapped"
	default:
		if strings.Contains(err.Error(), "invalid type") || strings.Contains(err.Error(), "not fixed-sized") {
			return "type"
		}
		return err.Error()
	}
}

func streamWire(xs []uint64) string {
	items := make([]string, len(xs))
	for i, v := range xs {
		items[i] = strconv.FormatUint(v, 10)
	}
	return strings.Join(items, ",")
}

func streamOracle(c streamFixture) string {
	initial := streamWire(c.Words)
	if c.Form == 3 {
		return "type:" + initial + ":|type:0:false:"
	}
	namedByte := c.Kind == 11
	byteAlias := c.Kind == 2 && c.Form == 2 && len(c.Words) > 0
	if namedByte {
		c.Kind = 2
	}
	width := len(c.Words) * [...]int{1, 1, 1, 2, 2, 4, 4, 8, 8, 4, 8}[c.Kind]
	var requests []uint64
	read, values := "ok", initial
	if namedByte && c.Form == 0 {
		read = "type"
	} else {
		need := width
		for consumed := 0; need > 0; consumed += c.Chunk {
			requests = append(requests, uint64(need))
			if need <= min(c.Length-consumed, c.Chunk) {
				break
			}
			if c.Length-consumed <= c.Chunk {
				read = []string{"eof", "other", "wrapped"}[c.ReaderEnd]
				if read == "eof" && c.Length > 0 {
					read = "unexpected"
				}
				break
			}
			need -= c.Chunk
		}
		if read == "ok" && c.Form == 0 {
			read = "type"
		} else if read == "ok" {
			input := c.codecFixture
			input.Length = width
			decoded, _, _ := strings.Cut(codecOracle(input), "|")
			_, values, _ = strings.Cut(decoded, ":")
		}
	}
	payload := c.codecFixture
	payload.Length, payload.Backing = width, make([]byte, width)
	_, e, _ := strings.Cut(codecOracle(payload), "|")
	_, output, _ := strings.Cut(e, ":")
	return fmt.Sprintf("%s:%s:%s|%s:1:%t:%s", read, values, streamWire(requests),
		[]string{"ok", "other", "panic"}[c.WriterEnd], byteAlias, output)
}

func TestGoLeanPublicStream(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_BINARY_PUBLIC_STREAM_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewScanner(f)
	count := 0
	for reader.Scan() {
		var c streamFixture
		if err := json.Unmarshal(reader.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		got, want := nativeStream(c), streamOracle(c)
		if got != c.Want || got != want {
			t.Fatalf("fixture %+v: Go %q, oracle %q", c, got, want)
		}
		count++
	}
	if err := reader.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("missing cases")
	}
	t.Logf("%d public stream observations per API matched Lean and independent oracle", count)
}
