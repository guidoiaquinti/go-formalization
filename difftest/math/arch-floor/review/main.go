package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const sign = uint64(1) << 63
	const big = uint64(0x4330000000000000)
	fixtures := []struct {
		word uint64
		want string
	}{
		{0, "0"}, {sign, "9223372036854775808"},
		{1, "0"}, {sign | 1, "13830554455654793216"},
		{math.Float64bits(0.75), "0"}, {math.Float64bits(-0.75), "13830554455654793216"},
		{math.Float64bits(-3.25), strconv.FormatUint(math.Float64bits(-4), 10)},
		{math.Float64bits(-3), strconv.FormatUint(math.Float64bits(-3), 10)},
		{big - 1, strconv.FormatUint(big-2, 10)},
		{sign | (big - 1), strconv.FormatUint(sign|big, 10)},
		{big, strconv.FormatUint(big, 10)}, {sign | big, strconv.FormatUint(sign|big, 10)},
		{math.Float64bits(math.Inf(1)), strconv.FormatUint(math.Float64bits(math.Inf(1)), 10)},
		{math.Float64bits(math.Inf(-1)), strconv.FormatUint(math.Float64bits(math.Inf(-1)), 10)},
		{0x7ff0000000000001, "NaN"}, {0xfff80000000000a5, "NaN"},
	}
	var fixture, leanInput bytes.Buffer
	var truncateFailures, lowerGuardFailures int
	for _, c := range fixtures {
		fmt.Fprintln(&fixture, c.word, c.want)
		fmt.Fprintln(&leanInput, 0, c.word)
		fmt.Fprintln(&leanInput, 1, c.word)
		x := math.Float64frombits(c.word)
		if observe(math.Trunc(x)) != c.want {
			truncateFailures++
		}
		// Moving the unchanged-value guard down by one binary64 word
		// returns the fractional predecessor of 2^52 without rounding.
		lowerGuard := c.want
		if c.word&^sign == big-1 {
			lowerGuard = observe(x)
		}
		if lowerGuard != c.want {
			lowerGuardFailures++
		}
	}
	for i, row := range check.LeanRows("godiff-math-arch-floor", 2*len(fixtures), &leanInput) {
		if c := fixtures[i/2]; row != c.want {
			log.Fatalf("profile %d input %#016x: Lean=%s want=%s", i%2, c.word, row, c.want)
		}
	}
	if truncateFailures != 4 || lowerGuardFailures != 2 {
		log.Fatalf("truncation failures=%d, lowered guard failures=%d", truncateFailures, lowerGuardFailures)
	}
	dir, err := os.MkdirTemp("", "golean-arch-floor-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_ARCH_FLOOR_CASES", path); err != nil {
		log.Fatal(err)
	}
	const data = "difftest/math/arch-floor/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanArchFloor", data+"bridge_test.go", data+"arch_floor_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("archFloor passed %d fixtures under native Go and both Lean profiles; truncation failed %d; lowered guard failed %d\n",
		len(fixtures), truncateFailures, lowerGuardFailures)
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(math.Float64bits(x), 10)
}
