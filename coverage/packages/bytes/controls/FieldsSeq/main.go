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
		input string
		want  []string
	}{
		{"", nil},
		{"\t \u0085", nil},
		{"  a\tb ", []string{"a", "b"}},
		{"a\u00a0b\u3000c", []string{"a", "b", "c"}},
		{"\xff \xfe\x00", []string{"\xff", "\xfe\x00"}},
		{"a\u200bb", []string{"a\u200bb"}},
	}
	var asciiFailures, ignoreStopFailures int
	for _, f := range fixtures {
		src := make([]byte, len(f.input), len(f.input)+2)
		copy(src, f.input)
		seq := bytes.FieldsSeq(src)
		var stopped, complete []string
		see := func(part []byte) string {
			if cap(part) != len(part) {
				log.Fatalf("field %q has spare capacity", part)
			}
			_ = append(part, 0xfc)
			if string(src) != f.input || !bytes.Equal(src[len(src):cap(src)], []byte{0, 0}) {
				log.Fatalf("append modified backing bytes for %q", f.input)
			}
			return string(part)
		}
		seq(func(part []byte) bool {
			stopped = append(stopped, see(part))
			return false
		})
		for part := range seq {
			complete = append(complete, see(part))
		}
		if !slices.Equal(complete, f.want) || !slices.Equal(stopped, f.want[:min(1, len(f.want))]) {
			log.Fatalf("FieldsSeq(%q) yielded %q then %q", f.input, stopped, complete)
		}
		ascii := strings.FieldsFunc(f.input, func(r rune) bool {
			return r == ' ' || r >= '\t' && r <= '\r'
		})
		if !slices.Equal(ascii, f.want) {
			asciiFailures++
		}
		if !slices.Equal(complete, stopped) {
			ignoreStopFailures++
		}
	}
	if asciiFailures != 2 || ignoreStopFailures != 3 {
		log.Fatalf("ASCII-only mutant failed %d, ignore-stop mutant failed %d", asciiFailures, ignoreStopFailures)
	}
	fmt.Printf("Go passed %d fixtures; mutants failed %d and %d\n", len(fixtures), asciiFailures, ignoreStopFailures)
}
