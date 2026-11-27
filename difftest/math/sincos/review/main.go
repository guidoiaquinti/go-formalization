package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	x, sine, cosine float64
	exact           bool
}

func agree(got, want float64, exact bool) bool {
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	if exact {
		return math.Float64bits(got) == math.Float64bits(want)
	}
	return math.Abs(got-want) <= 0x1p-50
}

func word(token string) float64 {
	if token == "NaN" {
		return math.NaN()
	}
	bits, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return math.Float64frombits(bits)
}

func main() {
	negativeZero := math.Copysign(0, -1)
	nan := math.NaN()
	d := math.Sqrt(0.5)
	fixtures := []fixture{
		{0, 0, 1, true}, {negativeZero, negativeZero, 1, true},
		{nan, nan, nan, true}, {math.Inf(1), nan, nan, true}, {math.Inf(-1), nan, nan, true},
		{math.Pi / 4, d, d, false}, {math.Pi / 2, 1, 0, false}, {3 * math.Pi / 4, d, -d, false},
		{math.Pi, 0, -1, false}, {5 * math.Pi / 4, -d, -d, false}, {3 * math.Pi / 2, -1, 0, false},
		{7 * math.Pi / 4, -d, d, false}, {-math.Pi / 2, -1, 0, false},
		// Platform libm values independently exercise large argument reduction.
		{0x1p29, 0.32656766301856333, -0.94517382606089662, false},
		{0x1p128, 0.97486851628605864, 0.22278100447349308, false},
		{0x1p512, 0.99057865836580183, 0.13694495825041525, false},
		{math.MaxFloat64, 0.004961954789184062, -0.99998768942655991, false},
	}
	var input strings.Builder
	var swapped, lostSigns int
	for _, c := range fixtures {
		for profile := range 3 {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(c.x))
		}
		sine, cosine := math.Sincos(c.x)
		if !agree(sine, c.sine, c.exact) || !agree(cosine, c.cosine, c.exact) {
			log.Fatalf("Go Sincos(%g) = (%g, %g), want (%g, %g)", c.x, sine, cosine, c.sine, c.cosine)
		}
		if !agree(cosine, c.sine, c.exact) || !agree(sine, c.cosine, c.exact) {
			swapped++
		}
		if !agree(math.Abs(sine), c.sine, c.exact) || !agree(math.Abs(cosine), c.cosine, c.exact) {
			lostSigns++
		}
	}
	rows := check.LeanRows("godiff-math-sincos", 3*len(fixtures), strings.NewReader(input.String()))
	for i, row := range rows {
		c := fixtures[i/3]
		sine, cosine, ok := strings.Cut(row, " ")
		if !ok || !agree(word(sine), c.sine, c.exact) || !agree(word(cosine), c.cosine, c.exact) {
			log.Fatalf("Lean profile %d Sincos(%g) = %s, want (%g, %g)", i%3, c.x, row, c.sine, c.cosine)
		}
	}
	if swapped == 0 || lostSigns == 0 {
		log.Fatalf("controls did not distinguish swapped outputs (%d) or lost signs (%d)", swapped, lostSigns)
	}
	fmt.Printf("Sincos passed %d native and %d Lean fixtures; swapping outputs failed %d, losing signs failed %d\n", len(fixtures), len(rows), swapped, lostSigns)
}
