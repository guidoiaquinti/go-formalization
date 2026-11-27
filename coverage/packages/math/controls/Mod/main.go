package main

import (
	"fmt"
	"log"
	"math"
)

func main() {
	const sign = uint64(1) << 63
	const inf = uint64(0x7ff0000000000000)
	const nan = inf + 1
	bits := math.Float64bits
	cases := []struct{ x, y, want uint64 }{
		{bits(5.5), bits(2), bits(1.5)},
		{bits(-5.5), bits(2), bits(-1.5)},
		{bits(5.5), bits(-2), bits(1.5)},
		{bits(-5.5), bits(-2), bits(-1.5)},
		{bits(-4), bits(2), sign},
		{sign, bits(2), sign},
		{sign | 7, 3, sign | 1},
		{7, sign | 3, 1},
		{0x0010000000000002, 3, 0},
		{0x7fefffffffffffff, 1, 0},
		{bits(-7), inf, bits(-7)},
		{inf, inf, nan},
		{bits(1), sign, nan},
		{nan, bits(1), nan},
		{bits(1), nan, nan},
	}
	matches := func(x float64, want uint64) bool {
		if want == nan {
			return math.IsNaN(x)
		}
		return bits(x) == want
	}
	var nearestFailures, zeroSignFailures int
	for _, c := range cases {
		x, y := math.Float64frombits(c.x), math.Float64frombits(c.y)
		actual := math.Mod(x, y)
		if !matches(actual, c.want) {
			log.Fatalf("Mod(%#016x, %#016x) = %#016x, want %#016x", c.x, c.y, bits(actual), c.want)
		}
		if !matches(math.Remainder(x, y), c.want) {
			nearestFailures++
		}
		if actual == 0 && !matches(0, c.want) {
			zeroSignFailures++
		}
	}
	if nearestFailures != 4 || zeroSignFailures != 2 {
		log.Fatalf("nearest quotient failed %d, lost zero sign failed %d", nearestFailures, zeroSignFailures)
	}
	fmt.Printf("Mod passed %d fixtures; nearest quotient failed %d; lost zero sign failed %d\n",
		len(cases), nearestFailures, zeroSignFailures)
}
