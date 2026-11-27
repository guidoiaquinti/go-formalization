package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type observation struct{ Input, Want, Steps uint64 }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	nz := math.Copysign(0, -1)
	quarter := math.Sqrt(0.5)
	samples := []struct {
		input, want float64
		steps       uint64
	}{
		{nz, 0, 0}, {0, nz, 0}, {math.SmallestNonzeroFloat64, -3 * math.SmallestNonzeroFloat64, 0},
		{0.25, -quarter, 2}, {0.5, -1, 0}, {0.75, -quarter, 2}, {1, nz, 0},
		{1.25, quarter, 2}, {1.5, 1, 0}, {1.75, quarter, 2}, {2, nz, 0},
		{math.Nextafter(0x1p52, 0), 1, 0}, {0x1p52, nz, 0}, {0x1p52 + 1, nz, 0},
		{0x1p53 - 1, nz, 0}, {0x1p53, nz, 0}, {math.MaxFloat64, nz, 0},
	}
	cases := make([]observation, len(samples))
	for i, c := range samples {
		cases[i] = observation{math.Float64bits(c.input), math.Float64bits(c.want), c.steps}
	}
	var rows bytes.Buffer
	for profile := range 3 {
		for _, c := range cases {
			fmt.Fprintln(&rows, profile, c.Input)
		}
	}
	results := check.LeanRows("godiff-math-sinpi", 3*len(cases), &rows)
	unreduced, unsigned := 0, 0
	for i, row := range results {
		c := cases[i%len(cases)]
		word, err := strconv.ParseUint(row, 10, 64)
		if err != nil {
			return err
		}
		if !matches(c, word) {
			return fmt.Errorf("Lean profile %d input %#016x: got %#016x, want %#016x", i/len(cases), c.Input, word, c.Want)
		}
		x := math.Float64frombits(c.Input)
		if !matches(c, math.Float64bits(-math.Sin(math.Pi*x))) {
			unreduced++
		}
		if !matches(c, math.Float64bits(math.Abs(math.Float64frombits(word)))) {
			unsigned++
		}
	}
	if unreduced == 0 || unsigned == 0 {
		return fmt.Errorf("expected mutant failures, got unreduced=%d unsigned=%d", unreduced, unsigned)
	}
	dir, err := os.MkdirTemp("", "golean-review-sinpi-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases.json")
	data, err := json.Marshal(cases)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_REVIEW_SINPI_CASES", path); err != nil {
		return err
	}
	const base = "difftest/math/sinpi/"
	if err := check.TestInPackage("math", "TestGoLeanReviewSinPi",
		base+"testdata/bridge_test.go", base+"review/testdata/review_test.go"); err != nil {
		return err
	}
	fmt.Printf("sinPi passed %d native and %d Lean profile rows; no reduction and lost sign failed %d and %d.\n",
		len(cases), len(results), unreduced, unsigned)
	return nil
}

func matches(c observation, word uint64) bool {
	return max(word, c.Want)-min(word, c.Want) <= c.Steps
}
