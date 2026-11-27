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
		name  string
		width int
		call  func(uint64) int
	}{
		{"LeadingZeros", 64, func(x uint64) int { return bits.LeadingZeros(uint(x)) }},
		{"LeadingZeros8", 8, func(x uint64) int { return bits.LeadingZeros8(uint8(x)) }},
		{"LeadingZeros16", 16, func(x uint64) int { return bits.LeadingZeros16(uint16(x)) }},
		{"LeadingZeros32", 32, func(x uint64) int { return bits.LeadingZeros32(uint32(x)) }},
		{"LeadingZeros64", 64, bits.LeadingZeros64},
	} {
		if got := entry.call(0); got != entry.width {
			log.Fatalf("%s(0) = %d, want %d", entry.name, got, entry.width)
		}
		var reversedFailures int
		for bit := range entry.width {
			word := uint64(1) << bit
			want := entry.width - 1 - bit
			if got := entry.call(word); got != want {
				log.Fatalf("%s(%#x) = %d, want %d", entry.name, word, got, want)
			}
			if bits.TrailingZeros64(word) != want {
				reversedFailures++
			}
		}
		if reversedFailures != entry.width {
			log.Fatalf("%s reversed mutant failed %d inputs, want %d", entry.name, reversedFailures, entry.width)
		}
		fmt.Printf("%s %s/%s: %s passed zero and all %d one-bit words, reversed scan rejected\n",
			runtime.Version(), runtime.GOOS, runtime.GOARCH, entry.name, entry.width)
	}
}
