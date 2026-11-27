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
	cases := []struct {
		x, want uint64
	}{
		{0, 0},
		{1 << 63, 1 << 63},
		{0x8000000000000001, 1 << 63},
		{0xbfe8000000000000, 1 << 63},
		{0xbffc000000000000, 0xbff0000000000000},
		{0x3ffc000000000000, 0x3ff0000000000000},
		{0x4330000000000001, 0x4330000000000001},
		{0xfff0000000000000, 0xfff0000000000000},
		{0x7ff0000000000001, 0x7ff0000000000001},
		{0xfff80000000000a5, 0xfff80000000000a5},
	}
	var rows bytes.Buffer
	var floorFailures, canonicalNaNFailures int
	for _, c := range cases {
		fmt.Fprintln(&rows, c.x, c.want)
		input := math.Float64frombits(c.x)
		if !math.IsNaN(input) && math.Float64bits(math.Floor(input)) != c.want {
			floorFailures++
		}
		if math.IsNaN(input) && math.Float64bits(math.NaN()) != c.want {
			canonicalNaNFailures++
		}
	}
	if floorFailures != 3 || canonicalNaNFailures != 2 {
		log.Fatalf("mutant failures floor=%d canonical NaN=%d, want three and two", floorFailures, canonicalNaNFailures)
	}

	dir, err := os.MkdirTemp("", "golean-trunc-generic-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_TRUNC_GENERIC_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/trunc-generic/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanTruncGeneric", fixtures+"bridge_test.go", fixtures+"trunc_generic_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s/%s: trunc passed 10 fixtures, floor and canonical-NaN mutants rejected three and two times\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
