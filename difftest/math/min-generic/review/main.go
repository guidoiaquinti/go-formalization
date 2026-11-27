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
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	const positiveInf, negativeInf = uint64(0x7ff0000000000000), uint64(0xfff0000000000000)
	const negativeZero, canonicalNaN = uint64(1 << 63), uint64(0x7ff8000000000001)
	const nan = uint64(0xfff8000000000042)
	cases := []struct{ x, y, want uint64 }{
		{negativeInf, nan, negativeInf},
		{nan, negativeInf, negativeInf},
		{nan, positiveInf, canonicalNaN},
		{0x7ff0000000000001, 0, canonicalNaN},
		{negativeZero, 0, negativeZero},
		{0, negativeZero, negativeZero},
		{negativeZero, negativeZero, negativeZero},
		{0xc000000000000000, 0xbff0000000000000, 0xc000000000000000},
	}
	var rows bytes.Buffer
	var nanPriorityFailures, leftTieFailures int
	for _, c := range cases {
		fmt.Fprintln(&rows, c.x, c.y, c.want)
		x, y := math.Float64frombits(c.x), math.Float64frombits(c.y)
		if (math.IsNaN(x) || math.IsNaN(y)) && canonicalNaN != c.want {
			nanPriorityFailures++
		}
		if x == y && c.x != c.want {
			leftTieFailures++
		}
	}
	if nanPriorityFailures != 2 || leftTieFailures != 1 {
		log.Fatalf("mutant failures NaN first=%d left tie=%d, want two and one", nanPriorityFailures, leftTieFailures)
	}

	dir, err := os.MkdirTemp("", "golean-min-generic-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_MIN_GENERIC_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/min-generic/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanMinGeneric", fixtures+"bridge_test.go", fixtures+"min_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s/%s: min passed 8 fixtures; NaN-first and left-tie mutants rejected twice and once\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
