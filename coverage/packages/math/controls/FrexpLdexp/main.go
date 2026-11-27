package main

import (
	"fmt"
	"log"
	"math"
)

func main() {
	frexp()
	ldexp()
}

func frexp() {
	cases := []struct {
		input, fraction uint64
		exponent        int
	}{
		{0, 0, 0},
		{0x8000000000000000, 0x8000000000000000, 0},
		{0x0000000000000001, 0x3fe0000000000000, -1073},
		{0x8000000000000003, 0xbfe8000000000000, -1072},
		{0x000fffffffffffff, 0x3feffffffffffffe, -1022},
		{0x0010000000000000, 0x3fe0000000000000, -1021},
		{0x4008000000000000, 0x3fe8000000000000, 2},
		{0x7fefffffffffffff, 0x3fefffffffffffff, 1024},
		{0xfff0000000000000, 0xfff0000000000000, 0},
		{0xfff0000000000001, 0xfff0000000000001, 0},
	}
	for _, candidate := range []struct {
		name string
		run  func(float64) (float64, int)
		want int
	}{
		{"Frexp", math.Frexp, 0},
		{"Frexp skips subnormal normalization", skipNormalization, 3},
		{"Frexp omits exponent increment", omitIncrement, 6},
	} {
		failures := 0
		for _, test := range cases {
			f, e := candidate.run(math.Float64frombits(test.input))
			if math.Float64bits(f) != test.fraction || e != test.exponent {
				failures++
			}
		}
		assertCount(candidate.name, failures, candidate.want, len(cases))
	}
}

func skipNormalization(x float64) (float64, int) {
	word := math.Float64bits(x)
	if mag := word & 0x7fffffffffffffff; mag != 0 && mag < 1<<52 {
		return math.Float64frombits(word | 0x3fe0000000000000), -1022
	}
	return math.Frexp(x)
}

func omitIncrement(x float64) (float64, int) {
	f, e := math.Frexp(x)
	if x != 0 && !math.IsNaN(x) && !math.IsInf(x, 0) {
		e--
	}
	return f, e
}

func ldexp() {
	cases := []struct {
		input    uint64
		exponent int
		want     uint64
	}{
		{0x4000000000000000, math.MaxInt, 0},
		{0x3fe0000000000000, math.MinInt, 0x7ff0000000000000},
		{0xbfe0000000000000, math.MinInt, 0xfff0000000000000},
		{0x3ff0000000000000, -1075, 0},
		{0xbff0000000000000, -1075, 0x8000000000000000},
		{0x3ff0000000000000, -1074, 1},
		{0x3ff8000000000000, -1074, 2},
		{0x3ff0000000000001, -1075, 1},
		{0x0010000000000000, -1, 0x0008000000000000},
		{0x0000000000000001, 1074, 0x3ff0000000000000},
		{0x8000000000000000, math.MaxInt, 0x8000000000000000},
		{0xfff0000000000001, 22, 0xfff0000000000001},
	}
	for _, candidate := range []struct {
		name string
		run  func(float64, int) float64
		want int
	}{
		{"Ldexp", math.Ldexp, 0},
		{"Ldexp clamps before wrapped addition", clampEarly, 3},
		{"Ldexp flushes subnormals", flushSubnormals, 4},
	} {
		failures := 0
		for _, test := range cases {
			if got := math.Float64bits(candidate.run(math.Float64frombits(test.input), test.exponent)); got != test.want {
				failures++
				if candidate.want == 0 {
					log.Fatalf("%s(%#016x, %d) = %#016x, want %#016x", candidate.name, test.input, test.exponent, got, test.want)
				}
			}
		}
		assertCount(candidate.name, failures, candidate.want, len(cases))
	}
}

func clampEarly(x float64, e int) float64 {
	if x != 0 && !math.IsNaN(x) && !math.IsInf(x, 0) {
		switch {
		case e > 4096:
			return math.Copysign(math.Inf(1), x)
		case e < -4096:
			return math.Copysign(0, x)
		}
	}
	return math.Ldexp(x, e)
}

func flushSubnormals(x float64, e int) float64 {
	result := math.Ldexp(x, e)
	if magnitude := math.Abs(result); magnitude > 0 && magnitude < 0x1p-1022 {
		return math.Copysign(0, result)
	}
	return result
}

func assertCount(name string, got, want, total int) {
	if got != want {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, got, want)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, got, total)
}
