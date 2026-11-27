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
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Word  uint64
	First float64
	Gamma float64
	Mode  string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	const threshold = 143.01608
	inputs := []float64{33, 34, 100, 143, math.Nextafter(threshold, 0), threshold,
		math.Nextafter(threshold, math.Inf(1)), 144, 170, 200, math.Nextafter(200, math.Inf(1)), 201, math.Inf(1)}
	cases := make([]fixture, 0, len(inputs))
	for _, x := range inputs {
		c := fixture{Word: math.Float64bits(x), Mode: "saturated"}
		if x <= 200 {
			logX := math.Log(x)
			c.Mode = "ordinary"
			c.First = math.Exp((x-0.5)*logX - x)
			if x > threshold {
				c.Mode = "split"
				c.First = math.Exp((0.5*x - 0.25) * logX)
			}
			if x <= 170 && math.Trunc(x) == x {
				c.Gamma = 1
				for k := 2; k < int(x); k++ {
					c.Gamma *= float64(k)
				}
			}
		}
		cases = append(cases, c)
	}
	var rows bytes.Buffer
	for profile := range 5 {
		for _, c := range cases {
			fmt.Fprintln(&rows, profile, c.Word)
		}
	}
	results := check.LeanRows("godiff-math-stirling", 5*len(cases), &rows)
	collapse, early := 0, 0
	for i, row := range results {
		a, b, ok := strings.Cut(row, " ")
		if !ok {
			return fmt.Errorf("invalid Lean row %q", row)
		}
		first, err := strconv.ParseUint(a, 10, 64)
		if err != nil {
			return err
		}
		second, err := strconv.ParseUint(b, 10, 64)
		if err != nil {
			return err
		}
		c := cases[i%len(cases)]
		f, s := math.Float64frombits(first), math.Float64frombits(second)
		if !matches(c, f, s) {
			return fmt.Errorf("Lean profile %d: %#x gave %g, %g, expected %s", i/len(cases), c.Word, f, s, c.Mode)
		}
		if !matches(c, f*s, 1) {
			collapse++
		}
		if math.Float64frombits(c.Word) > 170 && !matches(c, math.Inf(1), 1) {
			early++
		}
	}
	if collapse != 50 || early != 5 {
		return fmt.Errorf("mutation failures collapse=%d early=%d, want 50 and 5", collapse, early)
	}

	dir, err := os.MkdirTemp("", "golean-review-stirling-")
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
	if err := os.Setenv("GOLEAN_REVIEW_STIRLING_CASES", path); err != nil {
		return err
	}
	const base = "difftest/math/stirling/"
	if err := check.TestInPackage("math", "TestGoLeanReviewStirling",
		base+"testdata/bridge_test.go", base+"review/testdata/review_test.go"); err != nil {
		return err
	}
	fmt.Printf("Stirling passed %d native and %d Lean profile rows; collapsed-pair and early-cutoff mutants failed %d and %d.\n",
		len(cases), len(results), collapse, early)
	return nil
}

func matches(c fixture, first, second float64) bool {
	if c.Mode == "saturated" {
		return math.IsInf(first, 1) && second == 1
	}
	if !(math.Abs(first/c.First-1) <= 0x1p-40) {
		return false
	}
	if c.Mode == "ordinary" && !(2.5 < second && second < 2.6) {
		return false
	}
	if c.Mode == "split" && !(second > 1e50 && !math.IsInf(second, 0)) {
		return false
	}
	return c.Gamma == 0 || math.Abs(first*second/c.Gamma-1) <= 0x1p-40
}
