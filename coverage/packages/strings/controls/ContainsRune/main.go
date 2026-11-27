package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"unicode/utf8"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	cases := []struct {
		text string
		key  rune
		want bool
	}{
		{"", 0, false},
		{"é", 'é', true},
		{"é", 'e', false},
		{"a\xff", utf8.RuneError, true},
		{"\xe0\xa0", utf8.RuneError, true},
		{"a�", utf8.RuneError, true},
		{"a�", 0xd800, false},
		{"a�", -1, false},
	}
	var failures int
	for _, c := range cases {
		if got := strings.ContainsRune(c.text, c.key); got != c.want {
			log.Fatalf("ContainsRune(%q, %#x) = %t, want %t", c.text, c.key, got, c.want)
		}
		if strings.Contains(c.text, string(c.key)) != c.want {
			failures++
		}
	}
	if failures != 4 {
		log.Fatalf("encoded-string mutant failed %d cases, want four", failures)
	}
	fmt.Printf("%s %s/%s: ContainsRune passed eight fixtures; encoded-string mutant failed four\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
