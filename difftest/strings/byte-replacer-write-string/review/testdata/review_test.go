package strings_test

import (
	"errors"
	"io"
	"runtime"
	"slices"
	"strings"
	"testing"
)

var stopByteWriter = errors.New("byte writer stopped")

type byteReply struct {
	count int
	err   error
}

type byteWriteRecorder struct {
	calls        []string
	replies      []byteReply
	ignoreErrors bool
}

func (w *byteWriteRecorder) record(method, data string) (int, error) {
	i := len(w.calls)
	w.calls = append(w.calls, method+":"+data)
	if i < len(w.replies) {
		reply := w.replies[i]
		if w.ignoreErrors {
			reply.err = nil
		}
		return reply.count, reply.err
	}
	return len(data), nil
}

func (w *byteWriteRecorder) Write(data []byte) (int, error) { return w.record("B", string(data)) }

type byteWriteDual struct{ *byteWriteRecorder }

func (w byteWriteDual) WriteString(data string) (int, error) { return w.record("S", data) }

func observeByteWrite(mapping *[256]byte, writer io.Writer, text string) (n int, err error, panicked bool) {
	defer func() {
		if fault := recover(); fault != nil {
			runtimeFault, ok := fault.(runtime.Error)
			if !ok || runtimeFault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(fault)
			}
			panicked = true
		}
	}()
	n, err = strings.GoLeanByteReplacerWriteString(mapping, writer, text)
	return n, err, false
}

func TestGoLeanByteReplacerWriteStringReview(t *testing.T) {
	var mapping [256]byte
	for i := range mapping {
		mapping[i] = byte(i)
	}
	mapping['b'], mapping[0xff] = 0xff, 0
	fixtures := []struct {
		text         string
		dual, nilled bool
		replies      []byteReply
		count        int
		fault        error
		panic        bool
		trace        []string
	}{
		{"abcb!", true, false, nil, 5, nil, false, []string{"S:a", "B:\xff", "S:c", "B:\xff", "S:!"}},
		{"abx", false, false, []byteReply{{0, nil}, {1, stopByteWriter}}, 1, stopByteWriter, false, []string{"B:a", "B:\xff"}},
		{"\xff\xffx", true, false, []byteReply{{0, nil}}, 2, nil, false, []string{"B:\x00", "B:\x00", "S:x"}},
		{"abz", true, false, []byteReply{{1, stopByteWriter}}, 1, stopByteWriter, false, []string{"S:a"}},
		{"bb", true, false, nil, 2, nil, false, []string{"B:\xff", "B:\xff"}},
		{"", true, true, nil, 0, nil, false, nil},
		{"x", true, true, nil, 0, nil, true, nil},
	}
	for _, c := range []struct {
		name                     string
		forceBytes, ignoreErrors bool
		want                     int
	}{
		{"Go", false, false, 0},
		{"Route all spans through Write", true, false, 3},
		{"Continue after a writer error", false, true, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			rec := &byteWriteRecorder{replies: f.replies, ignoreErrors: c.ignoreErrors}
			var writer io.Writer = rec
			if f.dual && !c.forceBytes {
				writer = byteWriteDual{rec}
			}
			table := &mapping
			if f.nilled {
				table = nil
			}
			n, err, panicked := observeByteWrite(table, writer, f.text)
			if n != f.count || err != f.fault || panicked != f.panic || !slices.Equal(rec.calls, f.trace) {
				bad++
				if c.want == 0 {
					t.Errorf("WriteString(%x) got (%d,%v,%t) %q", f.text, n, err, panicked, rec.calls)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
