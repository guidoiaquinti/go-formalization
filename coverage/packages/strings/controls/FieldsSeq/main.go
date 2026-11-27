package main

import (
	"fmt"
	"iter"
	"log"
	"reflect"
	"strings"
)

func observe(seq iter.Seq[string]) [2][]string {
	var batches [2][]string
	seq(func(field string) bool {
		batches[0] = append(batches[0], field)
		return false
	})
	seq(func(field string) bool {
		batches[1] = append(batches[1], field)
		return true
	})
	return batches
}

func asciiOnly(s string) iter.Seq[string] {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || '\t' <= r && r <= '\r'
	})
	return func(yield func(string) bool) {
		for _, field := range fields {
			if !yield(field) {
				return
			}
		}
	}
}

func resume(s string) iter.Seq[string] {
	fields := strings.Fields(s)
	return func(yield func(string) bool) {
		for len(fields) != 0 {
			field := fields[0]
			fields = fields[1:]
			if !yield(field) {
				return
			}
		}
	}
}

func ignoreStop(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		strings.FieldsSeq(s)(func(field string) bool {
			yield(field)
			return true
		})
	}
}

func main() {
	fixtures := []struct {
		input string
		want  [2][]string
	}{
		{"", [2][]string{}},
		{"\t\r \v\n\f", [2][]string{}},
		{"\u0085\u3000", [2][]string{}},
		{" a\t\tb ", [2][]string{{"a"}, {"a", "b"}}},
		{"\u00a0x\u1680\xff\x00\u202fy", [2][]string{{"x"}, {"x", "\xff\x00", "y"}}},
		{"q\u200bq\ufeff", [2][]string{{"q\u200bq\ufeff"}, {"q\u200bq\ufeff"}}},
	}
	for _, candidate := range []struct {
		name   string
		fields func(string) iter.Seq[string]
		want   int
	}{
		{"Go", strings.FieldsSeq, 0},
		{"Only ASCII whitespace", asciiOnly, 2},
		{"Resume after stopping", resume, 3},
		{"Ignore stop", ignoreStop, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(candidate.fields(f.input)); !reflect.DeepEqual(got, f.want) {
				bad++
				if candidate.want == 0 {
					log.Fatalf("%s(%q) yielded %q, want %q", candidate.name, f.input, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, bad, len(fixtures))
		if bad != candidate.want {
			log.Fatalf("%s: want %d disagreements", candidate.name, candidate.want)
		}
	}
}
