package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type fixture struct {
	x, want, tolerance float64
}

func agrees(got float64, test fixture) bool {
	if math.IsNaN(test.want) {
		return math.IsNaN(got)
	}
	if test.want == 0 {
		return math.Float64bits(got) == math.Float64bits(test.want)
	}
	return math.Abs(got-test.want) <= test.tolerance*math.Abs(test.want)
}

func main() {
	negativeZero := math.Copysign(0, -1)
	nearPole := math.Nextafter(math.Pi/2, math.Inf(1))
	tests := []fixture{
		{0, 0, 0},
		{negativeZero, negativeZero, 0},
		{math.NaN(), math.NaN(), 0},
		{math.Inf(1), math.NaN(), 0},
		{math.Inf(-1), math.NaN(), 0},
		{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, 0},
		{math.Pi / 4, 1, 1e-14},
		{3 * math.Pi / 4, -1, 1e-14},
		{-math.Pi / 4, -1, 1e-14},
		{math.Pi / 2, 16331239353195370, 1e-13},
		{nearPole, -6218431163823738, 1e-13},
		{-nearPole, 6218431163823738, 1e-13},
	}
	var input strings.Builder
	for _, test := range tests {
		for profile := range 3 {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.x))
		}
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-math-tan")
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	results := bytes.Fields(output)
	if len(results) != 3*len(tests) {
		fmt.Fprintln(os.Stderr, "unexpected fixture result count:", len(results))
		os.Exit(1)
	}
	var modelFailures, failures, dropOddReciprocal, loseNegativeSign int
	for i, test := range tests {
		for profile := range 3 {
			got := math.NaN()
			token := string(results[3*i+profile])
			if token != "NaN" {
				word, err := strconv.ParseUint(token, 10, 64)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				got = math.Float64frombits(word)
			}
			if !agrees(got, test) {
				fmt.Fprintf(os.Stderr, "Lean profile %d Tan(%g) = %g, wanted %g\n", profile, test.x, got, test.want)
				modelFailures++
			}
		}
		if got := math.Tan(test.x); !agrees(got, test) {
			fmt.Fprintf(os.Stderr, "Tan(%g) = %g, wanted %g\n", test.x, got, test.want)
			failures++
		}
		reduced := math.Remainder(test.x, math.Pi/2)
		if !agrees(math.Tan(reduced), test) {
			dropOddReciprocal++
		}
		if !agrees(math.Tan(math.Abs(test.x)), test) {
			loseNegativeSign++
		}
	}
	fmt.Printf("Go disagreed on %d of %d fixtures\n", failures, len(tests))
	fmt.Printf("Substituting rounded remainder for the odd-quadrant reciprocal disagreed on %d of %d fixtures\n", dropOddReciprocal, len(tests))
	fmt.Printf("Losing the input sign disagreed on %d of %d fixtures\n", loseNegativeSign, len(tests))
	fmt.Printf("Lean disagreed on %d of %d profile fixtures\n", modelFailures, len(results))
	if modelFailures > 0 || failures > 0 || dropOddReciprocal == 0 || loseNegativeSign == 0 {
		os.Exit(1)
	}
}
