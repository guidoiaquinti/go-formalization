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
		{1, true},
		{-3, true},
		{1.5, false},
		{-3.5, false},
		{math.Nextafter(1<<52, 0), false},
		{1 << 52, false},
		{1<<52 + 1, true},
		{-(1<<53 - 1), true},
		{1 << 53, false},
		{math.Inf(1), false},
		{math.NaN(), false},
	}
	var rows bytes.Buffer
	var truncatedFailures, earlyCutoffFailures int
	for _, c := range cases {
		want := 0
		if c.want {
			want = 1
		}
		fmt.Fprintln(&rows, math.Float64bits(c.value), want)
		magnitude := math.Abs(c.value)
		whole, fraction := math.Modf(c.value)
		if got := magnitude < 1<<53 && int64(whole)&1 == 1; got != c.want {
			truncatedFailures++
		}
		if got := magnitude < 1<<52 && fraction == 0 && int64(whole)&1 == 1; got != c.want {
			earlyCutoffFailures++
		}
	}
	if truncatedFailures != 3 || earlyCutoffFailures != 2 {
		return fmt.Errorf("mutant failures: fraction ignored=%d, early cutoff=%d, want three and two", truncatedFailures, earlyCutoffFailures)
	}

	dir, err := os.MkdirTemp("", "golean-is-odd-int-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_IS_ODD_INT_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/is-odd-int/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanIsOddInt", fixtures+"bridge_test.go", fixtures+"is_odd_int_test.go"); err != nil {
		return err
	}
	fmt.Printf("%s %s/%s: isOddInt passed %d cases, fraction-ignored and early-cutoff mutants failed three and two\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases))
	return nil
}
