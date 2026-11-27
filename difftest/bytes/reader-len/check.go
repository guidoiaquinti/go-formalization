package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	length int
	offset int64
	mode   int
}

func main() {
	cases := []input{{mode: 0}, {mode: 1}, {mode: 1, offset: math.MaxInt64}}
	for length := range 65 {
		for offset := int64(0); offset <= int64(length)+1; offset++ {
			cases = append(cases, input{length, offset, 2})
		}
		cases = append(cases, input{length, math.MaxInt64, 2})
	}
	check.CompareSamples("godiff-bytes-reader-len", cases, sample, encode, native, expected)
}
func sample(rng *rand.Rand, i int) input {
	length := rng.Intn(257)
	offset := int64(rng.Intn(length + 3))
	if i%9 == 0 {
		offset = math.MaxInt64 - int64(i)
	}
	return input{length, offset, 2}
}
func encode(c input) string {
	if c.mode == 0 {
		return "0"
	}
	return fmt.Sprintf("1 %d %d", c.length, c.offset)
}
func native(c input) (out string) {
	var r *bytes.Reader
	switch c.mode {
	case 1:
		r = new(bytes.Reader)
	case 2:
		r = bytes.NewReader(make([]byte, c.length))
	}
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if c.mode != 0 || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			out = "panic"
		}
	}()
	if r == nil {
		return strconv.Itoa(r.Len())
	}
	if got, err := r.Seek(c.offset, io.SeekStart); err != nil || got != c.offset {
		panic("reader setup failed")
	}
	length := r.Len()
	if got, err := r.Seek(0, io.SeekCurrent); err != nil || got != c.offset {
		panic("Len moved the reader")
	}
	return strconv.Itoa(length)
}
func expected(c input) string {
	if c.mode == 0 {
		return "panic"
	}
	if c.offset >= int64(c.length) {
		return "0"
	}
	return strconv.FormatInt(int64(c.length)-c.offset, 10)
}
