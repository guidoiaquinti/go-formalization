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
	x, want float64
	exact   bool
}

func agree(result float64, c fixture) bool {
	if c.exact {
		return math.Float64bits(result) == math.Float64bits(c.want)
	}
	return math.Abs(result-c.want) <= 0x1p-51
}

func main() {
	// Golden values are rounded from an independent 110-digit decimal evaluation.
	fixtures := []fixture{
		{0.03125, 0x1.ffd559992b1dep-6, false}, {0.0625, 0x1.ff55997e030d7p-5, false},
		{0.3125, 0x1.35f98a0ea650ep-2, false}, {0.625, 0x1.1bf47eabb8f95p-1, false},
		{-0.4375, -0x1.a5729ee488037p-2, false}, {-0.875, -0x1.686650b8c2015p-1, false},
		{0.75, 0x1.45323e552f228p-1, false}, {1.5, 0x1.cf6f9786df577p-1, false},
		{2.5, 0x1.f9258260a71c2p-1, false}, {5, 0x1.fff419668df11p-1, false},
		{-9.25, -0x1.ffffff614d446p-1, false}, {-18.5, -0x1.ffffffffffffep-1, true},
		{19.5, 1, true}, {39, 1, true}, {math.Inf(-1), -1, true},
	}
	const profiles = 5
	var inputs, cases bytes.Buffer
	var linear, unit int
	for _, c := range fixtures {
		input := math.Float64bits(c.x)
		for profile := range profiles {
			fmt.Fprintln(&inputs, profile, input)
		}
		got := math.Tanh(c.x)
		if !agree(got, c) {
			log.Fatalf("Go Tanh(%g) = %.17g, want %.17g", c.x, got, c.want)
		}
		reflectionCheck := math.IsInf(c.x, 0)
		fmt.Fprintln(&cases, input, math.Float64bits(got), math.Float64bits(c.want), reflectionCheck)
		if math.Abs(c.x) < 0.625 && !agree(c.x, c) {
			linear++
		}
		if math.Abs(c.x) >= 0.625 && !agree(math.Copysign(1, c.x), c) {
			unit++
		}
	}
	lean := check.LeanRows("godiff-math-tanh", len(fixtures)*profiles, &inputs)
	results := make([]float64, len(lean))
	for i, row := range lean {
		word, err := strconv.ParseUint(strings.TrimSpace(row), 10, 64)
		must(err)
		results[i] = math.Float64frombits(word)
		c := fixtures[i/profiles]
		if !agree(results[i], c) {
			log.Fatalf("Lean profile %d Tanh(%g) = %.17g, want %.17g", i%profiles, c.x, results[i], c.want)
		}
	}
	for i := 0; i < len(fixtures)-1; i += 2 {
		for profile := range profiles {
			x := results[i*profiles+profile]
			doubled := results[(i+1)*profiles+profile]
			if math.Abs(doubled-2*x/(1+x*x)) > 0x1p-51 {
				log.Fatalf("Lean profile %d failed double-angle identity at %g", profile, fixtures[i].x)
			}
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-tanh-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_TANH_CASES", path))
	const overlay = "difftest/math/tanh/testdata/"
	must(check.TestInPackage("math", "TestGoLeanTanh", overlay+"bridge_test.go", overlay+"tanh_test.go"))
	if linear == 0 || unit == 0 {
		log.Fatalf("expected to detect both mutants, linear=%d premature saturation=%d", linear, unit)
	}
	fmt.Printf("Private tanh matched %d fixtures and %d Lean profile fixtures; linear failed %d and premature saturation failed %d\n",
		len(fixtures), len(lean), linear, unit)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
