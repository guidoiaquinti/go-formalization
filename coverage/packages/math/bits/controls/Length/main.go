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
	check(8, func(x uint32) int { return bits.Len8(uint8(x)) })
	check(16, func(x uint32) int { return bits.Len16(uint16(x)) })
}

func check(width int, goLen func(uint32) int) {
	total := uint32(1) << width
	var real, population, fixed int
	for x := range total {
		if !isBitLength(x, goLen(x), width) {
			real++
		}
		if !isBitLength(x, bits.OnesCount32(x), width) {
			population++
		}
		if !isBitLength(x, width, width) {
			fixed++
		}
	}
	if real != 0 || population != int(total)-width-1 || fixed != int(total)/2 {
		log.Fatalf("Len%d failures: Go=%d population=%d fixed=%d", width, real, population, fixed)
	}
	fmt.Printf("%s %s/%s: Len%d Go=%d/%d population=%d/%d fixed=%d/%d\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, width, real, total, population, total, fixed, total)
}

// Check the quantified fit criterion for all relevant k. Both sides remain
// true beyond width because the input fits in width and n has been bounded.
func isBitLength(x uint32, n, width int) bool {
	if n < 0 || n > width {
		return false
	}
	for k := 0; k <= width; k++ {
		if (n <= k) != (x < uint32(1)<<k) {
			return false
		}
	}
	return true
}
