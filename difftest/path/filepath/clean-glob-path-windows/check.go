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

type input struct {
	path   string
	volume int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []input{{"", 0}, {"/", 0}, {"\\", 0}, {"C:", 2}, {"C:/", 2}, {"C:\\", 2},
		{"C:one/", 2}, {"//", 2}, {"//host/share", 12}, {"//host/share/", 12},
		{"//host/../", 0}, {"//?", 3}, {"//?/unc", 7}, {"//?/unc/host/share", 18},
		{"/??/dev/", 7}, {"\xff:/\x00", 2}}
	for byteValue := range 256 {
		cases = append(cases, input{string([]byte{byte(byteValue), ':'}), 2},
			input{string([]byte{byte(byteValue), ':', '/'}), 2})
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	prefixes := []input{{"q_", 0}, {"X:", 2}, {"//host/share/", 12}, {"//?/dev/", 7},
		{"/??/UNC/host/share/", 18}, {"//host/../", 0}, {"/_", 0}}
	for _, seed := range seeds {
		r := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			raw := make([]byte, r.Intn(80))
			r.Read(raw)
			if i%2 == 0 && len(raw) > 0 {
				raw[len(raw)-1] = "/\\"[r.Intn(2)]
			}
			c := prefixes[i%len(prefixes)]
			c.path += string(raw)
			cases = append(cases, c)
		}
	}

	var lean, fixture bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&lean, 0)
		for i := range len(c.path) {
			fmt.Fprint(&lean, " ", c.path[i])
		}
		fmt.Fprintln(&lean)
		prefix, path := expected(c, runtime.GOOS == "windows")
		fmt.Fprintf(&fixture, "x%s %d x%s\n", hex.EncodeToString([]byte(c.path)), prefix,
			hex.EncodeToString([]byte(path)))
	}
	for i, row := range check.LeanRows("godiff-path-filepath-clean-glob-path-windows", len(cases), &lean) {
		c := cases[i]
		want := wire(expected(c, false)) + "|" + wire(expected(c, true))
		if row != want {
			return fmt.Errorf("path %q: Lean %q, want %q", c.path, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-clean-glob-windows-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_FILEPATH_CLEAN_GLOB_WINDOWS_CASES", path); err != nil {
		return err
	}
	const overlay = "difftest/path/filepath/clean-glob-path-windows/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanCleanGlobWindows", overlay+"bridge_test.go",
		overlay+"clean_glob_windows_test.go"); err != nil {
		return err
	}
	fmt.Printf("path/filepath.cleanGlobPathWindows: %d boundaries and %d seeded rows agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(c input, windows bool) (int, string) {
	if c.path == "" {
		return 0, "."
	}
	volume, separators := 0, "/"
	if windows {
		volume, separators = c.volume, "/\\"
	}
	last := len(c.path) - 1
	if last == volume && strings.ContainsRune(separators, rune(c.path[last])) {
		return len(c.path), c.path
	}
	if volume == 2 && len(c.path) == 2 {
		return 2, c.path + "."
	}
	return min(volume, last), c.path[:last]
}

func wire(prefix int, path string) string {
	fields := make([]string, len(path))
	for i := range len(path) {
		fields[i] = strconv.Itoa(int(path[i]))
	}
	return strconv.Itoa(prefix) + ":" + strings.Join(fields, ",")
}
