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
	const tiny = uint64(0x0020000000000000)
	const half = uint64(0x7fdfffffffffffff)
	bits := math.Float64bits
	cases := []struct{ x, y, want uint64 }{
		{bits(0.5), bits(1), bits(0.5)},
		{bits(1.5), bits(1), bits(-0.5)},
		{bits(2.5), bits(1), bits(0.5)},
		{bits(3.5), bits(1), bits(-0.5)},
		{bits(-2.5), bits(-1), bits(-0.5)},
		{bits(-3.5), bits(1), bits(0.5)},
		{bits(-2), bits(2), sign},
		{sign, bits(2), sign},
		{4, 3, 1},
		{5, 3, sign | 1},
		{9, 6, sign | 3},
		{tiny, tiny * 2, tiny},
		{0x7fefffffffffffff, half, 0},
		{0x7fefffffffffffff, half + 1, 0xfca0000000000000},
		{bits(-7), inf, bits(-7)},
		{inf, inf, nan},
		{bits(1), sign, nan},
		{bits(1), nan, nan},
	}
	matches := func(got float64, want uint64) bool {
		if want == nan {
			return math.IsNaN(got)
		}
		return bits(got) == want
	}
	var truncationFailures, awayFailures int
	for _, c := range cases {
		x, y := math.Float64frombits(c.x), math.Float64frombits(c.y)
		got := math.Remainder(x, y)
		if !matches(got, c.want) {
			log.Fatalf("Remainder(%#016x, %#016x) = %#016x, want %#016x", c.x, c.y, bits(got), c.want)
		}
		if !matches(math.Mod(x, y), c.want) {
			truncationFailures++
		}
		away := got
		if halfY := math.Abs(got); halfY != 0 && !math.IsInf(y, 0) && halfY+halfY == math.Abs(y) {
			away = -math.Copysign(halfY, x)
		}
		if !matches(away, c.want) {
			awayFailures++
		}
	}
	if truncationFailures != 6 || awayFailures != 3 {
		log.Fatalf("truncating quotient failed %d, ties away from zero failed %d", truncationFailures, awayFailures)
	}
	fmt.Printf("Remainder passed %d fixtures; truncating quotient failed %d; ties away from zero failed %d\n",
		len(cases), truncationFailures, awayFailures)
}
