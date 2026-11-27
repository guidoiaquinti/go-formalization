package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"unicode/utf8"
)

func main() {
	word := strconv.IntSize / 8
	fixtures := []struct {
		in   []byte
		want bool
	}{
		{nil, true},
		{[]byte{}, true},
		{bytes.Repeat([]byte{'A'}, 8*word+2), true},
		{[]byte("A\xef\xbf\xbd"), true},
		{[]byte("A\xe2\x82"), false},
		{[]byte("A\xed\xa0\x80"), false},
		{[]byte("A\xf4\x90\x80\x80"), false},
		{append([]byte("A\xff"), bytes.Repeat([]byte{'A'}, 4*word)...), false},
		{append(bytes.Repeat([]byte{'A'}, 4*word+1), 0xff), false},
	}
	var skipFailures, replacementFailures int
	for _, f := range fixtures {
		before := bytes.Clone(f.in)
		if got := utf8.Valid(f.in); got != f.want || !bytes.Equal(f.in, before) {
			log.Fatalf("Valid(% x) = %v, want %v and unchanged input", before, got, f.want)
		}
		skipped := f.in
		if len(skipped) > word+1 && skipped[0] < utf8.RuneSelf {
			skipped = skipped[word+1:]
		}
		if utf8.Valid(skipped) != f.want {
			skipFailures++
		}
		if (utf8.Valid(f.in) && !bytes.Contains(f.in, []byte("\xef\xbf\xbd"))) != f.want {
			replacementFailures++
		}
	}
	if skipFailures != 1 || replacementFailures != 1 {
		log.Fatalf("mutant failures: unchecked ASCII block %d, reject encoded replacement %d", skipFailures, replacementFailures)
	}
	fmt.Printf("Go passed %d fixtures; both mutants failed one fixture\n", len(fixtures))
}
