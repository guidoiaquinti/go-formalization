package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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
		high, low         uint64
		shift             uint
		wantHigh, wantLow uint64
	}{
		{1, 3, 0, 1, 3},
		{3, 1<<63 | 1, 1, 7, 2},
		{1, 3, 63, 1<<63 | 1, 1 << 63},
		{7, 3, 64, 3, 0},
		{7, 3, 65, 6, 0},
		{7, 3, 127, 1 << 63, 0},
		{1, 3, 128, 0, 0},
		{1, 3, ^uint(0), 0, 0},
	}
	var rows bytes.Buffer
	var independentFailures, maskedFailures int
	for _, c := range cases {
		fmt.Fprintln(&rows, c.high, c.low, c.shift, c.wantHigh, c.wantLow)
		if c.high<<c.shift != c.wantHigh || c.low<<c.shift != c.wantLow {
			independentFailures++
		}
		n := c.shift & 127
		maskedHigh := c.high<<n | c.low>>(64-n) | c.low<<(n-64)
		maskedLow := c.low << n
		if maskedHigh != c.wantHigh || maskedLow != c.wantLow {
			maskedFailures++
		}
	}
	if independentFailures != 5 || maskedFailures != 2 {
		return fmt.Errorf("mutant failures: independent=%d, masked=%d; want 5 and 2", independentFailures, maskedFailures)
	}
	dir, err := os.MkdirTemp("", "golean-shl-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_SHL_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/shl/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanShl", fixtures+"bridge_test.go", fixtures+"shl_test.go"); err != nil {
		return err
	}
	fmt.Printf("%s %s/%s (uint%d): shl passed %d cases, independent and masked shifts failed %d and %d\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, strconv.IntSize, len(cases), independentFailures, maskedFailures)
	return nil
}
