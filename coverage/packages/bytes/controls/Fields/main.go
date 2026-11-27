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
		want  []string
	}{
		{nil, nil},
		{[]byte{}, nil},
		{[]byte(" \t\u0085"), nil},
		{[]byte("  a\tb "), []string{"a", "b"}},
		{[]byte("a\u00a0b\u3000c"), []string{"a", "b", "c"}},
		{[]byte("\xff \xfe\x00"), []string{"\xff", "\xfe\x00"}},
	}
	var asciiFailures, nilFailures int
	for _, f := range fixtures {
		src := slices.Clone(f.input)
		if src != nil {
			src = slices.Grow(src, 2)
		}
		fields := bytes.Fields(src)
		if fields == nil || len(fields) != len(f.want) {
			log.Fatalf("Fields(%q) returned %#v", src, fields)
		}
		for i, field := range fields {
			if string(field) != f.want[i] || cap(field) != len(field) {
				log.Fatalf("field %d = %q with capacity %d", i, field, cap(field))
			}
		}
		ascii := strings.FieldsFunc(string(src), func(r rune) bool {
			return r == ' ' || r >= '\t' && r <= '\r'
		})
		if !slices.Equal(ascii, f.want) {
			asciiFailures++
		}
		mutant := fields
		if len(mutant) == 0 {
			mutant = nil
		}
		if mutant == nil {
			nilFailures++
		}
	}
	if asciiFailures != 2 || nilFailures != 3 {
		log.Fatalf("ASCII-only mutant failed %d, nil-on-empty mutant failed %d", asciiFailures, nilFailures)
	}
	fmt.Printf("Go passed %d fixtures; mutants failed %d and %d\n", len(fixtures), asciiFailures, nilFailures)
}
