package strings_test

import (
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
)

var stopped = errors.New("review writer stopped")

type recordingWriter struct {
	n     int
	err   error
	calls []string
}

func (w *recordingWriter) Write(payload []byte) (int, error) {
	if payload == nil {
		w.calls = append(w.calls, "Write:nil")
	} else {
		w.calls = append(w.calls, "Write:"+string(payload))
	}
	return w.n, w.err
}

type dualWriter struct{ *recordingWriter }

func (w dualWriter) WriteString(payload string) (int, error) {
	w.calls = append(w.calls, "WriteString:"+payload)
	return w.n, w.err
}

func dynamicDispatch(w io.Writer, s string) (int, error) {
	return io.WriteString(w, s)
}

func suppressEmpty(w io.Writer, s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strings.GoLeanReviewStringWriter(w, s)
}

func inventError(w io.Writer, s string) (int, error) {
	n, err := strings.GoLeanReviewStringWriter(w, s)
	if err == nil && n < len(s) {
		return n, io.ErrShortWrite
	}
	return n, err
}

func TestGoLeanStringWriterReview(t *testing.T) {
	fixtures := []struct {
		input string
		dual  bool
		n     int
		err   error
	}{
		{"", false, 0, nil}, {"", true, 0, stopped},
		{"\xff\x00🚀", false, 0, nil}, {"\xff\xfe", true, 1, stopped},
		{"abc", true, 3, nil}, {"xy", false, 1, nil},
		{"z", false, -7, stopped},
	}
	for _, c := range []struct {
		name  string
		write func(io.Writer, string) (int, error)
		want  int
	}{
		{"Go", strings.GoLeanReviewStringWriter, 0},
		{"Re-dispatch to WriteString", dynamicDispatch, 3},
		{"Omit empty writes", suppressEmpty, 2},
		{"Invent short write errors", inventError, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			rec := &recordingWriter{n: f.n, err: f.err}
			var writer io.Writer = rec
			if f.dual {
				writer = dualWriter{rec}
			}
			n, err := c.write(writer, f.input)
			if n != f.n || err != f.err || !slices.Equal(rec.calls, []string{"Write:" + f.input}) {
				bad++
				if c.want == 0 {
					t.Errorf("WriteString(%x) gave (%d,%v,%q)", f.input, n, err, rec.calls)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
