package main

import (
	"fmt"
	"math"
	"runtime"
)

func main() {
	words := []uint64{
		0, 1 << 63,
		1, 1<<63 | 1,
		0x3ff0000000000000, 0xbff0000000000000,
		0x7ff0000000000000, 0xfff0000000000000,
		0x7ff8000000000001, 0xfff8000000000001,
	}
	real, mutant := 0, 0
	for i, word := range words {
		x := math.Float64frombits(word)
		want := i%2 == 1 // Each pair contains the positive encoding first.
		if math.Signbit(x) != want {
			real++
		}
		if (x < 0) != want {
			mutant++
		}
	}
	fmt.Printf("%s %s/%s: Signbit Go=%d/%d numeric-comparison=%d/%d\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, len(words), mutant, len(words))
	if real != 0 || mutant != 2 {
		panic("unexpected Signbit control counts")
	}
}
