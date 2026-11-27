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
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected Go version %s", runtime.Version())
	}
	check("Len32", 32, 130, 80, 93, func(x uint64) int { return bits.Len32(uint32(x)) })
	check("Len64", 64, 258, 160, 189, bits.Len64)
	if bits.UintSize == 64 {
		check("Len", 64, 258, 160, 189, func(x uint64) int { return bits.Len(uint(x)) })
	}
}

func check(name string, width, wantTotal, wantHalf, wantPopulation int, native func(uint64) int) {
	mask := ^uint64(0) >> (64 - width)
	halfMask := mask >> (width / 2)
	cases := []uint64{0, mask}
	for i := range width {
		bit := uint64(1) << i
		cases = append(cases, bit-1, bit, bit+1, mask^bit)
	}
	var halfFailures, populationFailures int
	for _, x := range cases {
		n := native(x)
		if !fits(x, n, width) {
			log.Fatalf("%s(%#x)=%d fails the bit-length criterion", name, x, n)
		}
		if !fits(x, native(x&halfMask), width) {
			halfFailures++
		}
		if !fits(x, bits.OnesCount64(x), width) {
			populationFailures++
		}
	}
	if len(cases) != wantTotal || halfFailures != wantHalf || populationFailures != wantPopulation {
		log.Fatalf("%s: %d cases, low-half mutant=%d, population mutant=%d", name, len(cases), halfFailures, populationFailures)
	}
	fmt.Printf("%s %s/%s: %s Go=0/%d low-half=%d/%d population=%d/%d\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, name,
		len(cases), halfFailures, len(cases), populationFailures, len(cases))
}

func fits(x uint64, n, width int) bool {
	if n < 0 || n > width {
		return false
	}
	// Right shifts also express the fit criterion for k == 64, where
	// a uint64 left shift could not represent the bound 2^64.
	for k := 0; k <= width; k++ {
		if (n <= k) != (x>>k == 0) {
			return false
		}
	}
	return true
}
