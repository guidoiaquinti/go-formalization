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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []struct {
		length int
		want   bool
	}{
		{math.MinInt, false}, {math.MinInt + 1, false},
		{-256, false}, {-9, false}, {-8, false}, {-1, false},
		{0, false}, {1, false}, {7, false}, {8, false},
		{9, true}, {10, true}, {16, true},
		{math.MaxInt - 1, true}, {math.MaxInt, true},
	}
	var input, data bytes.Buffer
	var inclusiveEight, strictNine, unsigned int
	for _, c := range cases {
		fmt.Fprintln(&input, c.length)
		fmt.Fprintf(&data, "%d\t%t\n", c.length, c.want)
		if (c.length >= 8) != c.want {
			inclusiveEight++
		}
		if (c.length > 9) != c.want {
			strictNine++
		}
		if (uint(c.length) > 8) != c.want {
			unsigned++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-should-use-ascii-set", len(cases), &input) {
		if want := strconv.FormatBool(cases[i].want); row != want {
			return fmt.Errorf("length %d: Lean=%s want=%s", cases[i].length, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-bytes-ascii-threshold-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_SHOULD_USE_ASCII_SET_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/bytes/should-use-ascii-set/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanShouldUseASCIISet",
		fixtures+"bridge_test.go", fixtures+"should_use_ascii_set_test.go"); err != nil {
		return err
	}
	if inclusiveEight != 1 || strictNine != 1 || unsigned != 6 {
		return fmt.Errorf("inclusive eight failed %d, strict nine failed %d, unsigned failed %d", inclusiveEight, strictNine, unsigned)
	}
	fmt.Printf("shouldUseASCIISet passed %d fixtures; inclusive eight failed %d, strict nine failed %d, unsigned failed %d\n",
		len(cases), inclusiveEight, strictNine, unsigned)
	return nil
}
