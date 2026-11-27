package main

import (
	"fmt"
	"iter"
	"log"
	"reflect"
	"strings"
)

type splitter func(string, string) iter.Seq[string]

func observe(seq iter.Seq[string]) [3][]string {
	var batches [3][]string
	for i := range batches {
		seq(func(part string) bool {
			batches[i] = append(batches[i], part)
			return i == 1
		})
	}
	return batches
}

func advanceBeforeYield(s, sep string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for {
			i := strings.Index(s, sep)
			if i < 0 || sep == "" {
				strings.SplitSeq(s, sep)(yield)
				return
			}
			part := s[:i]
			s = s[i+len(sep):]
			if !yield(part) {
				return
			}
		}
	}
}

func keepSeparator(s, sep string) iter.Seq[string] { return strings.SplitAfterSeq(s, sep) }

func main() {
	fixtures := []struct {
		s, sep string
		want   [3][]string
	}{
		{"", ":", [3][]string{{""}, {""}, {""}}},
		{"", "", [3][]string{}},
		{"::a::", "::", [3][]string{{""}, {"", "a", ""}, {""}}},
		{"aaaaa", "aa", [3][]string{{""}, {"", "", "a"}, {"a"}}},
		{"é\xffx", "", [3][]string{{"é"}, {"é", "\xff", "x"}, nil}},
		{"abc", ":", [3][]string{{"abc"}, {"abc"}, {"abc"}}},
	}
	for _, candidate := range []struct {
		name string
		fn   splitter
		want int
	}{
		{"Go", strings.SplitSeq, 0},
		{"Advance state before yield", advanceBeforeYield, 3},
		{"Keep separators", keepSeparator, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(candidate.fn(f.s, f.sep)); !reflect.DeepEqual(got, f.want) {
				bad++
				if candidate.want == 0 {
					log.Fatalf("SplitSeq(%q,%q): got %q, want %q", f.s, f.sep, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, bad, len(fixtures))
		if bad != candidate.want {
			log.Fatalf("%s: wanted %d mismatches", candidate.name, candidate.want)
		}
	}
}
