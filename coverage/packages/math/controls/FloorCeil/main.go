package main

import (
	"fmt"
	"log"
	"math"
)

type fixture struct {
	input, floor, ceil float64
}

func main() {
	signZero := math.Float64frombits(1 << 63)
	nan := math.Float64frombits(0xfff0000000000100)
	cases := []fixture{
		{0, 0, 0},
		{signZero, signZero, signZero},
		{math.SmallestNonzeroFloat64, 0, 1},
		{-math.SmallestNonzeroFloat64, -1, signZero},
		{0.75, 0, 1},
		{-0.75, -1, signZero},
		{1.75, 1, 2},
		{-1.75, -2, -1},
		{2, 2, 2},
		{-2, -2, -2},
		{0x1.fffffffffffffp+51, 0x1.ffffffffffffep+51, 0x1p52},
		{-0x1.fffffffffffffp+51, -0x1p52, -0x1.ffffffffffffep+51},
		{math.Inf(1), math.Inf(1), math.Inf(1)},
		{math.Inf(-1), math.Inf(-1), math.Inf(-1)},
		{nan, nan, nan},
	}
	count := func(round func(float64) float64, floor bool) int {
		mismatches := 0
		for _, test := range cases {
			want := test.ceil
			if floor {
				want = test.floor
			}
			got := round(test.input)
			if !(math.IsNaN(want) && math.IsNaN(got)) && math.Float64bits(got) != math.Float64bits(want) {
				mismatches++
			}
		}
		return mismatches
	}
	candidates := []struct {
		name  string
		run   func(float64) float64
		floor bool
		want  int
	}{
		{"Floor", math.Floor, true, 0},
		{"Floor replaced with truncation", math.Trunc, true, 4},
		{"Floor drops signed zero", func(x float64) float64 { return 0 + math.Floor(x) }, true, 1},
		{"Ceil", math.Ceil, false, 0},
		{"Ceil replaced with truncation", math.Trunc, false, 4},
		{"Ceil drops signed zero", func(x float64) float64 { return 0 + math.Ceil(x) }, false, 3},
	}
	for _, candidate := range candidates {
		if got := count(candidate.run, candidate.floor); got != candidate.want {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, got, candidate.want)
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, candidate.want, len(cases))
	}
}
