package main

import (
	"fmt"
	"log"
	"math/bits"
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
	if bits.UintSize != 64 {
		log.Fatalf("native uint width %d, want 64", bits.UintSize)
	}
	for _, entry := range []struct {
		name      string
		width     int
		byteOrder bool
		call      func(uint64) uint64
	}{
		{"Reverse", 64, false, func(x uint64) uint64 { return uint64(bits.Reverse(uint(x))) }},
		{"Reverse8", 8, false, func(x uint64) uint64 { return uint64(bits.Reverse8(uint8(x))) }},
		{"Reverse16", 16, false, func(x uint64) uint64 { return uint64(bits.Reverse16(uint16(x))) }},
		{"Reverse32", 32, false, func(x uint64) uint64 { return uint64(bits.Reverse32(uint32(x))) }},
		{"Reverse64", 64, false, bits.Reverse64},
		{"ReverseBytes", 64, true, func(x uint64) uint64 { return uint64(bits.ReverseBytes(uint(x))) }},
	} {
		var mutantFailures int
		for bit := range entry.width {
			word := uint64(1) << bit
			flipped := uint64(1) << (entry.width - 1 - bit)
			swapped := uint64(1) << ((entry.width/8-1-bit/8)*8 + bit%8)
			want, mutant := flipped, swapped
			if entry.byteOrder {
				want, mutant = swapped, flipped
			}
			if got := entry.call(word); got != want {
				log.Fatalf("%s(%#x) = %#x, want %#x", entry.name, word, got, want)
			}
			if mutant != want {
				mutantFailures++
			}
		}
		if mutantFailures != entry.width {
			log.Fatalf("%s mutant failed %d inputs, want %d", entry.name, mutantFailures, entry.width)
		}
		fmt.Printf("%s %s/%s: %s passed all %d one-bit fixtures and rejected the opposite reversal\n",
			runtime.Version(), runtime.GOOS, runtime.GOARCH, entry.name, entry.width)
	}
}
