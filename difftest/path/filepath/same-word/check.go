package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ left, right string }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []input{{"", ""}, {"", "a"}, {"a", "A"}, {"K", "K"}, {"S", "ſ"}, {"σ", "ς"},
		{"Σ", "σ"}, {"I", "ı"}, {"ß", "ss"}, {"é", "e\u0301"}, {"a/b", "a\\b"},
		{"\xff", "�"}, {"\xff\x80", "��"}, {"\xff\x80", "�"}, {"\x00", ""}, {"A\x00", "a\x00"}}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		cases = append(cases, input{raw, raw}, input{raw, "�"})
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	pairs := []input{{"ꭰ", "Ꭰ"}, {"𐐀", "𐐨"}, {"K", "K"}, {"\xff", "�"}, {"a", "A"}, {"a", "a"}, {"x", "y"}}
	for _, seed := range seeds {
		r := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			if i%2 == 0 {
				var left, right strings.Builder
				for range r.Intn(32) {
					p := pairs[r.Intn(len(pairs))]
					left.WriteString(p.left)
					right.WriteString(p.right)
				}
				cases = append(cases, input{left.String(), right.String()})
			} else {
				left, right := make([]byte, r.Intn(64)), make([]byte, r.Intn(64))
				r.Read(left)
				r.Read(right)
				cases = append(cases, input{string(left), string(right)})
			}
		}
	}

	var lean, fixture bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&lean, len(c.left))
		for _, part := range []string{c.left, c.right} {
			for i := range len(part) {
				fmt.Fprint(&lean, " ", part[i])
			}
		}
		fmt.Fprintln(&lean)
		fmt.Fprintf(&fixture, "x%s x%s %t\n", hex.EncodeToString([]byte(c.left)), hex.EncodeToString([]byte(c.right)),
			expected(c, runtime.GOOS == "windows"))
	}
	for i, row := range check.LeanRows("godiff-path-filepath-same-word", len(cases), &lean) {
		c := cases[i]
		want := strconv.FormatBool(expected(c, false)) + "|" + strconv.FormatBool(expected(c, true))
		if row != want {
			return fmt.Errorf("sameWord(%q,%q): Lean %q, want %q", c.left, c.right, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-filepath-same-word-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_FILEPATH_SAME_WORD_CASES", path); err != nil {
		return err
	}
	const overlay = "difftest/path/filepath/same-word/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanSameWord", overlay+"bridge_test.go",
		overlay+"same_word_test.go"); err != nil {
		return err
	}
	fmt.Printf("path/filepath.sameWord: %d boundaries and %d seeded rows agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(c input, windows bool) bool {
	if windows {
		return strings.EqualFold(c.left, c.right)
	}
	return c.left == c.right
}
