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
	run("Mul", 64, func(x, y uint64) (uint64, uint64) {
		hi, lo := bits.Mul(uint(x), uint(y))
		return uint64(hi), uint64(lo)
	})
	run("Mul32", 32, func(x, y uint64) (uint64, uint64) {
		hi, lo := bits.Mul32(uint32(x), uint32(y))
		return uint64(hi), uint64(lo)
	})
}

func run(name string, width int, multiply func(uint64, uint64) (uint64, uint64)) {
	mask := ^uint64(0) >> (64 - width)
	highBit := uint64(1) << (width - 1)
	limb := uint64(1) << (width / 2)
	cases := []struct{ x, y, high, low uint64 }{
		{0, mask, 0, 0},
		{1, mask, 0, mask},
		{mask, mask, mask - 1, 1},
		{highBit, 2, 1, 0},
		{limb, limb, 1, 0},
	}
	var truncatedFailures, swappedFailures int
	for _, c := range cases {
		high, low := multiply(c.x, c.y)
		if high != c.high || low != c.low {
			log.Fatalf("%s(%#x, %#x) = (%#x, %#x), want (%#x, %#x)", name, c.x, c.y, high, low, c.high, c.low)
		}
		if c.high != 0 || ((c.x*c.y)&mask) != c.low {
			truncatedFailures++
		}
		if low != c.high || high != c.low {
			swappedFailures++
		}
	}
	if truncatedFailures != 3 || swappedFailures != 4 {
		log.Fatalf("%s mutant failures truncated=%d swapped=%d, want three and four", name, truncatedFailures, swappedFailures)
	}
	fmt.Printf("%s %s/%s: %s passed five fixtures; truncated and swapped products rejected\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, name)
}
