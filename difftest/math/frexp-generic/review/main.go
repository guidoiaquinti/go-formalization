package main

import (
	"bytes"
	"fmt"
	"log"
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
	const fractionBase = uint64(1022) << 52
	const fractionMask = uint64(1<<52 - 1)
	const sign = uint64(1) << 63
	cases := []struct {
		x, fraction uint64
		exponent    int
	}{
		{0, 0, 0},
		{sign, sign, 0},
		{1, fractionBase, -1073},
		{sign | 3, sign | fractionBase | 1<<51, -1072},
		{fractionMask, fractionBase | (fractionMask - 1), -1022},
		{1 << 52, fractionBase, -1021},
		{0x3ff8000000000000, fractionBase | 1<<51, 1},
		{0xfff0000000000000, 0xfff0000000000000, 0},
		{0xfff8000000000042, 0xfff8000000000042, 0},
		{0x7ff0000000000001, 0x7ff0000000000001, 0},
	}
	var rows bytes.Buffer
	var unnormalizedFailures int
	for _, c := range cases {
		fmt.Fprintln(&rows, c.x, c.fraction, c.exponent)
		if f := c.x & fractionMask; f != 0 && (c.x>>52)&2047 == 0 {
			if (c.x&sign|fractionBase|f) != c.fraction || c.exponent != -1022 {
				unnormalizedFailures++
			}
		}
	}
	if unnormalizedFailures != 3 {
		log.Fatalf("skip-normalization mutant failed %d cases, want three", unnormalizedFailures)
	}

	dir, err := os.MkdirTemp("", "golean-frexp-generic-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_FREXP_GENERIC_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/frexp-generic/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanFrexpGeneric", fixtures+"bridge_test.go", fixtures+"frexp_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s/%s: frexp passed 10 fixtures; skip-normalization mutant rejected three times\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
