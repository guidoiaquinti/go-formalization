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
	const tiny = uint64(0x0020000000000000)
	const half = uint64(0x7fdfffffffffffff)
	bits := math.Float64bits
	cases := []struct{ x, y, want uint64 }{
		{bits(0.5), bits(1), bits(0.5)},
		{bits(1.5), bits(1), bits(-0.5)},
		{bits(2.5), bits(1), bits(0.5)},
		{bits(3.5), bits(1), bits(-0.5)},
		{bits(-2.5), bits(-1), bits(-0.5)},
		{bits(-3.5), bits(1), bits(0.5)},
		{bits(-2), bits(2), sign},
		{sign, bits(2), sign},
		{3, 2, sign | 1},
		{5, 2, 1},
		{4, 3, 1},
		{5, 3, sign | 1},
		{9, 6, sign | 3},
		{tiny, tiny * 2, tiny},
		{0x7fefffffffffffff, half, 0},
		{0x7fefffffffffffff, half + 1, 0xfca0000000000000},
		{bits(-7), inf, bits(-7)},
		{inf, inf, nan},
		{bits(1), sign, nan},
		{bits(1), nan, nan},
	}
	matches := func(got float64, want uint64) bool {
		if want == nan {
			return math.IsNaN(got)
		}
		return bits(got) == want
	}
	var truncationFailures, awayFailures int
	var inputs, fixture bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		want := strconv.FormatUint(c.want, 10)
		if c.want == nan {
			want = "NaN"
		}
		wants[i] = want
		fmt.Fprintln(&inputs, c.x, c.y)
		fmt.Fprintln(&fixture, c.x, c.y, want)
		x, y := math.Float64frombits(c.x), math.Float64frombits(c.y)
		got := math.Remainder(x, y)
		if !matches(got, c.want) {
			log.Fatalf("Remainder(%#016x, %#016x) = %#016x, want %#016x", c.x, c.y, bits(got), c.want)
		}
		if !matches(math.Mod(x, y), c.want) {
			truncationFailures++
		}
		away := got
		if halfY := math.Abs(got); halfY != 0 && !math.IsInf(y, 0) && halfY+halfY == math.Abs(y) {
			away = -math.Copysign(halfY, x)
		}
		if !matches(away, c.want) {
			awayFailures++
		}
	}
	if truncationFailures != 7 || awayFailures != 4 {
		log.Fatalf("truncating quotient failed %d, ties away from zero failed %d", truncationFailures, awayFailures)
	}
	for i, row := range check.LeanRows("godiff-math-remainder", len(cases), &inputs) {
		if row != wants[i] {
			log.Fatalf("case %d: Lean=%s want=%s", i, row, wants[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-remainder-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_REMAINDER_CASES", path); err != nil {
		log.Fatal(err)
	}
	const data = "difftest/math/remainder/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanRemainder", data+"bridge_test.go", data+"remainder_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Remainder passed %d fixtures; truncating quotient failed %d; ties away from zero failed %d\n",
		len(cases), truncationFailures, awayFailures)
}
