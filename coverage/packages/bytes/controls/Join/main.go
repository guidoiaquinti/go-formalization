package main

import (
	"bytes"
	"fmt"
	"log"
	"slices"
)

func main() {
	fixtures := []struct {
		parts [][]byte
		sep   string
		want  string
		isNil bool
	}{
		{nil, ":", "", false},
		{[][]byte{nil}, ":", "", true},
		{[][]byte{{}}, "", "", true},
		{[][]byte{[]byte("A")}, "skip", "A", false},
		{[][]byte{[]byte("a"), nil, []byte("b")}, ":", "a::b", false},
		{[][]byte{nil, nil}, "-", "-", false},
		{[][]byte{[]byte("\xff"), []byte("é")}, "\x00", "\xff\x00é", false},
	}
	var omitFailures, singletonAliasFailures int
	for _, f := range fixtures {
		sep := []byte(f.sep)
		input := make([][]byte, len(f.parts))
		for i, part := range f.parts {
			input[i] = slices.Clone(part)
		}
		got := bytes.Join(input, sep)
		if string(got) != f.want || (got == nil) != f.isNil {
			log.Fatalf("Join(%q, %q) = %q with nil=%t", input, sep, got, got == nil)
		}
		if len(input) > 1 && cap(got) != len(got) {
			log.Fatalf("result capacity %d, want %d", cap(got), len(got))
		}
		if string(bytes.Join(input, nil)) != f.want {
			omitFailures++
		}
		if len(got) > 0 {
			got[0] ^= 0xff
			for i, part := range input {
				if !bytes.Equal(part, f.parts[i]) {
					log.Fatal("result mutation changed input")
				}
			}
			if string(sep) != f.sep {
				log.Fatal("result mutation changed separator")
			}
			snapshot := bytes.Clone(got)
			for _, part := range input {
				if len(part) > 0 {
					part[0] ^= 0xff
				}
			}
			if len(sep) > 0 {
				sep[0] ^= 0xff
			}
			if !bytes.Equal(got, snapshot) {
				log.Fatal("input or separator mutation changed result")
			}
		}
		if len(f.parts) == 1 && len(f.parts[0]) > 0 && string(input[0]) != f.want {
			singletonAliasFailures++
		}
	}
	if omitFailures != 3 || singletonAliasFailures != 1 {
		log.Fatalf("omit-separator mutant failed %d, singleton alias mutant failed %d", omitFailures, singletonAliasFailures)
	}
	fmt.Printf("Go passed %d fixtures; mutants failed %d and %d\n", len(fixtures), omitFailures, singletonAliasFailures)
}
