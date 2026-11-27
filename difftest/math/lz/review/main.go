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
	cases := []struct {
		high, low uint64
		want      int
	}{
		{0, 0, 128},
		{0, 1, 127},
		{0, 1 << 63, 64},
		{1, 0, 63},
		{1, 1 << 63, 63},
		{1 << 63, 1, 0},
		{0x100, 1, 55},
	}
	var rows bytes.Buffer
	var reversedFailures, missingOffsetFailures int
	for _, c := range cases {
		fmt.Fprintln(&rows, c.high, c.low, c.want)
		reversed := bits.LeadingZeros64(c.low)
		if c.low == 0 {
			reversed += bits.LeadingZeros64(c.high)
		}
		if reversed != c.want {
			reversedFailures++
		}
		missingOffset := bits.LeadingZeros64(c.high)
		if c.high == 0 {
			missingOffset = bits.LeadingZeros64(c.low)
		}
		if missingOffset != c.want {
			missingOffsetFailures++
		}
	}
	if reversedFailures != 6 || missingOffsetFailures != 3 {
		return fmt.Errorf("mutant failures: reversed=%d, missing offset=%d; want 6 and 3", reversedFailures, missingOffsetFailures)
	}
	dir, err := os.MkdirTemp("", "golean-lz-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_LZ_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/lz/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanLz", fixtures+"bridge_test.go", fixtures+"lz_test.go"); err != nil {
		return err
	}
	fmt.Printf("%s %s/%s: lz passed %d cases, reversed-word and missing-offset mutants failed %d and %d\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases), reversedFailures, missingOffsetFailures)
	return nil
}
