package zip

import (
	"encoding/json"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func reportConfig(w *Writer, action func() error) (output string) {
	status := "ok"
	defer func() {
		if p := recover(); p != nil {
			if p == "zip: SetOffset called after data was written" {
				status = "afterWrite"
			} else if e, ok := p.(runtime.Error); ok && strings.Contains(e.Error(), "nil pointer") {
				status = "nil"
			} else {
				panic(p)
			}
		}
		if w == nil {
			output = status + ":none"
			return
		}
		count := "none"
		if w.cw != nil {
			count = strconv.FormatInt(w.cw.count, 10)
		}
		raw := []byte(w.comment)
		cells := make([]string, len(raw))
		for i, b := range raw {
			cells[i] = strconv.Itoa(int(b))
		}
		output = status + ":" + count + ":" + strings.Join(cells, ",")
	}()
	if err := action(); err != nil {
		if err.Error() != "zip: Writer.Comment too long" {
			panic(err)
		}
		status = "tooLong"
	}
	return
}

func TestGoLeanZipWriterConfig(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ZIP_WRITER_CONFIG"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := json.NewDecoder(f)
	for row := 0; ; row++ {
		var c struct {
			Tag, Repeat   int
			Count, Offset int64
			Fill          byte
			Old, Tail     []byte
			Want          string
		}
		err := in.Decode(&c)
		if err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		makeWriter := func() *Writer {
			if c.Tag == 0 {
				return nil
			}
			w := &Writer{comment: string(c.Old), last: nil}
			if c.Tag == 2 {
				w.cw = &countWriter{count: c.Count}
			}
			return w
		}
		w := makeWriter()
		comment := strings.Repeat(string([]byte{c.Fill}), c.Repeat) + string(c.Tail)
		first := reportConfig(w, func() error { return w.SetComment(comment) })
		w = makeWriter()
		second := reportConfig(w, func() error { w.SetOffset(c.Offset); return nil })
		if got := first + "|" + second; got != c.Want {
			t.Fatalf("row %d: observed state differed from expectation", row)
		}
	}
}
