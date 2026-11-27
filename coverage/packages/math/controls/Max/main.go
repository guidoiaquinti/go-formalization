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
	negativeZero := math.Float64frombits(1 << 63)
	inf, nan := math.Inf(1), math.NaN()
	cases := []struct{ x, y, want float64 }{
		{inf, nan, inf}, {nan, inf, inf}, {-inf, nan, nan}, {nan, 1, nan},
		{negativeZero, 0, 0}, {0, negativeZero, 0}, {negativeZero, negativeZero, negativeZero},
		{-10, -5, -5},
	}
	var real, nanFirst, firstEqual int
	for _, c := range cases {
		got := math.Max(c.x, c.y)
		if !same(got, c.want) {
			real++
		}
		mutant := got
		if math.IsNaN(c.x) || math.IsNaN(c.y) {
			mutant = nan
		}
		if !same(mutant, c.want) {
			nanFirst++
		}
		mutant = got
		if c.x == c.y {
			mutant = c.x
		}
		if !same(mutant, c.want) {
			firstEqual++
		}
	}
	if real != 0 || nanFirst != 2 || firstEqual != 1 {
		log.Fatalf("Max mismatches: Go=%d NaN-first=%d first-equal=%d", real, nanFirst, firstEqual)
	}
	fmt.Printf("%s %s/%s: Max Go=%d/8 NaN-first=%d/8 first-equal=%d/8\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, nanFirst, firstEqual)
}

func same(x, y float64) bool {
	if math.IsNaN(y) {
		return math.IsNaN(x)
	}
	return math.Float64bits(x) == math.Float64bits(y)
}
