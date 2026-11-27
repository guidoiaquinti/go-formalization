package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	cases := []struct {
		x    float64
		want uint64
	}{
		{0, 0},
		{math.Copysign(0, -1), 1 << 63},
		{1, 0x3ff0000000000000},
		{2, 0x3ff6a09e667f3bcd},
		{4, 0x4000000000000000},
		{8, 0x4006a09e667f3bcd},
		{math.SmallestNonzeroFloat64, uint64(1023-537) << 52},
		{math.Inf(1), 0x7ff0000000000000},
	}
	var roundDownFailures, zeroFailures int
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&rows, math.Float64bits(c.x), c.want)
		if c.x > 0 && !math.IsInf(c.x, 1) {
			down := math.Float64frombits(c.want)
			exactRoot := new(big.Rat).SetFloat64(down)
			square := new(big.Rat).Mul(exactRoot, exactRoot)
			if square.Cmp(new(big.Rat).SetFloat64(c.x)) > 0 {
				down = math.Nextafter(down, 0)
			}
			if math.Float64bits(down) != c.want {
				roundDownFailures++
			}
		}
		if c.x == 0 && c.want != 0 {
			zeroFailures++
		}
	}
	for _, x := range []float64{-1, math.Inf(-1), math.Float64frombits(0xfff8000000000001)} {
		fmt.Fprintln(&rows, math.Float64bits(x), "NaN")
	}
	if roundDownFailures != 2 || zeroFailures != 1 {
		log.Fatalf("mutant failures down=%d zero=%d, want two and one", roundDownFailures, zeroFailures)
	}
	dir, err := os.MkdirTemp("", "golean-sqrt-generic-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_SQRT_GENERIC_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/sqrt-generic/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanSqrtGeneric", fixtures+"bridge_test.go", fixtures+"sqrt_generic_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s/%s: sqrt passed 11 fixtures, round-down and canonical-zero mutants rejected twice and once\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
