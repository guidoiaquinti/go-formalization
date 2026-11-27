package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const sign = uint64(1) << 63
	const inf = uint64(0x7ff0000000000000)
	const nan = inf + 1
	bits := math.Float64bits
	cases := []struct{ x, y, z, want uint64 }{
		{bits(1 + 0x1p-27), bits(1 - 0x1p-27), bits(-1), bits(-0x1p-54)},
		{bits(math.MaxFloat64), bits(2), bits(-math.MaxFloat64), bits(math.MaxFloat64)},
		{bits(1), bits(1), bits(0x1p-53), bits(1)},
		{bits(1), bits(1) + 1, bits(0x1p-53), bits(1) + 2},
		{bits(1), bits(1), bits(0x1p-53) + 1, bits(1) + 1},
		{bits(3), bits(2), bits(-6), 0},
		{sign | 1, bits(0.5), 0, sign},
		{sign, bits(1), sign, sign},
		{sign, bits(1), 0, 0},
		{inf, 0, bits(1), nan},
		{inf, bits(2), sign | inf, nan},
		{bits(2), bits(3), sign | inf, sign | inf},
		{bits(2), nan, bits(3), nan},
	}
	observe := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
	}
	var input bytes.Buffer
	var separateFailures int
	for _, c := range cases {
		fmt.Fprintln(&input, c.x, c.y, c.z)
		x, y, z := math.Float64frombits(c.x), math.Float64frombits(c.y), math.Float64frombits(c.z)
		want := observe(c.want)
		if got := observe(bits(math.FMA(x, y, z))); got != want {
			log.Fatalf("FMA(%#x, %#x, %#x)=%s want %s", c.x, c.y, c.z, got, want)
		}
		if observe(bits(float64(x*y)+z)) != want {
			separateFailures++
		}
	}
	for i, row := range check.LeanRows("godiff-math-fma", len(cases), &input) {
		if want := observe(cases[i].want); row != want {
			log.Fatalf("case %d: Lean=%s want %s", i, row, want)
		}
	}
	if separateFailures != 3 {
		log.Fatalf("separate multiply and add failed %d cases, want 3", separateFailures)
	}
	fmt.Printf("FMA passed %d fixtures; separately rounded multiply and add failed %d\n", len(cases), separateFailures)
}
