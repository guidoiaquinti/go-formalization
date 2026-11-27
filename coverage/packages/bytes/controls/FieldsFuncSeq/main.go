package main

import (
	"bytes"
	"fmt"
	"iter"
	"log"
	"reflect"
)

type fieldSeq func([]byte, func(rune) bool) iter.Seq[[]byte]

func observe(fn fieldSeq, s string) [2][]string {
	var out [2][]string
	batch := 0
	backing := make([]byte, len(s), len(s)+1)
	copy(backing, s)
	seq := fn(backing, func(r rune) bool {
		out[batch] = append(out[batch], fmt.Sprintf("p:%U", r))
		return r == '|' || r == '\u3000'
	})
	for batch = range out {
		seq(func(field []byte) bool {
			if cap(field) != len(field) {
				log.Fatalf("field %q capacity is %d", field, cap(field))
			}
			_ = append(field, 0xfc)
			if string(backing) != s || backing[:cap(backing)][len(s)] != 0 {
				log.Fatal("append changed source backing")
			}
			out[batch] = append(out[batch], "y:"+string(field))
			return batch == 1
		})
	}
	return out
}

func ignoreStop(s []byte, p func(rune) bool) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		bytes.FieldsFuncSeq(s, p)(func(field []byte) bool {
			yield(field)
			return true
		})
	}
}

func eager(s []byte, p func(rune) bool) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for _, field := range bytes.FieldsFunc(s, p) {
			if !yield(field) {
				return
			}
		}
	}
}

func main() {
	fixtures := []struct {
		s    string
		want [2][]string
	}{
		{"", [2][]string{}},
		{"||", [2][]string{{"p:U+007C", "p:U+007C"}, {"p:U+007C", "p:U+007C"}}},
		{"a|b", [2][]string{{"p:U+0061", "p:U+007C", "y:a"}, {"p:U+0061", "p:U+007C", "y:a", "p:U+0062", "y:b"}}},
		{"\xff\u3000é|z", [2][]string{
			{"p:U+FFFD", "p:U+3000", "y:\xff"},
			{"p:U+FFFD", "p:U+3000", "y:\xff", "p:U+00E9", "p:U+007C", "y:é", "p:U+007A", "y:z"},
		}},
		{"a", [2][]string{{"p:U+0061", "y:a"}, {"p:U+0061", "y:a"}}},
	}
	for _, candidate := range []struct {
		name string
		fn   fieldSeq
		want int
	}{
		{"Go", bytes.FieldsFuncSeq, 0},
		{"Ignore stop", ignoreStop, 2},
		{"Evaluate predicates eagerly", eager, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(candidate.fn, f.s); !reflect.DeepEqual(got, f.want) {
				bad++
				if candidate.want == 0 {
					log.Fatalf("%s(%q) trace = %q, want %q", candidate.name, f.s, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, bad, len(fixtures))
		if bad != candidate.want {
			log.Fatalf("%s: wanted %d disagreements", candidate.name, candidate.want)
		}
	}
}
