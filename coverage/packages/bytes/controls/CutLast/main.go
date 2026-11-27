package main

import (
	"bytes"
	"fmt"
	"log"
)

type view struct {
	content  string
	capacity int
	isNil    bool
}

func describe(s []byte) view {
	return view{string(s), cap(s), s == nil}
}

func main() {
	data := make([]byte, 5, 9)
	copy(data, "ababa")
	nilView := view{isNil: true}
	cases := []struct {
		s, sep        []byte
		found         bool
		before, after view
	}{
		{nil, nil, true, nilView, nilView},
		{nil, []byte("a"), false, nilView, nilView},
		{make([]byte, 0, 3), nil, true, view{"", 3, false}, view{"", 3, false}},
		{data, []byte("aba"), true, view{"ab", 9, false}, view{"", 4, false}},
		{data, nil, true, view{"ababa", 9, false}, view{"", 4, false}},
		{data, []byte("abb"), false, view{"ababa", 9, false}, nilView},
		{data, []byte("ababa"), true, view{"", 9, false}, view{"", 4, false}},
		{[]byte{0xff, 0, 0xff, 1}, []byte{0xff, 1}, true, view{"\xff\x00", 4, false}, view{"", 0, false}},
	}
	var firstMatch, clipped int
	for i, c := range cases {
		b, a, found := bytes.CutLast(c.s, c.sep)
		if found != c.found || describe(b) != c.before || describe(a) != c.after {
			log.Fatalf("fixture %d: got %+v %+v %t", i, describe(b), describe(a), found)
		}
		mb, ma, mf := bytes.Cut(c.s, c.sep)
		if mf != c.found || describe(mb) != c.before || describe(ma) != c.after {
			firstMatch++
		}
		if describe(b[:len(b):len(b)]) != c.before || describe(a[:len(a):len(a)]) != c.after {
			clipped++
		}
	}
	if firstMatch != 2 || clipped != 6 {
		log.Fatalf("mutation failures: first match=%d clipped capacity=%d", firstMatch, clipped)
	}
	fmt.Printf("CutLast: %d fixtures pass, mutant failures first match=%d clipped capacity=%d\n",
		len(cases), firstMatch, clipped)
}
