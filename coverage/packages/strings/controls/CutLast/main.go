package main

import (
	"fmt"
	"log"
	"strings"
)

type split struct {
	before string
	after  string
	found  bool
}

func main() {
	mutants := []struct {
		name string
		cut  func(string, string) split
		want int
	}{
		{"first occurrence", firstOccurrence, 5},
		{"nonoverlapping search", nonoverlappingSearch, 1},
		{"separator included in suffix", separatorInSuffix, 5},
		{"missing separator discards input", discardInputOnMiss, 2},
	}
	tests := []struct {
		input     string
		separator string
		want      split
	}{
		{"prefix::middle::suffix", "::", split{"prefix::middle", "suffix", true}},
		{"ababa", "aba", split{"ab", "", true}},
		{"front", "", split{"front", "", true}},
		{"", "", split{"", "", true}},
		{"ab", "abc", split{"ab", "", false}},
		{"no match", "::", split{"no match", "", false}},
		{"\x00\xff\x00\xfftail", "\x00\xff", split{"\x00\xff", "tail", true}},
		{"same", "same", split{"", "", true}},
		{"a\x00b\x00c", "\x00", split{"a\x00b", "c", true}},
	}
	disagreements := make([]int, len(mutants))
	for _, test := range tests {
		before, after, found := strings.CutLast(test.input, test.separator)
		if got := (split{before, after, found}); got != test.want {
			log.Fatalf("CutLast(%q, %q) = %+v, want %+v", test.input, test.separator, got, test.want)
		}
		for i, mutant := range mutants {
			if mutant.cut(test.input, test.separator) != test.want {
				disagreements[i]++
			}
		}
	}
	fmt.Printf("strings.CutLast matched all %d fixtures\n", len(tests))
	for i, mutant := range mutants {
		if disagreements[i] != mutant.want {
			log.Fatalf("mutant %q rejected on %d fixtures, want %d", mutant.name, disagreements[i], mutant.want)
		}
		fmt.Printf("mutant %q rejected on %d fixtures\n", mutant.name, disagreements[i])
	}
}

func firstOccurrence(input, separator string) split {
	before, after, found := strings.Cut(input, separator)
	return split{before, after, found}
}

func nonoverlappingSearch(input, separator string) split {
	if separator == "" {
		return split{input, "", true}
	}
	lastEnd := 0
	for {
		next := strings.Index(input[lastEnd:], separator)
		if next < 0 {
			break
		}
		lastEnd += next + len(separator)
	}
	if lastEnd == 0 {
		return split{input, "", false}
	}
	return split{input[:lastEnd-len(separator)], input[lastEnd:], true}
}

func separatorInSuffix(input, separator string) split {
	index := strings.LastIndex(input, separator)
	if index < 0 {
		return split{input, "", false}
	}
	return split{input[:index], input[index:], true}
}

func discardInputOnMiss(input, separator string) split {
	before, after, found := strings.CutLast(input, separator)
	if !found {
		return split{}
	}
	return split{before, after, found}
}
