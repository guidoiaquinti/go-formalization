package main

import (
	"bytes"
	"fmt"
	"log"
	"slices"
)

func main() {
	fixtures := []struct {
		input, sep string
		limit      int
		fragments  []string
		capacity   []int
	}{
		{"x", ":", 0, nil, nil},
		{"", "", -1, []string{}, nil},
		{"", ":", -1, []string{""}, []int{4}},
		{"a:b:c", ":", 2, []string{"a", "b:c"}, []int{1, 7}},
		{"a:b:c", ":", -1, []string{"a", "b", "c"}, []int{1, 1, 5}},
		{"é", "", 1, []string{"é"}, []int{6}},
		{"é", "", -1, []string{"é"}, []int{2}},
		{"\xff\x80", "", -1, []string{"\xff", "\x80"}, []int{1, 5}},
	}
	var capAllFailures, ignoreLimitFailures int
	for _, f := range fixtures {
		input := make([]byte, len(f.input), len(f.input)+4)
		copy(input, f.input)
		parts := bytes.SplitN(input, []byte(f.sep), f.limit)
		if (parts == nil) != (f.fragments == nil) || len(parts) != len(f.fragments) {
			log.Fatalf("SplitN(%q, %q, %d) returned %#v", f.input, f.sep, f.limit, parts)
		}
		capAllDiffers := false
		for i, part := range parts {
			if string(part) != f.fragments[i] || cap(part) != f.capacity[i] {
				log.Fatalf("fragment %d = (%q, cap %d), want (%q, cap %d)", i, part, cap(part), f.fragments[i], f.capacity[i])
			}
			capAllDiffers = capAllDiffers || len(part) != f.capacity[i]
		}
		if capAllDiffers {
			capAllFailures++
		}
		mutant := bytes.Split(input, []byte(f.sep))
		mutantText := make([]string, len(mutant))
		for i, part := range mutant {
			mutantText[i] = string(part)
		}
		if (mutant == nil) != (f.fragments == nil) || !slices.Equal(mutantText, f.fragments) {
			ignoreLimitFailures++
		}
	}
	if capAllFailures != 5 || ignoreLimitFailures != 2 {
		log.Fatalf("cap-all mutant failed %d, ignore-limit mutant failed %d", capAllFailures, ignoreLimitFailures)
	}
	fmt.Printf("Go passed %d fixtures; mutants failed %d and %d\n", len(fixtures), capAllFailures, ignoreLimitFailures)
}
