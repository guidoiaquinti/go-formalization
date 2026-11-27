package zip

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

type writerIORecorder struct {
	writes [][]byte
	closes int
	reply  int
	err    error
}

func (w *writerIORecorder) Write(p []byte) (int, error) {
	w.writes = append(w.writes, bytes.Clone(p))
	return w.reply, w.err
}
func (w *writerIORecorder) Close() error {
	w.closes++
	return w.err
}

func TestGoLeanZipWriterIO(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ZIP_WRITER_IO_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := gob.NewDecoder(bufio.NewReader(f))
	for row := 0; ; row++ {
		var c struct {
			Count             int64
			Written, Existing int
			Failed            bool
			Data              []byte
			Want              string
		}
		if err := in.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row > 0 {
				t.Logf("%d rows agree for all four pinned writer methods", row)
				return
			}
			t.Fatal(err)
		}
		var failure error
		if c.Failed {
			failure = errors.New("scripted write error")
		}
		sink := &writerIORecorder{reply: c.Written, err: failure}
		cw := &countWriter{w: sink, count: c.Count}
		n, err := cw.Write(c.Data)
		if err != failure || len(sink.writes) != 1 || !bytes.Equal(sink.writes[0], c.Data) {
			t.Fatalf("row %d: countWriter altered callback input or error identity", row)
		}
		counted := fmt.Sprintf("%d:%d:%t:%d", cw.count, n, err != nil, len(sink.writes))

		sink.writes = nil
		buffer := bufio.NewWriter(sink)
		if n, err := buffer.Write(c.Data); n != len(c.Data) || err != nil || len(sink.writes) != 0 {
			t.Fatal("fixture must start with pending data and no callback")
		}
		z := &Writer{cw: &countWriter{w: buffer, count: c.Count}}
		flush := func() string {
			offset := len(sink.writes)
			err := z.Flush()
			if (err != nil && err != failure) || len(sink.writes) > offset+1 {
				t.Fatalf("row %d: flush changed callback error identity or count", row)
			}
			if len(sink.writes) > offset && !bytes.Equal(sink.writes[offset], c.Data) {
				t.Fatalf("row %d: flush changed callback payload", row)
			}
			return fmt.Sprintf("%d:%t:%d:%d", buffer.Buffered(), err != nil, len(sink.writes)-offset, z.cw.count)
		}
		first, second := flush(), flush()

		sink.writes = nil
		var original *bufio.Writer
		var underlying io.Writer
		if c.Existing >= 0 {
			underlying = sink
		}
		if c.Existing > 0 {
			original = bufio.NewWriterSize(sink, c.Existing)
			if err := original.WriteByte(17); err != nil {
				t.Fatal(err)
			}
			underlying = original
		}
		created := NewWriter(underlying)
		actual, ok := created.cw.w.(*bufio.Writer)
		if !ok || created.dir != nil || created.last != nil || created.closed || created.compressors != nil ||
			created.comment != "" || created.testHookCloseSizeOffset != nil {
			t.Fatalf("row %d: constructor state differs", row)
		}
		construction := fmt.Sprintf("%d:%d:%t:%d:%d", actual.Size(), actual.Buffered(), actual == original, created.cw.count, len(sink.writes))
		sink.writes = nil
		closeErr := (nopCloser{Writer: sink}).Close()
		closing := fmt.Sprintf("%t:%d", closeErr != nil, len(sink.writes)+sink.closes)
		if got := fmt.Sprintf("%s|%s|%s|%s|%s", counted, first, second, construction, closing); got != c.Want {
			t.Fatalf("row %d: pinned Go %q, expected %q", row, got, c.Want)
		}
	}
}
