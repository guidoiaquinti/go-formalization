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
	words := []uint64{0, 0x7fefffffffffffff, 0x7ff0000000000000, 0x7ff0000000000001,
		0x7ff8000000000000, 0x7fffffffffffffff}
	var actual, exponentOnly, quietOnly int
	for _, sign := range []uint64{0, 1 << 63} {
		for _, code := range words {
			word := code | sign
			exponent, fraction := (word>>52)&0x7ff, word&((1<<52)-1)
			want := exponent == 0x7ff && fraction != 0
			if math.IsNaN(math.Float64frombits(word)) != want {
				actual++
			}
			if (exponent == 0x7ff) != want {
				exponentOnly++
			}
			if (exponent == 0x7ff && fraction&(1<<51) != 0) != want {
				quietOnly++
			}
		}
	}
	if actual != 0 || exponentOnly != 2 || quietOnly != 2 {
		log.Fatalf("IsNaN mismatches: Go=%d exponent-only=%d quiet-only=%d", actual, exponentOnly, quietOnly)
	}
	fmt.Printf("%s %s/%s: IsNaN Go=%d/12 exponent-only=%d/12 quiet-only=%d/12\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, actual, exponentOnly, quietOnly)
}
