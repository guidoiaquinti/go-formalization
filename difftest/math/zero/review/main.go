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
	cases := []uint64{0, 1, 2, 3, 1<<63 - 1, 1 << 63, 1<<63 + 1, ^uint64(0)}
	var rows bytes.Buffer
	var omittedLowFailures, unfoldedFailures int
	for _, word := range cases {
		want := uint64(0)
		if word == 0 {
			want = 1
		}
		fmt.Fprintln(&rows, word, want)
		if ((word>>1)-1)>>63 != want {
			omittedLowFailures++
		}
		if (word-1)>>63 != want {
			unfoldedFailures++
		}
	}
	if omittedLowFailures != 1 || unfoldedFailures != 2 {
		return fmt.Errorf("mutant failures: missing low bit=%d, unfolded subtraction=%d, want one and two", omittedLowFailures, unfoldedFailures)
	}

	dir, err := os.MkdirTemp("", "golean-zero-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_FMA_ZERO_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/zero/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanFmaZero", fixtures+"bridge_test.go", fixtures+"zero_test.go"); err != nil {
		return err
	}
	fmt.Printf("%s %s/%s: zero passed %d cases, missing-low-bit and unfolded-subtraction mutants failed one and two\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases))
	return nil
}
