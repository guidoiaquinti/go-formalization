package strings_test

import (
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

type writer func(*[]byte, string) (int, error)

func observe(fn writer, dst *[]byte, src string) (n int, err error, panicked bool) {
	defer func() {
		if fault := recover(); fault != nil {
			value, ok := fault.(runtime.Error)
			if !ok || value.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(fault)
			}
			panicked = true
		}
	}()
	n, err = fn(dst, src)
	return n, err, false
}

func normalize(dst *[]byte, src string) (int, error) {
	return strings.GoLeanAppendSliceWriterWriteString(dst, strings.ToValidUTF8(src, "�"))
}

func allocateEmpty(dst *[]byte, src string) (int, error) {
	n, err := strings.GoLeanAppendSliceWriterWriteString(dst, src)
	if *dst == nil {
		*dst = make([]byte, 0)
	}
	return n, err
}

func TestGoLeanAppendSliceWriterWriteStringReview(t *testing.T) {
	for _, c := range []struct {
		name  string
		write writer
		want  int
	}{
		{"Go", strings.GoLeanAppendSliceWriterWriteString, 0},
		{"Normalize malformed string bytes", normalize, 1},
		{"Always allocate an empty buffer", allocateEmpty, 1},
	} {
		disjoint := []byte("abcde")
		aliased := disjoint[:2:5]
		exact := []byte("é")
		grow := exact[:len(exact):len(exact)]
		var noStorage []byte
		empty := make([]byte, 0, 4)
		full := []byte{0xff}
		fixtures := []struct {
			dst                *[]byte
			src                string
			want               string
			capacity           int
			wantNil, wantPanic bool
		}{
			{&aliased, unsafe.String(unsafe.SliceData(disjoint), 2), "abab", 5, false, false},
			{&grow, unsafe.String(unsafe.SliceData(exact), len(exact)), "éé", -1, false, false},
			{&noStorage, "", "", 0, true, false},
			{&empty, "", "", 4, false, false},
			{nil, "", "", 0, false, true},
			{&full, "\x00\x80", "\xff\x00\x80", -1, false, false},
		}
		bad := 0
		for _, f := range fixtures {
			n, err, panicked := observe(c.write, f.dst, f.src)
			if panicked != f.wantPanic || (!panicked && (n != len(f.src) || err != nil ||
				(*f.dst == nil) != f.wantNil || string(*f.dst) != f.want ||
				(f.capacity >= 0 && cap(*f.dst) != f.capacity) || cap(*f.dst) < len(*f.dst))) {
				bad++
				if c.want == 0 {
					t.Errorf("WriteString(%x) got count=%d err=%v panic=%t", f.src, n, err, panicked)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
