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
	const inf = uint64(0x7ff0000000000000)
	const nan = inf + 1
	bits := math.Float64bits
	cases := []struct{ x, y, want uint64 }{
		{bits(6.5), bits(3), bits(0.5)},
		{bits(7.5), bits(2), bits(1.5)},
		{bits(-7.5), bits(2), bits(-1.5)},
		{bits(3.25), bits(-2), bits(1.25)},
		{bits(-6), bits(3), sign},
		{sign, bits(2), sign},
		{sign | 7, 3, sign | 1},
		{0x0010000000000001, 3, 2},
		{0x7fefffffffffffff, 1, 0},
		{bits(-7), sign | inf, bits(-7)},
		{inf, inf, nan},
		{bits(1), sign, nan},
		{bits(1), nan, nan},
	}
	var input, fixture bytes.Buffer
	want := make([]string, len(cases))
	var nearestFailures, zeroSignFailures int
	for i, c := range cases {
		row := strconv.FormatUint(c.want, 10)
		if c.want == nan {
			row = "NaN"
		}
		want[i] = row
		fmt.Fprintln(&input, c.x, c.y)
		fmt.Fprintln(&fixture, c.x, c.y, row)
		x, y := math.Float64frombits(c.x), math.Float64frombits(c.y)
		if observe(math.Remainder(x, y)) != row {
			nearestFailures++
		}
		lostSign := math.Mod(x, y)
		if lostSign == 0 {
			lostSign = 0
		}
		if observe(lostSign) != row {
			zeroSignFailures++
		}
	}
	if nearestFailures != 4 || zeroSignFailures != 2 {
		log.Fatalf("nearest quotient failed %d, lost zero sign failed %d", nearestFailures, zeroSignFailures)
	}
	for i, row := range check.LeanRows("godiff-math-mod", len(cases), &input) {
		if row != want[i] {
			log.Fatalf("case %d: Lean=%q want=%q", i, row, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-mod-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_MOD_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/mod/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanMod", fixtures+"bridge_test.go", fixtures+"mod_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("private math.mod passed %d fixtures; nearest quotient failed %d; lost zero sign failed %d\n",
		len(cases), nearestFailures, zeroSignFailures)
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(math.Float64bits(x), 10)
}
