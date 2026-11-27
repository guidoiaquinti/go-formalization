package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"unicode/utf8"
	"unsafe"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	backing := runtime.GOOS + "\x00\xffpayload"
	input := backing[len(runtime.GOOS):]
	got := strings.Clone(input)
	if got != input || !disjoint(got, backing) {
		log.Fatalf("Clone did not preserve bytes with disjoint storage: %q", got)
	}
	reused := input
	if disjoint(reused, backing) {
		log.Fatal("identity mutant unexpectedly used disjoint storage")
	}
	if empty := strings.Clone(backing[1:1]); empty != "" {
		log.Fatalf("empty Clone returned %q", empty)
	}

	cases := []struct {
		text string
		r    rune
		want int
	}{
		{"aéé", 'é', 1},
		{"a\xff�", utf8.RuneError, 1},
		{"猫�", utf8.RuneError, 3},
		{"猫\xe0\xa0", utf8.RuneError, 3},
		{"x�", 0xd800, -1},
		{"x�", -1, -1},
		{"x\x00y", 0, 1},
	}
	var failures int
	for _, c := range cases {
		if got := strings.IndexRune(c.text, c.r); got != c.want {
			log.Fatalf("IndexRune(%q, %#x) = %d, want %d", c.text, c.r, got, c.want)
		}
		if strings.Index(c.text, string(c.r)) != c.want {
			failures++
		}
	}
	if failures != 4 {
		log.Fatalf("encoded-string-only mutant failed %d fixtures, want four", failures)
	}
	fmt.Printf("%s %s/%s: Clone and IndexRune controls passed, identity and encoded-string-only mutants rejected\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func disjoint(a, b string) bool {
	first := uintptr(unsafe.Pointer(unsafe.StringData(a)))
	second := uintptr(unsafe.Pointer(unsafe.StringData(b)))
	result := first+uintptr(len(a)) <= second || second+uintptr(len(b)) <= first
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	return result
}
