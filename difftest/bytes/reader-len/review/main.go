package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	present bool
	length  int
	offset  int64
	want    string
}

func main() {
	cases := []fixture{
		{false, 0, 0, "panic"},
		{true, 0, 0, "0"},
		{true, 0, math.MaxInt64, "0"},
		{true, 6, 0, "6"},
		{true, 6, 5, "1"},
		{true, 6, 6, "0"},
		{true, 6, 7, "0"},
		{true, 6, 1 << 32, "0"},
		{true, 6, math.MaxInt64, "0"},
	}
	var input bytes.Buffer
	for _, f := range cases {
		if f.present {
			fmt.Fprintln(&input, 1, f.length, f.offset)
		} else {
			fmt.Fprintln(&input, 0)
		}
	}
	for i, got := range check.LeanRows("godiff-bytes-reader-len", len(cases), &input) {
		f := cases[i]
		if native := observe(f); got != f.want || native != f.want {
			log.Fatalf("fixture %d: Go=%s Lean=%s want=%s", i, native, got, f.want)
		}
	}
	r := bytes.NewReader([]byte("é!"))
	if value, width, err := r.ReadRune(); value != 'é' || width != 2 || err != nil ||
		r.Len() != 1 || r.Len() != 1 || r.UnreadRune() != nil || r.Len() != 3 {
		log.Fatal("Len changed observable cursor or previous-rune state")
	}
	if _, err := r.Seek(-1, io.SeekStart); err == nil || r.Len() != 3 {
		log.Fatal("reader admitted a negative position")
	}
	if cases[6].want == strconv.FormatInt(int64(cases[6].length)-cases[6].offset, 10) ||
		cases[7].want == strconv.Itoa(cases[7].length-int(int32(cases[7].offset))) {
		log.Fatal("incorrect clamp or cursor truncation survived")
	}
	fmt.Printf("Reader.Len: %d native and Lean fixtures plus cursor and rune-state sequence passed\n", len(cases))
}

func observe(f fixture) (result string) {
	defer func() {
		if failure := recover(); failure != nil {
			e, ok := failure.(runtime.Error)
			if f.present || !ok || e.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			result = "panic"
		}
	}()
	var r *bytes.Reader
	if f.present {
		r = bytes.NewReader(make([]byte, f.length))
		if _, err := r.Seek(f.offset, io.SeekStart); err != nil {
			log.Fatal(err)
		}
	}
	return strconv.Itoa(r.Len())
}
