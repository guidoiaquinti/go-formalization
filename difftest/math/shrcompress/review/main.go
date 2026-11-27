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
		high, low        uint64
		shift            uint
		plainHigh        uint64
		plainLow         uint64
		discardedNonzero uint64
		wantLow          uint64
	}{
		{1, 3, 0, 1, 3, 0, 3},
		{0, 1, 1, 0, 0, 1, 1},
		{0, 3, 1, 0, 1, 1, 1},
		{0, 2, 1, 0, 1, 0, 1},
		{2, 1, 64, 0, 2, 1, 3},
		{3, 1, 64, 0, 3, 1, 3},
		{1, 0, 65, 0, 0, 1, 1},
		{2, 0, 65, 0, 1, 0, 1},
		{1, 0, 128, 0, 0, 1, 1},
		{0, 0, ^uint(0), 0, 0, 0, 0},
	}
	var rows bytes.Buffer
	var droppedFailures, addedFailures int
	for _, c := range cases {
		fmt.Fprintln(&rows, c.high, c.low, c.shift, c.plainHigh, c.wantLow)
		if c.plainLow != c.wantLow {
			droppedFailures++
		}
		if c.plainLow+c.discardedNonzero != c.wantLow {
			addedFailures++
		}
	}
	if droppedFailures != 4 || addedFailures != 2 {
		return fmt.Errorf("mutant failures: dropped=%d, added=%d; want 4 and 2", droppedFailures, addedFailures)
	}
	dir, err := os.MkdirTemp("", "golean-shrcompress-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_SHRCOMPRESS_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/shrcompress/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanShrCompress", fixtures+"bridge_test.go", fixtures+"shrcompress_test.go"); err != nil {
		return err
	}
	fmt.Printf("%s %s/%s: shrcompress passed %d cases, dropped and added sticky mutants failed %d and %d\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases), droppedFailures, addedFailures)
	return nil
}
