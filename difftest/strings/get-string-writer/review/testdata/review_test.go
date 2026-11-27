package strings_test

import (
	"io"
	"strings"
	"testing"
)

type dual struct{ last string }

func (w *dual) Write(p []byte) (int, error) {
	w.last = "Write:" + string(p)
	return len(p), nil
}

func (w *dual) WriteString(s string) (int, error) {
	w.last = "WriteString:" + s
	return len(s), nil
}

type plain struct{ last string }

func (w *plain) Write(p []byte) (int, error) {
	w.last = string(p)
	return len(p), nil
}

type writeAdapter struct{ w io.Writer }

func (w writeAdapter) WriteString(s string) (int, error) { return w.w.Write([]byte(s)) }

func directOnly(w io.Writer) io.StringWriter {
	sw, _ := w.(io.StringWriter)
	return sw
}

func TestGoLeanGetStringWriterReview(t *testing.T) {
	text := "\xff\x00🚀"
	for _, c := range []struct {
		name         string
		selectWriter func(io.Writer) io.StringWriter
		want         int
	}{
		{"Go", strings.GoLeanReviewGetStringWriter, 0},
		{"Always wrap", func(w io.Writer) io.StringWriter { return writeAdapter{w} }, 2},
		{"Return only direct implementation", directOnly, 3},
	} {
		d, p := new(dual), new(plain)
		var nilDual *dual
		var nilPlain *plain
		bad := 0
		for _, f := range []struct {
			name     string
			w        io.Writer
			identity io.StringWriter
			count    int
			err      error
			log      *string
			payload  string
		}{
			{"dual", d, d, len(text), nil, &d.last, "WriteString:" + text},
			{"plain", p, nil, len(text), nil, &p.last, text},
			{"typed nil dual", nilDual, nilDual, 0, nil, nil, ""},
			{"typed nil plain", nilPlain, nil, 0, nil, nil, ""},
			{"nil interface", nil, nil, 0, nil, nil, ""},
		} {
			sw := c.selectWriter(f.w)
			var n int
			var err error
			if sw != nil && f.log != nil {
				n, err = sw.WriteString(text)
			}
			if sw == nil || (f.identity != nil && sw != f.identity) || n != f.count ||
				err != f.err || (f.log != nil && *f.log != f.payload) {
				bad++
				if c.want == 0 {
					t.Errorf("%s: selected %T, got (%d,%v)", f.name, sw, n, err)
				}
			}
		}
		t.Logf("%s disagreed on %d of 5", c.name, bad)
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
