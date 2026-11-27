package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		return err
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		return fmt.Errorf("Go version %s, want %s", runtime.Version(), want)
	}
	cases := []struct {
		value float64
		want  bool
	}{
		{0, false},
		{math.Copysign(0, -1), false},
		{-math.SmallestNonzeroFloat64, false},
		{math.Nextafter(-1, 0), false},
		{-1, true},
		{-1.5, false},
		{-2, true},
		{2, false},
		{-math.MaxFloat64, true},
		{math.Inf(-1), false},
		{math.Float64frombits(0xfff8000000000001), false},
	}
	var rows bytes.Buffer
	var signBitFailures, noFractionFailures int
	for _, c := range cases {
		want := 0
		if c.want {
			want = 1
		}
		fmt.Fprintln(&rows, math.Float64bits(c.value), want)
		_, fraction := math.Modf(c.value)
		if got := math.Signbit(c.value) && fraction == 0; got != c.want {
			signBitFailures++
		}
		if got := c.value < 0 && !math.IsInf(c.value, -1); got != c.want {
			noFractionFailures++
		}
	}
	if signBitFailures != 1 || noFractionFailures != 3 {
		return fmt.Errorf("mutant failures: sign bit=%d, fraction ignored=%d, want one and three", signBitFailures, noFractionFailures)
	}

	dir, err := os.MkdirTemp("", "golean-is-neg-int-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_IS_NEG_INT_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/is-neg-int/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanIsNegInt", fixtures+"bridge_test.go", fixtures+"is_neg_int_test.go"); err != nil {
		return err
	}
	fmt.Printf("%s %s/%s: isNegInt passed %d cases, sign-bit and fraction-ignored mutants failed one and three\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases))
	return nil
}
