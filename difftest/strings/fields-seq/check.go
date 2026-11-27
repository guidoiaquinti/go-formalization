package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s      string
	budget int
}

var spaces = []rune("\t\n\v\f\r \u0085\u00a0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u2028\u2029\u202f\u205f\u3000")

func main() {
	values := []string{"", "a", "  a   b  c  ", "one two", "one\xa0two", "\xc2 \xc2\xa0a\xff\xe2\x80", "a\u200bb\uFEFFc"}
	for b := range 256 {
		values = append(values, "a"+string([]byte{byte(b)})+"b")
	}
	for _, space := range spaces {
		for _, r := range []rune{space - 1, space, space + 1} {
			values = append(values, string(r), "a"+string(r)+"b"+string(r)+string(r)+"c")
		}
	}
	var edges []input
	for _, s := range values {
		for _, budget := range []int{0, 1, 2, 3, len(s) + 1} {
			edges = append(edges, input{s, budget})
		}
	}
	check.CompareSamples("godiff-strings-fields-seq", edges, sample, encodeInput, actual, expected)
}

func sample(rng *rand.Rand, _ int) input {
	var s strings.Builder
	for range rng.Intn(65) {
		switch rng.Intn(5) {
		case 0:
			s.WriteRune(spaces[rng.Intn(len(spaces))])
		case 1:
			s.WriteRune(rune(rng.Intn(0x110000)))
		default:
			s.WriteByte(byte(rng.Intn(256)))
		}
	}
	return input{s.String(), rng.Intn(s.Len() + 2)}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(c.budget), 10)
	for i := range len(c.s) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.s[i]), 10)
	}
	return string(row)
}

func describe(fields []string) string {
	encoded := make([]string, len(fields))
	for i, field := range fields {
		parts := make([]string, len(field))
		for j := range len(field) {
			parts[j] = strconv.Itoa(int(field[j]))
		}
		encoded[i] = strings.Join(parts, ",")
	}
	return strings.Join(encoded, ";")
}

func actual(c input) string {
	seq := strings.FieldsSeq(c.s)
	var limited, full []string
	if c.budget != 0 {
		seq(func(field string) bool {
			if field == "" || len(limited) >= c.budget {
				log.Fatalf("strings.FieldsSeq(%q) yielded %q past the callback contract", c.s, field)
			}
			limited = append(limited, field)
			return len(limited) < c.budget
		})
	}
	seq(func(field string) bool {
		if field == "" {
			log.Fatalf("strings.FieldsSeq(%q) yielded an empty field", c.s)
		}
		full = append(full, field)
		return true
	})
	return describe(limited) + "|" + describe(full)
}

func expected(c input) string {
	fields := strings.Fields(c.s)
	return describe(fields[:min(c.budget, len(fields))]) + "|" + describe(fields)
}
