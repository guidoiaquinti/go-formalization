package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []struct {
		text string
		key  rune
		want bool
	}{
		{"", 0, false},
		{"é\x00", 0, true},
		{"©", 0xa9, true},
		{"©", 0xc2, false},
		{"🚀é", 'é', true},
		{"é\xf0\x9f\x9a", '🚀', false},
		{"é🚀", '🚀', true},
		{"\xf0a🚀", 'a', true},
		{"é\xe0\xa0", utf8.RuneError, true},
		{"é\xed\xa0\x80", utf8.RuneError, true},
		{"é\xf4\x90\x80\x80", utf8.RuneError, true},
		{"é\xc0\xaf", utf8.RuneError, true},
		{"é�", utf8.RuneError, true},
		{"é�", 0xd800, false},
		{"é�", -1, false},
		{"é�", math.MinInt32, false},
		{"é�", 0x110000, false},
		{"é�", math.MaxInt32, false},
		{strings.Repeat("©", 512), utf8.RuneError, false},
	}
	var input, expected bytes.Buffer
	encoded, sanitize := 0, 0
	for _, c := range cases {
		fmt.Fprint(&input, c.key)
		for _, b := range []byte(c.text) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		fmt.Fprintf(&expected, "%x\t%d\t%t\n", c.text, c.key, c.want)
		if strings.Contains(c.text, string(c.key)) != c.want {
			encoded++
		}
		found := false
		for _, r := range strings.ToValidUTF8(c.text, "") {
			found = found || r == c.key
		}
		if found != c.want {
			sanitize++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-contains-rune-private", len(cases), &input) {
		if row != strconv.FormatBool(cases[i].want) {
			return fmt.Errorf("Lean fixture %d returned %s, want %t", i, row, cases[i].want)
		}
	}
	if encoded != 9 || sanitize != 4 {
		return fmt.Errorf("mutation counts encoded=%d sanitize=%d, want 9 and 4", encoded, sanitize)
	}
	dir, err := os.MkdirTemp("", "golean-bytes-contains-rune-private-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, expected.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_CONTAINS_RUNE_PRIVATE_CASES", file); err != nil {
		return err
	}
	const path = "difftest/bytes/contains-rune-private/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanContainsRunePrivate", path+"bridge_test.go", path+"contains_rune_private_test.go"); err != nil {
		return err
	}
	fmt.Printf("private bytes.containsRune passed %d fixed Go and Lean fixtures; encoded search and deleting invalid input mutants failed %d and %d.\n",
		len(cases), encoded, sanitize)
	return nil
}
