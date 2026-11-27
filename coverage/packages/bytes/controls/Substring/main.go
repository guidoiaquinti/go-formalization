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
	var nilView = view{isNil: true}
	cases := []struct {
		data, sep     []byte
		first, last   int
		before, after view
	}{
		{nil, nil, 0, 0, nilView, nilView},
		{nil, []byte("a"), -1, -1, nilView, nilView},
		{make([]byte, 0, 3), nil, 0, 0, view{"", 3, false}, view{"", 3, false}},
		{data, []byte("aba"), 0, 2, view{"", 9, false}, view{"ba", 6, false}},
		{data, nil, 0, 5, view{"", 9, false}, view{"ababa", 9, false}},
		{data, []byte("abb"), -1, -1, view{"ababa", 9, false}, nilView},
		{data, []byte("ababa"), 0, 0, view{"", 9, false}, view{"", 4, false}},
		{[]byte{0xff, 0, 0xff, 1}, []byte{0xff, 1}, 2, 2, view{"\xff\x00", 4, false}, view{"", 0, false}},
		{[]byte("abac"), []byte("ac"), 2, 2, view{"ab", 4, false}, view{"", 0, false}},
	}
	var wrongDirection, ignoresSuffix, wrongCut, losesCapacity int
	for i, c := range cases {
		gotFirst, gotLast := bytes.Index(c.data, c.sep), bytes.LastIndex(c.data, c.sep)
		present := bytes.Contains(c.data, c.sep)
		before, after, found := bytes.Cut(c.data, c.sep)
		want := c.first >= 0
		if gotFirst != c.first || gotLast != c.last || present != want || found != want ||
			describe(before) != c.before || describe(after) != c.after {
			log.Fatalf("fixture %d: first=%d last=%d present=%t cut=%+v %+v %t",
				i, gotFirst, gotLast, present, describe(before), describe(after), found)
		}
		if gotFirst != gotLast {
			wrongDirection++
		}
		if firstByte := len(c.sep) == 0 || bytes.IndexByte(c.data, c.sep[0]) >= 0; firstByte != want {
			ignoresSuffix++
		}
		if found {
			if describe(c.data[:gotLast]) != c.before || describe(c.data[gotLast+len(c.sep):]) != c.after {
				wrongCut++
			}
			if describe(before[:len(before):len(before)]) != c.before || describe(after[:len(after):len(after)]) != c.after {
				losesCapacity++
			}
		}
	}
	if wrongDirection != 2 || ignoresSuffix != 1 || wrongCut != 2 || losesCapacity != 6 {
		log.Fatalf("mutation failures: direction=%d suffix=%d cut=%d capacity=%d",
			wrongDirection, ignoresSuffix, wrongCut, losesCapacity)
	}
	fmt.Printf("substring: %d fixtures pass, mutant failures direction=%d suffix=%d cut=%d capacity=%d\n",
		len(cases), wrongDirection, ignoresSuffix, wrongCut, losesCapacity)
}
