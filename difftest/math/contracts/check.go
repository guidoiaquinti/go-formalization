// Command contracts checks selected public math contracts directly against Go.
package main

import (
	"flag"
	"fmt"
	"math"
	"math/big"
	"os"
	"runtime"
	"strconv"
)

type observation struct {
	input string
	want  float64
	got   float64
	// A nonempty witness names the pinned result to replay for a known discrepancy.
	witness string
}

func main() {
	replay := flag.Bool("replay", false, "require the recorded pinned witnesses and contract agreement everywhere else")
	flag.Parse()
	fmt.Printf("%s %s/%s, int=%d\n", runtime.Version(), runtime.GOOS, runtime.GOARCH, strconv.IntSize)
	failed := false
	for _, group := range []struct {
		name string
		rows []observation
	}{
		{"Erfcinv", erfcinvObservations()},
		{"Ldexp", ldexpObservations()},
	} {
		disagreements := 0
		for _, row := range group.rows {
			want, got := result(row.want), result(row.got)
			if got != want {
				disagreements++
				fmt.Printf("%s(%s): contract=%s Go=%s\n", group.name, row.input, want, got)
			}
			expected := want
			if *replay && row.witness != "" {
				expected = row.witness
			}
			if got != expected {
				failed = true
				if *replay {
					fmt.Printf("%s(%s): replay expected %s, got %s\n", group.name, row.input, expected, got)
				}
			}
		}
		fmt.Printf("%s: %d cases, %d contract disagreements\n", group.name, len(group.rows), disagreements)
	}
	if failed {
		os.Exit(1)
	}
}

func erfcinvObservations() []observation {
	// The public special cases provide the oracle for these inputs.
	nan, positiveInfinity, negativeInfinity := math.NaN(), math.Inf(1), math.Inf(-1)
	cases := []struct {
		x, want float64
		witness string
	}{
		{negativeInfinity, nan, ""},
		{-1, nan, ""},
		{math.Nextafter(-0x1p-53, negativeInfinity), nan, ""},
		{-0x1p-53, nan, result(positiveInfinity)},
		{-math.SmallestNonzeroFloat64, nan, result(positiveInfinity)},
		{math.Copysign(0, -1), positiveInfinity, ""},
		{0, positiveInfinity, ""},
		{2, negativeInfinity, ""},
		{math.Nextafter(2, positiveInfinity), nan, ""},
		{positiveInfinity, nan, ""},
		{nan, nan, ""},
	}
	var rows []observation
	for _, c := range cases {
		rows = append(rows, observation{result(c.x), c.want, math.Erfcinv(c.x), c.witness})
	}
	return rows
}

func ldexpObservations() []observation {
	// Adjacent native exponents distinguish wraparound from ordinary scaling.
	var rows []observation
	for _, frac := range []float64{0.5, -0.5, 2, -2} {
		for _, exp := range []int{
			math.MinInt, math.MinInt + 1, -1076, -1075, -1074, -1073,
			-1022, 0, 1022, 1023, 1024, math.MaxInt - 1, math.MaxInt,
		} {
			witness := ""
			switch {
			case math.Abs(frac) == 0.5 && exp == math.MinInt:
				witness = result(math.Copysign(math.Inf(1), frac))
			case math.Abs(frac) == 2 && exp == math.MaxInt:
				witness = result(math.Copysign(0, frac))
			}
			rows = append(rows, observation{
				input:   fmt.Sprintf("%s, %d", result(frac), exp),
				want:    scaled(frac, exp),
				got:     math.Ldexp(frac, exp),
				witness: witness,
			})
		}
	}
	return rows
}

// scaled evaluates frac × 2**exp for finite nonzero binary64 frac, then
// rounds to binary64. Bound exp before adding the mantissa exponent so native
// integer wraparound cannot affect the reference.
func scaled(frac float64, exp int) float64 {
	if exp > 2098 {
		return math.Copysign(math.Inf(1), frac)
	}
	if exp < -2098 {
		return math.Copysign(0, frac)
	}
	value := new(big.Float).SetFloat64(frac)
	leading := value.MantExp(value)
	value.SetMantExp(value, leading+exp)
	rounded, _ := value.Float64()
	return rounded
}

func result(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return fmt.Sprintf("%#016x", math.Float64bits(x))
}
