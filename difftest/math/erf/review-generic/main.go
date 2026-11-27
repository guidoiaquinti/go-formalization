package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	in, want uint64
	exact    bool
}

func observed(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(math.Float64bits(x), 10)
}

func agrees(c fixture, got float64) bool {
	want := math.Float64frombits(c.want)
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	if math.IsNaN(got) || math.Signbit(got) != math.Signbit(want) {
		return false
	}
	a, b := math.Float64bits(got), c.want
	if c.exact {
		return a == b
	}
	return max(a, b)-min(a, b) <= 4
}

func main() {
	// Finite, nonexact reference words come from the reviewer's 384-bit
	// positive-term series for exp(x²) integral_0^x exp(-t²)dt, followed
	// by division by exp(x²), normalization, and binary64 rounding.
	// Exact tiny fixtures separately preserve the source's scaled steps.
	const expProbe uint64 = 0x3ff6c2d5940ad8fb
	positive := []fixture{
		{0x0000000000000000, 0x0000000000000000, true},
		{0x0000000000000001, 0x0000000000000001, true},
		{0x0000000000000004, 0x0000000000000004, true},
		{0x0000000000000013, 0x0000000000000016, true},
		{0x007fffffffffffff, 0x00820dd750429b6d, false},
		{0x0080000000000000, 0x00820dd750429b6d, false},
		{0x0080000000000003, 0x00820dd750429b70, false},
		{0x3e2fffffffffffff, 0x3e320dd750429b6c, false},
		{0x3e30000000000000, 0x3e320dd750429b6d, false},
		{0x3feaffffffffffff, 0x3fe88d1cd474a2e0, false},
		{0x3feb000000000000, 0x3fe88d1cd474a2e0, false},
		{0x3ff3ffffffffffff, 0x3fed8865d98abe00, false},
		{0x3ff4000000000000, 0x3fed8865d98abe01, false},
		{expProbe, 0x3fee9596d5317465, false},
		{0x4006db6db6db6db6, 0x3fefff90322bb733, false},
		{0x4006db6db6db6db7, 0x3fefff90322bb733, false},
		{0x4008000000000000, 0x3fefffd1ac4135f9, false},
		{0x4016000000000000, 0x3fefffffffffffbe, false},
		{0x4017ffffffffffff, 0x3ff0000000000000, false},
		{0x4018000000000000, 0x3ff0000000000000, true},
		{0x7ff0000000000000, 0x3ff0000000000000, true},
		{0x7ff0000000000001, 0x7ff8000000000000, true},
	}
	const sign = 1 << 63
	const profiles = 5
	fixtures := make([]fixture, 0, 2*len(positive))
	var modelInput, nativeFixture bytes.Buffer
	for _, value := range positive {
		fixtures = append(fixtures, value, fixture{value.in | sign, value.want | sign, value.exact})
	}
	for _, c := range fixtures {
		for profile := range profiles {
			fmt.Fprintln(&modelInput, profile, c.in)
		}
	}
	rows := check.LeanRows("godiff-math-erf", profiles*len(fixtures), &modelInput)
	var noScaling, earlySaturation, complementCancellation int
	for i, c := range fixtures {
		x := math.Float64frombits(c.in)
		if c.in&^sign == expProbe && rows[i*profiles] == rows[i*profiles+1] {
			log.Fatalf("split and fused Exp profiles collapsed at %#016x", c.in)
		}
		for profile := range profiles {
			row := rows[i*profiles+profile]
			got := math.NaN()
			if row != "NaN" {
				word, err := strconv.ParseUint(row, 10, 64)
				must(err)
				got = math.Float64frombits(word)
			}
			if !agrees(c, got) {
				log.Fatalf("Lean private erf(%#016x), profile %d: got=%#016x want=%#016x",
					c.in, profile, math.Float64bits(got), c.want)
			}
		}
		native := observed(math.Erf(x))
		if runtime.GOARCH == "arm64" {
			native = rows[i*profiles+4]
		}
		kind := "numeric"
		if c.exact {
			kind = "exact"
		}
		fmt.Fprintln(&nativeFixture, c.in, native, observed(math.Float64frombits(c.want)), kind)

		magnitude := math.Abs(x)
		if c.exact && 0 < magnitude && magnitude < 0x1p-1015 {
			const efx = 0.128379167095512586316
			if !agrees(c, x+float64(efx*x)) {
				noScaling++
			}
		}
		if 1.25 <= magnitude && magnitude < 6 && !agrees(c, math.Copysign(1, x)) {
			earlySaturation++
		}
		if magnitude < 0x1p-28 && !agrees(c, 1-math.Erfc(x)) {
			complementCancellation++
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-erf-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, nativeFixture.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_ERF_CASES", path))
	const overlay = "difftest/math/erf/testdata/"
	must(check.TestInPackage("math", "TestGoLeanErf", overlay+"bridge_test.go", overlay+"erf_test.go"))
	if noScaling != 4 || earlySaturation != 12 || complementCancellation != 15 {
		log.Fatalf("unexpected discrimination: omitted scaling=%d, early saturation=%d, complement cancellation=%d",
			noScaling, earlySaturation, complementCancellation)
	}
	fmt.Printf("private erf: %d native and %d Lean profile fixtures; omitted scaling=%d, early saturation=%d, complement cancellation=%d\n",
		len(fixtures), len(rows), noScaling, earlySaturation, complementCancellation)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
