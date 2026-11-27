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
	negativeZero := math.Copysign(0, -1)
	cases := []struct{ x, y, want float64 }{
		{negativeZero, 0, 0},
		{0, negativeZero, 0},
		{1, 2, 0},
		{2, 1, 1},
		{math.SmallestNonzeroFloat64, 0, math.SmallestNonzeroFloat64},
		{math.MaxFloat64, -math.MaxFloat64, math.Inf(1)},
		{math.Inf(1), math.Inf(1), math.NaN()},
		{math.Inf(-1), math.Inf(-1), math.NaN()},
		{math.NaN(), 1, math.NaN()},
		{1, math.NaN(), math.NaN()},
		{math.Inf(-1), 0, 0},
	}
	var loseNaN, preserveNegativeZero int
	for _, c := range cases {
		if got := math.Dim(c.x, c.y); !same(got, c.want) {
			log.Fatalf("Dim(%#016x, %#016x) = %#016x, want %#016x",
				math.Float64bits(c.x), math.Float64bits(c.y), math.Float64bits(got), math.Float64bits(c.want))
		}
		difference := c.x - c.y
		positiveOnly := 0.0
		if difference > 0 {
			positiveOnly = difference
		}
		if !same(positiveOnly, c.want) {
			loseNaN++
		}
		strictClamp := difference
		if difference < 0 {
			strictClamp = 0
		}
		if !same(strictClamp, c.want) {
			preserveNegativeZero++
		}
	}
	if loseNaN != 4 || preserveNegativeZero != 1 {
		log.Fatalf("mutant failures: NaN swallowed=%d, negative zero preserved=%d, want four and one", loseNaN, preserveNegativeZero)
	}
	fmt.Printf("%s %s/%s: Dim passed %d cases, positive-only and strict-clamp mutants failed four and one\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases))
}

func same(got, want float64) bool {
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	return math.Float64bits(got) == math.Float64bits(want)
}
