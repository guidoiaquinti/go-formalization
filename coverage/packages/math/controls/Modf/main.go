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
	cases := []struct{ input, whole, fraction float64 }{
		{0, 0, 0},
		{negativeZero, negativeZero, negativeZero},
		{-0.25, negativeZero, -0.25},
		{math.SmallestNonzeroFloat64, 0, math.SmallestNonzeroFloat64},
		{-1.25, -1, -0.25},
		{-2, -2, negativeZero},
		{math.Nextafter(1<<52, 0), 1<<52 - 1, 0.5},
		{math.Inf(1), math.Inf(1), math.NaN()},
		{math.Inf(-1), math.Inf(-1), math.NaN()},
		{math.Float64frombits(0xfff8000000000001), math.NaN(), math.NaN()},
	}
	var signFailures, floorFailures int
	for _, c := range cases {
		whole, fraction := math.Modf(c.input)
		if !same(whole, c.whole) || !same(fraction, c.fraction) {
			log.Fatalf("Modf(%#016x) = (%#016x, %#016x)", math.Float64bits(c.input),
				math.Float64bits(whole), math.Float64bits(fraction))
		}
		if !same(c.input-math.Trunc(c.input), c.fraction) {
			signFailures++
		}
		if !same(math.Floor(c.input), c.whole) {
			floorFailures++
		}
	}
	if signFailures != 2 || floorFailures != 2 {
		log.Fatalf("mutant failures: no sign copy=%d, flooring=%d, want two each", signFailures, floorFailures)
	}
	fmt.Printf("%s %s/%s: Modf passed %d fixtures, sign-copy and floor mutants rejected twice each\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases))
}

func same(got, want float64) bool {
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	return math.Float64bits(got) == math.Float64bits(want)
}
