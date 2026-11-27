package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	y, x, want float64
	kind       string
}

func agrees(got float64, test fixture) bool {
	if math.IsNaN(test.want) {
		return math.IsNaN(got)
	}
	if test.kind != "numeric" {
		return math.Float64bits(got) == math.Float64bits(test.want)
	}
	return math.Abs(got-test.want) <= 0x1p-50
}

func main() {
	negZero := math.Copysign(0, -1)
	inf := math.Inf(1)
	tiny := math.SmallestNonzeroFloat64
	tests := []fixture{
		{0, -1, math.Pi, "exact"}, {negZero, 1, negZero, "exact"}, {negZero, -1, -math.Pi, "exact"},
		{inf, 1, math.Pi / 2, "exact"}, {-inf, -1, -math.Pi / 2, "exact"},
		{1, -inf, math.Pi, "exact"}, {-1, inf, negZero, "exact"},
		{inf, -inf, 3 * math.Pi / 4, "exact"}, {-inf, inf, -math.Pi / 4, "exact"},
		{math.NaN(), inf, math.NaN(), "exact"},
		{2, 4, 0x1.dac670561bb50p-2, "numeric"},
		{-2, 4, -0x1.dac670561bb50p-2, "numeric"},
		{2, -4, math.Pi - 0x1.dac670561bb50p-2, "numeric"},
		{-2, -4, 0x1.dac670561bb50p-2 - math.Pi, "numeric"},
		{math.MaxFloat64, tiny, math.Pi / 2, "numeric"},
		{-math.MaxFloat64, -tiny, -math.Pi / 2, "numeric"},
		{-2 * tiny, -math.Nextafter(4, 0), -math.Pi, "numeric"},
		{-2 * tiny, -4, math.Pi, "underflow"},
		{-2 * tiny, -math.Nextafter(4, inf), math.Pi, "underflow"},
		{-2 * tiny, -inf, -math.Pi, "exact"},
	}
	const profiles = 2
	var inputs strings.Builder
	var cases bytes.Buffer
	var lostSign, mathematicalQuadrant int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&inputs, "%d %d %d\n", profile, math.Float64bits(test.y), math.Float64bits(test.x))
		}
		got := math.Atan2(test.y, test.x)
		if !agrees(got, test) {
			log.Fatalf("Go Atan2(%g, %g) = %.17g, want %.17g", test.y, test.x, got, test.want)
		}
		reference := test.want
		if test.kind == "underflow" {
			reference = -math.Pi
		}
		observation := "NaN"
		if !math.IsNaN(got) {
			observation = strconv.FormatUint(math.Float64bits(got), 10)
		}
		fmt.Fprintln(&cases, math.Float64bits(test.y), math.Float64bits(test.x),
			observation, math.Float64bits(reference), test.kind)
		if !agrees(math.Abs(got), test) {
			lostSign++
		}
		if test.y < 0 && test.x < 0 && test.y/test.x == 0 && !agrees(-math.Pi, test) {
			mathematicalQuadrant++
		}
	}
	rows := check.LeanRows("godiff-math-atan2", profiles*len(tests), strings.NewReader(inputs.String()))
	for i, row := range rows {
		got := math.NaN()
		if row != "NaN" {
			word, err := strconv.ParseUint(row, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(word)
		}
		test := tests[i/profiles]
		if !agrees(got, test) {
			log.Fatalf("Lean profile %d Atan2(%g, %g) = %.17g, want %.17g", i%profiles, test.y, test.x, got, test.want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-atan2-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_ATAN2_CASES", path))
	const overlay = "difftest/math/atan2/testdata/"
	must(check.TestInPackage("math", "TestGoLeanAtan2", overlay+"bridge_test.go", overlay+"atan2_test.go"))
	if lostSign == 0 || mathematicalQuadrant != 2 {
		log.Fatalf("expected sign and underflow differences, got %d and %d", lostSign, mathematicalQuadrant)
	}
	fmt.Printf("Private atan2 passed %d native and %d Lean fixtures; lost sign failed %d, substituting the mathematical quadrant failed %d\n",
		len(tests), len(rows), lostSign, mathematicalQuadrant)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
