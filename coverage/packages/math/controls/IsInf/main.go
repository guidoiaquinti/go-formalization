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
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected Go version %s", runtime.Version())
	}
	words := []uint64{0, 0x7fefffffffffffff, 0x7ff0000000000000, 0x7ff8000000000001}
	var actual, ignoresSign, exponentOnly int
	for _, high := range []uint64{0, 1 << 63} {
		for _, code := range words {
			word := high | code
			for _, sign := range []int{-1, 0, 1} {
				matchesSign := sign == 0 || (sign < 0) == (high != 0)
				isInfinity := code == 0x7ff0000000000000
				want := matchesSign && isInfinity
				if math.IsInf(math.Float64frombits(word), sign) != want {
					actual++
				}
				if isInfinity != want {
					ignoresSign++
				}
				if (matchesSign && code>>52 == 0x7ff) != want {
					exponentOnly++
				}
			}
		}
	}
	if actual != 0 || ignoresSign != 2 || exponentOnly != 4 {
		log.Fatalf("IsInf mismatches: Go=%d sign ignored=%d exponent only=%d", actual, ignoresSign, exponentOnly)
	}
	fmt.Printf("%s %s/%s: IsInf Go=%d/24 sign ignored=%d/24 exponent only=%d/24\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, actual, ignoresSign, exponentOnly)
}
