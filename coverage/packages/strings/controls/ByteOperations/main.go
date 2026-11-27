package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	cmpCases := []struct {
		a, b string
		want int
	}{
		{"", "", 0},
		{"a", "aa", -1},
		{"z", "aa", 1},
		{"\xff", "\x01", 1},
		{"a\x00z", "a\x00a", 1},
	}
	var signedFailures int
	for _, c := range cmpCases {
		if got := strings.Compare(c.a, c.b); got != c.want {
			log.Fatalf("Compare(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
		mutant := cmp.Compare(len(c.a), len(c.b))
		for i := range min(len(c.a), len(c.b)) {
			if order := cmp.Compare(int8(c.a[i]), int8(c.b[i])); order != 0 {
				mutant = order
				break
			}
		}
		if mutant != c.want {
			signedFailures++
		}
	}
	findCases := []struct {
		text        string
		key         byte
		first, last int
	}{
		{"", 'a', -1, -1},
		{"abc", 'd', -1, -1},
		{"\x00\xff\x00", 0, 0, 2},
		{"éé", 0xa9, 1, 3},
		{"abc", 'b', 1, 1},
	}
	var forwardFailures, backwardFailures int
	for _, c := range findCases {
		first, last := strings.IndexByte(c.text, c.key), strings.LastIndexByte(c.text, c.key)
		if first != c.first || last != c.last {
			log.Fatalf("byte search(%q, %d) = (%d,%d), want (%d,%d)", c.text, c.key, first, last, c.first, c.last)
		}
		if last != c.first {
			forwardFailures++
		}
		if first != c.last {
			backwardFailures++
		}
	}
	if signedFailures != 1 || forwardFailures != 2 || backwardFailures != 2 {
		log.Fatalf("mutant failure counts signed=%d forward=%d backward=%d, want 1, 2, 2",
			signedFailures, forwardFailures, backwardFailures)
	}
	fmt.Printf("%s %s/%s: byte operation controls passed; signed-byte and swapped directions rejected\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
