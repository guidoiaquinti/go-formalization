package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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
	cases := []string{"", "/", "\\", "//", "a/", "a", "abc", "\x00\xff"}
	for b := range 256 {
		cases = append(cases, string([]byte{byte(b)}), string([]byte{filepath.Separator, byte(b)}))
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		r := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			raw := make([]byte, r.Intn(128))
			r.Read(raw)
			if i%2 == 0 && len(raw) > 0 {
				raw[len(raw)-1] = filepath.Separator
			}
			cases = append(cases, string(raw))
		}
	}

	var lean, fixture bytes.Buffer
	for _, path := range cases {
		fmt.Fprint(&lean, uint8(filepath.Separator))
		for i := range len(path) {
			fmt.Fprint(&lean, " ", path[i])
		}
		fmt.Fprintln(&lean)
		fmt.Fprintf(&fixture, "x%s x%s\n", hex.EncodeToString([]byte(path)), hex.EncodeToString([]byte(expected(path))))
	}
	for i, row := range check.LeanRows("godiff-path-filepath-clean-glob-path", len(cases), &lean) {
		if want := wire(expected(cases[i])); row != want {
			return fmt.Errorf("path %q: Lean %q, want %q", cases[i], row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-clean-glob-path-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_CLEAN_GLOB_PATH_CASES", path); err != nil {
		return err
	}
	if err := check.TestInPackage("path/filepath", "TestGoLeanCleanGlobPath",
		"difftest/path/filepath/clean-glob-path/testdata/bridge_test.go",
		"difftest/path/filepath/clean-glob-path/testdata/clean_glob_path_test.go"); err != nil {
		return err
	}
	fmt.Printf("path/filepath.cleanGlobPath: %d boundaries and %d seeded rows agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(path string) string {
	if len(path) == 0 {
		return "."
	}
	if path == string(filepath.Separator) {
		return path
	}
	last := string([]byte{path[len(path)-1]})
	return strings.TrimSuffix(path, last)
}

func wire(s string) string {
	fields := make([]string, len(s))
	for i := range len(s) {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return "bytes:" + strings.Join(fields, ",")
}
