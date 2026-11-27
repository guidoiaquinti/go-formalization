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
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected Go version %s", runtime.Version())
	}
	cases := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"ascii", 5},
		{"Aé世🌍", 4},
		{"\xef\xbf\xbd", 1},     // Encoded RuneError occupies three bytes.
		{"\xc0\x80", 2},         // Overlong form.
		{"\xed\xa0\x80", 3},     // Surrogate.
		{"\xf4\x90\x80\x80", 4}, // Beyond Unicode.
		{"\xf0\x90\x80", 3},     // Truncated form.
		{"\xe2\x82Aé", 4},       // Malformed prefix followed by ASCII and a valid rune.
	}
	var sliceFailures, stringFailures, byteLength, runeStarts int
	for _, c := range cases {
		if got := utf8.RuneCount([]byte(c.input)); got != c.want || got > len(c.input) {
			sliceFailures++
		}
		if got := utf8.RuneCountInString(c.input); got != c.want || got > len(c.input) {
			stringFailures++
		}
		if len(c.input) != c.want {
			byteLength++
		}
		starts := 0
		for i := range len(c.input) {
			if utf8.RuneStart(c.input[i]) {
				starts++
			}
		}
		if starts != c.want {
			runeStarts++
		}
	}
	if sliceFailures != 0 || stringFailures != 0 || byteLength != 3 || runeStarts != 5 {
		log.Fatalf("rune-count failures: slice=%d string=%d byte-length=%d rune-start=%d", sliceFailures, stringFailures, byteLength, runeStarts)
	}
	fmt.Printf("%s %s/%s: RuneCount=%d/%d RuneCountInString=%d/%d byte-length=%d/%d rune-start=%d/%d\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, sliceFailures, len(cases), stringFailures, len(cases),
		byteLength, len(cases), runeStarts, len(cases))
}
