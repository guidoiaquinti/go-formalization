package main

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
)

func main() {
	fixtures := []struct {
		in, sep string
		want    []string
		repeat  []string
	}{
		{"", "", nil, nil},
		{"", ":", []string{""}, []string{""}},
		{"left::right", ":", []string{"left:", ":", "right"}, []string{"right"}},
		{"a\x00b\x00", "\x00", []string{"a\x00", "b\x00", ""}, []string{""}},
		{"é\xff", "", []string{"é", "\xff"}, nil},
		{"\xff!\xfe", "!", []string{"\xff!", "\xfe"}, []string{"\xfe"}},
	}
	var uncappedFailures, separatorFailures int
	for _, f := range fixtures {
		backing := make([]byte, len(f.in), len(f.in)+4)
		copy(backing, f.in)
		seq := bytes.SplitAfterSeq(backing, []byte(f.sep))
		var stopped, rest, repeated []string
		observe := func(part []byte) string {
			if cap(part) != len(part) {
				log.Fatalf("%q capacity = %d, want %d", part, cap(part), len(part))
			}
			_ = append(part, 0xfc)
			if string(backing) != f.in || !bytes.Equal(backing[len(backing):cap(backing)], []byte{0, 0, 0, 0}) {
				log.Fatal("append changed source backing bytes")
			}
			return string(part)
		}
		seq(func(part []byte) bool {
			stopped = append(stopped, observe(part))
			return false
		})
		seq(func(part []byte) bool {
			rest = append(rest, observe(part))
			return true
		})
		seq(func(part []byte) bool {
			repeated = append(repeated, observe(part))
			return true
		})
		var first []string
		if len(f.want) != 0 {
			first = f.want[:1]
		}
		if !reflect.DeepEqual(stopped, first) || !reflect.DeepEqual(rest, f.want) || !reflect.DeepEqual(repeated, f.repeat) {
			log.Fatalf("SplitAfterSeq(%q, %q) yielded %q, %q, %q", f.in, f.sep, stopped, rest, repeated)
		}
		if len(first) != 0 {
			part := backing[:len(first[0])]
			_ = append(part, 0xfc)
			if !bytes.Equal(backing[:cap(backing)], append([]byte(f.in), 0, 0, 0, 0)) {
				uncappedFailures++
			}
		}
		if !reflect.DeepEqual(f.want, bytesWithoutSeparators(f.in, f.sep)) {
			separatorFailures++
		}
	}
	if uncappedFailures != 5 || separatorFailures != 3 {
		log.Fatalf("uncapped mutant failed %d, omit-separator mutant failed %d", uncappedFailures, separatorFailures)
	}
	fmt.Printf("Go passed %d fixtures; uncapped and omit-separator mutants failed %d and %d\n", len(fixtures), uncappedFailures, separatorFailures)
}

func bytesWithoutSeparators(input, separator string) []string {
	var results []string
	for part := range bytes.SplitSeq([]byte(input), []byte(separator)) {
		results = append(results, string(part))
	}
	return results
}
