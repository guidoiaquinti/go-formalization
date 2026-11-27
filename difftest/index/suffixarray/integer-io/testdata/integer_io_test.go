package suffixarray

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

var ioOther = errors.New("suffixarray test I/O failure")

type observedReader struct {
	scratch, data           []byte
	limit, offset, terminal int
	requests                []string
}

func (r *observedReader) Read(p []byte) (int, error) {
	r.requests = append(r.requests, strconv.Itoa(len(p)))
	if &p[0] != &r.scratch[r.offset] || cap(p) != len(r.scratch)-r.offset {
		panic("reader lost its scratch alias")
	}
	n := copy(p[:min(len(p), r.limit)], r.data)
	r.data = r.data[n:]
	r.offset += n
	if len(r.data) > 0 {
		return n, nil
	}
	if r.terminal == 1 {
		return n, ioOther
	}
	return n, io.EOF
}

type observedWriter struct {
	scratch []byte
	offer   string
	other   bool
}

func (w *observedWriter) Write(p []byte) (int, error) {
	if cap(p) != len(w.scratch) || &p[0] != &w.scratch[0] {
		panic("writer lost its scratch alias")
	}
	if w.offer != "none" {
		panic("multiple writes")
	}
	w.offer = fmt.Sprintf("%d,%d:%s", len(p), cap(p), wireBytes(p))
	if w.other {
		return 0, ioOther
	}
	return 0, nil
}

func TestGoLeanSuffixarrayIntegerIO(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SUFFIXARRAY_INTEGER_IO_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows := json.NewDecoder(f)
	for row := 0; ; row++ {
		var c struct {
			Length, Chunk, Terminal, Reply, Profile int
			Value                                   int
			Backing, Source                         []byte
			Integers                                []int64
			Want                                    string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if row == 0 {
				t.Fatal("missing records")
			}
			t.Logf("three pinned suffixarray framing entries matched %d rows", row)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		var data ints
		if c.Profile == 0 {
			data.int32 = make([]int32, len(c.Integers))
			for i, v := range c.Integers {
				data.int32[i] = int32(v)
			}
		} else {
			data.int64 = c.Integers
		}
		write := func(pack bool) string {
			buf := slices.Clone(c.Backing)
			w := &observedWriter{scratch: buf, offer: "none", other: c.Reply == 1}
			var writer io.Writer = w
			if c.Reply == 2 {
				writer = nil
			}
			var count int
			status := observeIO(func() error {
				if !pack {
					return writeInt(writer, buf[:c.Length:len(buf)], c.Value)
				}
				var err error
				count, err = writeSlice(writer, buf[:c.Length:len(buf)], data)
				return err
			})
			return fmt.Sprintf("%s/%d/%s/%s", status, count, wireBytes(buf), w.offer)
		}
		buf := slices.Clone(c.Backing)
		r := &observedReader{scratch: buf, data: c.Source, limit: c.Chunk, terminal: c.Terminal}
		var reader io.Reader = r
		if c.Terminal == 2 {
			reader = nil
		}
		var v int64
		status := observeIO(func() error {
			var err error
			v, err = readInt(reader, buf[:c.Length:len(buf)])
			return err
		})
		got := fmt.Sprintf("%s|%s/%d/%s/%s|%s", write(false), status, v, wireBytes(buf), strings.Join(r.requests, ","), write(true))
		if got != c.Want {
			t.Fatalf("row %d %+v: Go %q want %q", row, c, got, c.Want)
		}
	}
}

func observeIO(call func() error) (status string) {
	defer func() {
		if failure := recover(); failure != nil {
			if err, ok := failure.(runtime.Error); ok {
				switch msg := err.Error(); {
				case msg == "runtime error: invalid memory address or nil pointer dereference":
					status = "nilIO"
					return
				case strings.HasPrefix(msg, "runtime error: slice bounds out of range"),
					strings.HasPrefix(msg, "runtime error: index out of range"):
					status = "bounds"
					return
				}
			}
			panic(failure)
		}
	}()
	switch err := call(); err {
	case nil:
		return "ok"
	case io.EOF:
		return "eof"
	case io.ErrUnexpectedEOF:
		return "unexpected"
	case ioOther:
		return "other"
	default:
		panic(err)
	}
}

func wireBytes(data []byte) string {
	fields := make([]string, len(data))
	for i, b := range data {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}
