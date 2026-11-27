package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"
	"unicode/utf16"
)

func control(name string, actual, mutant bool) {
	if !actual || mutant {
		log.Fatalf("%s: Go matched=%t mutant matched=%t, want true and false", name, actual, mutant)
	}
	fmt.Printf("%s: Go matched, mutant rejected\n", name)
}

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected version %s", runtime.Version())
	}
	fmt.Printf("%s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	const firstOutside rune = 0xe000
	control("IsSurrogate", !utf16.IsSurrogate(firstOutside),
		!(0xd800 <= firstOutside && firstOutside <= 0xe000))

	const firstHigh rune = 0xd800
	acceptEveryBMP := func(r rune) int {
		if 0 <= r && r < 0x10000 {
			return 1
		}
		return utf16.RuneLen(r)
	}
	control("RuneLen", utf16.RuneLen(firstHigh) == -1, acceptEveryBMP(firstHigh) == -1)

	const scalar rune = 0x10401
	hi, lo := utf16.EncodeRune(scalar)
	matchesPair := func(hi, lo rune) bool { return hi == 0xd801 && lo == 0xdc01 }
	control("EncodeRune", matchesPair(hi, lo), matchesPair(hi, 0xdc00))

	const high, low rune = 0xd801, 0xdc02
	const point rune = 0x10402
	noScale := 0x10000 + high - 0xd800 + low - 0xdc00
	control("DecodeRune", utf16.DecodeRune(high, low) == point, noScale == point)

	prefix := []uint16{0, 0xffff}
	control("AppendRune", slices.Equal(utf16.AppendRune(prefix, scalar), []uint16{0, 0xffff, 0xd801, 0xdc01}),
		slices.Equal(utf16.AppendRune(nil, scalar), []uint16{0, 0xffff, 0xd801, 0xdc01}))

	encoded := utf16.Encode([]rune{0x110000})
	trimmed := encoded[:len(encoded):len(encoded)]
	control("Encode", slices.Equal(encoded, []uint16{0xfffd}) && cap(encoded) == 2,
		slices.Equal(trimmed, []uint16{0xfffd}) && cap(trimmed) == 2)

	words := []uint16{0xd800, 0xd801, 0xdc02, 0xdc00}
	want := []rune{0xfffd, 0x10402, 0xfffd}
	output := utf16.Decode(words)
	dropFollowingHigh := append([]rune{0xfffd}, utf16.Decode(words[2:])...)
	control("Decode", slices.Equal(output, want), slices.Equal(dropFollowingHigh, want))
}
