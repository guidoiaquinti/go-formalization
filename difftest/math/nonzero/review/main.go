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
	cases := []struct{ word, want uint64 }{
		{0, 0},
		{1, 1},
		{2, 1},
		{1 << 31, 1},
		{1 << 32, 1},
		{1 << 63, 1},
		{1<<63 | 1, 1},
		{^uint64(0), 1},
	}
	var rows bytes.Buffer
	var truncatedFailures, identityFailures int
	for _, c := range cases {
		fmt.Fprintln(&rows, c.word, c.want)
		truncated := uint64(0)
		if uint32(c.word) != 0 {
			truncated = 1
		}
		if truncated != c.want {
			truncatedFailures++
		}
		if c.word != c.want {
			identityFailures++
		}
	}
	if truncatedFailures != 2 || identityFailures != 6 {
		return fmt.Errorf("mutant failures: truncated=%d, identity=%d, want two and six", truncatedFailures, identityFailures)
	}

	dir, err := os.MkdirTemp("", "golean-nonzero-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_FMA_NONZERO_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/nonzero/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanFmaNonzero", fixtures+"bridge_test.go", fixtures+"nonzero_test.go"); err != nil {
		return err
	}
	fmt.Printf("%s %s/%s: nonzero passed %d cases, truncation and identity mutants failed two and six\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases))
	return nil
}
