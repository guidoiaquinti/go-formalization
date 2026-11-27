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
	const signWeight uint64 = 1 << 63
	var words []uint64
	for _, magnitude := range []uint64{0, 1, 0x3ff0000000000000, 0x7ff8000000000001} {
		words = append(words, magnitude, magnitude+signWeight)
	}
	var real, orderedSign int
	for _, x := range words {
		for _, y := range words {
			got := math.Float64bits(math.Copysign(math.Float64frombits(x), math.Float64frombits(y)))
			if got%signWeight != x%signWeight || got/signWeight != y/signWeight {
				real++
			}
			mutant := x % signWeight
			if math.Float64frombits(y) < 0 {
				mutant += signWeight
			}
			if mutant/signWeight != y/signWeight {
				orderedSign++
			}
		}
	}
	if real != 0 || orderedSign != 16 {
		log.Fatalf("Copysign failures: Go=%d, ordered sign=%d, want 0 and 16", real, orderedSign)
	}
	fmt.Printf("%s %s/%s: Copysign Go=%d/64, ordered sign mutant=%d/64\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, orderedSign)
}
