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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []string{"", "*", "?", "[", "]", "\\", "\\[", "\xff\x00"}
	for b := range 256 {
		cases = append(cases, string([]byte{byte(b)}), string([]byte{0xff, byte(b), 0x80}))
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		r := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			raw := make([]byte, r.Intn(128))
			r.Read(raw)
			switch {
			case i%3 == 0 && len(raw) > 0:
				raw[len(raw)-1] = '\\'
			case i%3 == 1:
				for j, b := range raw {
					if marks[b] {
						raw[j] = ']'
					}
				}
			}
			cases = append(cases, string(raw))
		}
	}

	var lean, fixture bytes.Buffer
	windows := runtime.GOOS == "windows"
	for _, path := range cases {
		fmt.Fprint(&lean, 0)
		for i := range len(path) {
			fmt.Fprint(&lean, " ", path[i])
		}
		fmt.Fprintln(&lean)
		fmt.Fprintf(&fixture, "x%s %t\n", hex.EncodeToString([]byte(path)), expected(path, windows))
	}
	for i, row := range check.LeanRows("godiff-path-filepath-has-meta", len(cases), &lean) {
		path := cases[i]
		want := strconv.FormatBool(expected(path, false)) + "|" + strconv.FormatBool(expected(path, true))
		if row != want {
			return fmt.Errorf("path %q: Lean %q, want %q", path, row, want)
		}
	}

	dir, err := os.MkdirTemp("", "golean-has-meta-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_FILEPATH_HAS_META_CASES", path); err != nil {
		return err
	}
	const overlay = "difftest/path/filepath/has-meta/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanHasMeta",
		overlay+"bridge_test.go", overlay+"has_meta_test.go"); err != nil {
		return err
	}
	fmt.Printf("path/filepath.hasMeta: %d boundaries and %d seeded rows agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

var marks = [256]bool{'*': true, '?': true, '[': true, '\\': true}

func expected(path string, windows bool) bool {
	for i := range len(path) {
		if marks[path[i]] && (!windows || path[i] != '\\') {
			return true
		}
	}
	return false
}
