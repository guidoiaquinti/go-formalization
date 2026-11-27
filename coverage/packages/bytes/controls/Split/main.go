package main

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"strings"
)

func main() {
	fixtures := []struct {
		input []byte
		sep   string
		want  []string
		caps  []int
	}{
		{nil, "", []string{}, nil},
		{nil, ":", []string{""}, []int{0}},
		{[]byte{}, ":", []string{""}, []int{4}},
		{[]byte("a::b"), ":", []string{"a", "", "b"}, []int{1, 0, 5}},
		{[]byte("x\x00"), "\x00", []string{"x", ""}, []int{1, 4}},
		{[]byte("é"), "", []string{"é"}, []int{2}},
		{[]byte("\xffA"), "", []string{"\xff", "A"}, []int{1, 5}},
	}
	var capAll, normalized int
	for _, f := range fixtures {
		var input []byte
		if f.input != nil {
			input = make([]byte, len(f.input), len(f.input)+4)
			copy(input, f.input)
		}
		got := bytes.Split(input, []byte(f.sep))
		if got == nil || len(got) != len(f.want) {
			log.Fatalf("Split(%q, %q) returned %#v", input, f.sep, got)
		}
		changed := false
		for i, part := range got {
			if string(part) != f.want[i] || cap(part) != f.caps[i] {
				log.Fatalf("Split(%q, %q)[%d] = %q with capacity %d", input, f.sep, i, part, cap(part))
			}
			changed = changed || len(part) != f.caps[i]
		}
		if changed {
			capAll++
		}
		if !slices.Equal(strings.Split(strings.ToValidUTF8(string(input), "\ufffd"), f.sep), f.want) {
			normalized++
		}
	}
	if capAll != 4 || normalized != 1 {
		log.Fatalf("cap-all mutant failed %d, normalize-malformed mutant failed %d", capAll, normalized)
	}
	fmt.Printf("Go passed %d fixtures; mutants failed %d and %d\n", len(fixtures), capAll, normalized)
}
