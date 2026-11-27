package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	type fixture struct {
		x, y float64
		word uint64
	}
	nan := math.Float64frombits(0x7ff8000000000042)
	sub := math.SmallestNonzeroFloat64
	cases := []fixture{
		{math.Copysign(0, -1), 0, 0}, {-3, 4, math.Float64bits(5)},
		{5, -12, math.Float64bits(13)}, {3 * sub, 4 * sub, 5},
		{math.MaxFloat64, 0, math.Float64bits(math.MaxFloat64)},
		{math.MaxFloat64, math.MaxFloat64, 0x7ff0000000000000},
		{nan, 1, 0x7ff8000000000001}, {1, nan, 0x7ff8000000000001},
		{math.Inf(-1), nan, 0x7ff0000000000000}, {nan, math.Inf(1), 0x7ff0000000000000},
		{math.Float64frombits(0xfff0000000000042), 1, 0x7ff8000000000001},
		{1, sub, math.Float64bits(1)},
	}
	var input bytes.Buffer
	for p := range 2 {
		for _, c := range cases {
			fmt.Fprintln(&input, p, math.Float64bits(c.x), math.Float64bits(c.y))
		}
	}
	rows := check.LeanRows("godiff-math-arch-hypot", 2*len(cases), &input)
	var native bytes.Buffer
	for i, c := range cases {
		want := strconv.FormatUint(c.word, 10)
		fallback := rows[i+len(cases)]
		if rows[i] != want || fallback != "panic:not implemented" {
			log.Fatalf("case %d: assembly %s, fallback %s, expected assembly %s", i, rows[i], fallback, want)
		}
		if runtime.GOARCH != "amd64" {
			want = fallback
		}
		fmt.Fprintf(&native, "%d %d %q\n", math.Float64bits(c.x), math.Float64bits(c.y), want)
	}
	if runtime.GOARCH == "386" {
		return
	}
	file, err := os.CreateTemp("", "golean-arch-hypot-review-*")
	must(err)
	defer os.Remove(file.Name())
	_, err = file.Write(native.Bytes())
	must(err)
	must(file.Close())
	must(os.Setenv("GOLEAN_MATH_ARCH_HYPOT_CASES", file.Name()))
	const data = "difftest/math/arch-hypot/testdata/"
	must(check.TestInPackage("math", "TestGoLeanArchHypot", data+"bridge_test.go", data+"arch_hypot_test.go"))
	fmt.Printf("archHypot: %d native and %d Lean observations passed.\n", len(cases), len(rows))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
