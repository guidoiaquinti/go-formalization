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
	calls   []string
	replies []byteReply
}

func (w *byteWriteRecorder) record(method, data string) (int, error) {
	i := len(w.calls)
	w.calls = append(w.calls, method+":"+data)
	if i < len(w.replies) {
		reply := w.replies[i]
		return reply.count, reply.err
	}
	return len(data), nil
}

func (w *byteWriteRecorder) Write(data []byte) (int, error) {
	if data == nil {
		return w.record("nil", "")
	}
	return w.record("B", string(data))
}

type byteWriteDual struct{ *byteWriteRecorder }

func (w byteWriteDual) WriteString(data string) (int, error) { return w.record("S", data) }

func observeByteWrite(mapping *[256][]byte, writer io.Writer, text string) (n int, err error, panicked bool) {
	defer func() {
		if fault := recover(); fault != nil {
			runtimeFault, ok := fault.(runtime.Error)
			if !ok || runtimeFault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(fault)
			}
			panicked = true
		}
	}()
	n, err = strings.GoLeanByteStringReplacerWriteString(mapping, writer, text)
	return n, err, false
}

func TestGoLeanByteStringReplacerWriteStringReview(t *testing.T) {
	var mapping [256][]byte
	mapping['b'], mapping['i'], mapping['x'] = []byte{}, []byte("i"), []byte{0xff, 0}
	fixtures := []struct {
		text         string
		dual, nilled bool
		replies      []byteReply
		count        int
		fault        error
		panic        bool
		trace        []string
	}{
		{"zbin", true, false, nil, 3, nil, false, []string{"S:z", "B:", "B:i", "S:n"}},
		{"bb", true, false, nil, 0, nil, false, []string{"B:", "B:"}},
		{"bX", true, false, []byteReply{{0, stopByteWriter}}, 0, stopByteWriter, false, []string{"B:"}},
		{"xi", true, false, []byteReply{{1, nil}}, 2, nil, false, []string{"B:\xff\x00", "B:i"}},
		{"qx", false, false, nil, 3, nil, false, []string{"B:q", "B:\xff\x00"}},
		{"x!b", true, false, []byteReply{{1, stopByteWriter}}, 1, stopByteWriter, false, []string{"B:\xff\x00"}},
		{"", true, true, nil, 0, nil, false, nil},
		{"x", true, true, nil, 0, nil, true, nil},
	}
	for _, c := range []struct {
		name string
		drop byte
		want int
	}{
		{"Go", 0, 0},
		{"Treat present empty mapping as absent", 'b', 3},
		{"Treat identity mapping as absent", 'i', 2},
	} {
		bad := 0
		active := mapping
		if c.drop != 0 {
			active[c.drop] = nil
		}
		for _, f := range fixtures {
			rec := &byteWriteRecorder{replies: f.replies}
			var writer io.Writer = rec
			if f.dual {
				writer = byteWriteDual{rec}
			}
			table := &active
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
