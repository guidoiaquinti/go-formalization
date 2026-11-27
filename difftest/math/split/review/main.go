package main

import (
	"bytes"
	"fmt"
	"log"
	"math/bits"
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
	const implicit = uint64(1) << 52
	cases := []struct {
		word        uint64
		sign        uint32
		exponent    int32
		significand uint64
	}{
		{0, 0, -52, 0},
		{1 << 63, 1, -52, 0},
		{1, 0, -51, implicit},
		{1<<63 | implicit - 1, 1, 0, 2*implicit - 2},
		{implicit, 0, 1, implicit},
		{0xbff8000000000000, 1, 1023, implicit + implicit/2},
		{0x7ff0000000000000, 0, 2047, implicit},
		{0xfff8000000000001, 1, 2047, implicit + implicit/2 + 1},
	}
	var rows bytes.Buffer
	var shortShiftFailures, missingBitFailures int
	for _, c := range cases {
		fmt.Fprintln(&rows, c.word, c.sign, c.exponent, c.significand)
		fraction := c.word & (implicit - 1)
		if c.word>>52&2047 == 0 {
			shift := bits.LeadingZeros64(fraction) - 12
			if gotExponent, gotSignificand := 1-int32(shift), fraction<<shift; gotExponent != c.exponent || gotSignificand != c.significand {
				shortShiftFailures++
			}
		} else if fraction != c.significand {
			missingBitFailures++
		}
	}
	if shortShiftFailures != 4 || missingBitFailures != 4 {
		return fmt.Errorf("mutant failures: short shift=%d, missing implicit bit=%d, want four each", shortShiftFailures, missingBitFailures)
	}

	dir, err := os.MkdirTemp("", "golean-split-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_FMA_SPLIT_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/split/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanFmaSplit", fixtures+"bridge_test.go", fixtures+"split_test.go"); err != nil {
		return err
	}
	fmt.Printf("%s %s/%s: split passed %d cases, short-shift and missing-bit mutants failed four each\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases))
	return nil
}
