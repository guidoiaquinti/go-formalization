package strings_test

import (
	"runtime"
	"strings"
	"testing"
)

type writer func(*[]byte, []byte) (int, error)

func observe(fn writer, dst *[]byte, src []byte) (n int, err error, panicked bool) {
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

func appendOneByteAtATime(dst *[]byte, src []byte) (int, error) {
	if dst == nil {
		return strings.GoLeanAppendSliceWriterWrite(dst, src)
	}
	for _, b := range src {
		*dst = append(*dst, b)
	}
	return len(src), nil
}

func allocateEmpty(dst *[]byte, src []byte) (int, error) {
	n, err := strings.GoLeanAppendSliceWriterWrite(dst, src)
	if *dst == nil {
		*dst = make([]byte, 0)
	}
	return n, err
}

func TestGoLeanAppendSliceWriterWriteReview(t *testing.T) {
	for _, c := range []struct {
		name  string
		write writer
		want  int
	}{
		{"Go", strings.GoLeanAppendSliceWriterWrite, 0},
		{"Read aliased bytes one at a time", appendOneByteAtATime, 1},
		{"Always allocate an empty buffer", allocateEmpty, 1},
	} {
		backing := []byte("abcde")
		aliased := backing[:2:5]
		var noStorage []byte
		empty := make([]byte, 0, 4)
		full := []byte{0xff}
		fixtures := []struct {
			dst                *[]byte
			src                []byte
			want               string
			capacity           int
			wantNil, wantPanic bool
		}{
			{&aliased, backing[1:4], "abbcd", 5, false, false},
			{&noStorage, []byte{}, "", 0, true, false},
			{&empty, nil, "", 4, false, false},
			{nil, []byte{}, "", 0, false, true},
			{&full, []byte{0x00, 0x80}, "\xff\x00\x80", -1, false, false},
		}
		bad := 0
		for _, f := range fixtures {
			n, err, panicked := observe(c.write, f.dst, f.src)
			if panicked != f.wantPanic || (!panicked && (n != len(f.src) || err != nil ||
				(*f.dst == nil) != f.wantNil || string(*f.dst) != f.want ||
				(f.capacity >= 0 && cap(*f.dst) != f.capacity) || cap(*f.dst) < len(*f.dst))) {
				bad++
				if c.want == 0 {
					t.Errorf("AppendSliceWriter.Write(%x) got count=%d err=%v panic=%t", f.src, n, err, panicked)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
