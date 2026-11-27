package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const casesEnv = "GOLEAN_STRINGS_IS_SEPARATOR_CASES"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-is-separator-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	fixture, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fixture.Close()
	rows := bufio.NewWriter(fixture)
	if err := os.Setenv(casesEnv, path); err != nil {
		return err
	}

	// Check both the classification and the category disjointness used in
	// the model for every code point, including the surrogate interval.
	for first := rune(0); first <= unicode.MaxRune; first += 65_536 {
		last := min(first+65_535, unicode.MaxRune)
		runes := make([]rune, 0, last-first+1)
		for r := first; r <= last; r++ {
			if unicode.Is(unicode.White_Space, r) && (unicode.IsLetter(r) || unicode.IsDigit(r)) {
				return fmt.Errorf("White_Space overlaps a word category at U+%04X", r)
			}
			runes = append(runes, r)
		}
		if err := compare(rows, fmt.Sprintf("U+%04X–U+%04X", first, last), runes); err != nil {
			return err
		}
	}
	edges := []rune{math.MinInt32, math.MinInt32 + 1, -256, -129, -128, -2, -1,
		unicode.MaxRune + 1, unicode.MaxRune + 2, 1 << 24, math.MaxInt32 - 1, math.MaxInt32}
	if err := compare(rows, "outside Unicode boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		runes := make([]rune, 20_000)
		for i := range runes {
			runes[i] = rune(rng.Uint32())
		}
		if err := compare(rows, fmt.Sprintf("seed=%d", seed), runes); err != nil {
			return err
		}
	}
	if err := rows.Flush(); err != nil {
		return err
	}
	if err := fixture.Close(); err != nil {
		return err
	}
	const fixtures = "difftest/strings/is-separator/testdata/"
	return check.TestInPackage("strings", "TestGoLeanIsSeparator",
		fixtures+"bridge_test.go", fixtures+"is_separator_test.go")
}

func expected(r rune) bool {
	if r <= unicode.MaxASCII {
		return !strings.ContainsRune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_", r)
	}
	return unicode.Is(unicode.White_Space, r)
}

func compare(fixture *bufio.Writer, label string, runes []rune) error {
	var modelInput bytes.Buffer
	for _, r := range runes {
		fmt.Fprintln(&modelInput, r)
		if _, err := fmt.Fprintf(fixture, "%d\t%t\n", r, expected(r)); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-is-separator", len(runes), &modelInput) {
		r := runes[i]
		want := strconv.FormatBool(expected(r))
		if got != want {
			return fmt.Errorf("%s isSeparator(%d): Lean=%s want=%s", label, r, got, want)
		}
	}
	fmt.Printf("strings.isSeparator %s: %d Lean inputs agree.\n", label, len(runes))
	return nil
}
