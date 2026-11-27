package main

import (
	"fmt"
	"log"
	"math"
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

	cases := []struct {
		input, want uint64
	}{
		{0x3fdfffffffffffff, 0},                  // Just below 0.5.
		{0xbfdfffffffffffff, 0x8000000000000000}, // Just above -0.5.
		{0x3fe0000000000000, 0x3ff0000000000000}, // 0.5.
		{0xbfe0000000000000, 0xbff0000000000000}, // -0.5.
		{0x4003ffffffffffff, 0x4000000000000000}, // Just below 2.5.
		{0x4004000000000000, 0x4008000000000000}, // 2.5.
		{0xc004000000000000, 0xc008000000000000}, // -2.5.
		{0x4004000000000001, 0x4008000000000000}, // Just above 2.5.
		{0x432fffffffffffff, 0x4330000000000000}, // Last half-integer below 2^52.
		{0x4330000000000001, 0x4330000000000001}, // Large odd integer.
		{0x7ff0000000000000, 0x7ff0000000000000}, // Positive infinity.
		{0xfff0000000000000, 0xfff0000000000000}, // Negative infinity.
	}
	var evenFailures int
	for _, test := range cases {
		input := math.Float64frombits(test.input)
		if got := math.Float64bits(math.Round(input)); got != test.want {
			log.Fatalf("Round(%#016x) produced %#016x, want %#016x", test.input, got, test.want)
		}
		if math.Float64bits(math.RoundToEven(input)) != test.want {
			evenFailures++
		}
	}
	if evenFailures != 4 {
		log.Fatalf("ties-to-even mutant failed %d cases, want 4", evenFailures)
	}
	for _, word := range []uint64{0x7ff0000000000001, 0xfff8000000000123} {
		if !math.IsNaN(math.Round(math.Float64frombits(word))) {
			log.Fatalf("Round(%#016x) did not return NaN", word)
		}
	}
	fmt.Printf("%s %s/%s: Round passed %d cases, ties-to-even mutant failed %d of %d exact-word cases\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases)+2, evenFailures, len(cases))
}
