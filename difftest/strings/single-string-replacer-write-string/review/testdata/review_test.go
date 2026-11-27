package strings_test

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

var stopped = errors.New("review writer stopped")

type response struct {
	count int
	err   error
}

type writeRecorder struct {
	calls       []string
	script      []response
	omitEmpty   bool
	rejectShort bool
}

func (w *writeRecorder) record(method, payload string) (int, error) {
	if w.omitEmpty && payload == "" {
		return 0, nil
	}
	i := len(w.calls)
	w.calls = append(w.calls, method+":"+payload)
	if i < len(w.script) {
		reply := w.script[i]
		if w.rejectShort && reply.err == nil && reply.count < len(payload) {
			return reply.count, io.ErrShortWrite
		}
		return reply.count, reply.err
	}
	return len(payload), nil
}

func (w *writeRecorder) Write(p []byte) (int, error) {
	return w.record("Write", string(p))
}

type recordingStringWriter struct{ *writeRecorder }

func (w recordingStringWriter) WriteString(p string) (int, error) {
	return w.record("WriteString", p)
}

func TestGoLeanSingleStringReplacerWriteStringReview(t *testing.T) {
	fixtures := []struct {
		pattern, value, input string
		direct                bool
		script                []response
		n                     int
		err                   error
		calls                 []string
	}{
		{"aa", "Z", "aaXaa", true, nil, 3, nil, []string{"WriteString:", "WriteString:Z", "WriteString:X", "WriteString:Z", "WriteString:"}},
		{"aa", "Z", "aaXaa", false, []response{{0, stopped}}, 0, stopped, []string{"Write:"}},
		{"aa", "ZZ", "XaaYaa", true, []response{{1, nil}, {1, stopped}}, 2, stopped, []string{"WriteString:X", "WriteString:ZZ"}},
		{"aa", "ZZ", "XaaYaa", false, []response{{1, nil}, {1, nil}}, 5, nil, []string{"Write:X", "Write:ZZ", "Write:Y", "Write:ZZ", "Write:"}},
		{"z", "A", "\xffX", false, []response{{1, stopped}}, 1, stopped, []string{"Write:\xffX"}},
		{"z", "", "z", true, nil, 0, nil, []string{"WriteString:", "WriteString:", "WriteString:"}},
		{"z", "A", "", false, nil, 0, nil, []string{"Write:"}},
	}
	for _, c := range []struct {
		name                   string
		omitEmpty, rejectShort bool
		want                   int
	}{
		{"Go", false, false, 0},
		{"Suppress empty writes", true, false, 5},
		{"Convert a short nil-error reply to an error", false, true, 1},
	} {
		bad := 0
		for _, f := range fixtures {
			rec := &writeRecorder{script: f.script, omitEmpty: c.omitEmpty, rejectShort: c.rejectShort}
			var writer io.Writer = rec
			if f.direct {
				writer = recordingStringWriter{rec}
			}
			n, err := strings.GoLeanSingleStringReplacerWriteString(f.pattern, f.value, f.input, writer)
			if n != f.n || err != f.err || !reflect.DeepEqual(rec.calls, f.calls) {
				bad++
				if c.want == 0 {
					t.Errorf("WriteString(%x,%x,%x) gave (%d,%v,%q), want (%d,%v,%q)",
						f.pattern, f.value, f.input, n, err, rec.calls, f.n, f.err, f.calls)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
