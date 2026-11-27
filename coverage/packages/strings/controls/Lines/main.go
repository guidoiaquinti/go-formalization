package main

import (
	"fmt"
	"iter"
	"log"
	"reflect"
	"strings"
)

func batches(seq iter.Seq[string]) [3][]string {
	var result [3][]string
	seq(func(line string) bool {
		result[0] = append(result[0], line)
		return false
	})
	seq(func(line string) bool {
		result[1] = append(result[1], line)
		return true
	})
	seq(func(line string) bool {
		result[2] = append(result[2], line)
		return false
	})
	return result
}

func restart(s string) iter.Seq[string] {
	return func(yield func(string) bool) { strings.Lines(s)(yield) }
}

func omitLF(s string) iter.Seq[string] {
	seq := strings.Lines(s)
	return func(yield func(string) bool) {
		seq(func(line string) bool { return yield(strings.TrimSuffix(line, "\n")) })
	}
}

func main() {
	fixtures := []struct {
		s    string
		want [3][]string
	}{
		{"", [3][]string{}},
		{"\n", [3][]string{{"\n"}}},
		{"\n\n", [3][]string{{"\n"}, {"\n"}}},
		{"last", [3][]string{{"last"}}},
		{"first\r\n\xff\x00\nlast", [3][]string{{"first\r\n"}, {"\xff\x00\n", "last"}}},
		{"\u2028\r\x00", [3][]string{{"\u2028\r\x00"}}},
	}
	for _, implementation := range []struct {
		name  string
		lines func(string) iter.Seq[string]
		want  int
	}{
		{"Go", strings.Lines, 0},
		{"Restart on each call", restart, 5},
		{"Omit terminating LF", omitLF, 3},
	} {
		mismatches := 0
		for _, f := range fixtures {
			if actual := batches(implementation.lines(f.s)); !reflect.DeepEqual(actual, f.want) {
				mismatches++
				if implementation.want == 0 {
					log.Fatalf("Lines(%q) yielded %q, want %q", f.s, actual, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", implementation.name, mismatches, len(fixtures))
		if mismatches != implementation.want {
			log.Fatalf("%s: want %d mismatches", implementation.name, implementation.want)
		}
	}
}
