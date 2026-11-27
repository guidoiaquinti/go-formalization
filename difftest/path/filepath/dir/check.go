package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const executor = "godiff-path-filepath-dir"

func main() {
	fixed := []struct{ path, windows string }{
		{"", "."}, {"a", "."}, {"/", "\\"}, {"a/b", "a"}, {"a\\b", "a"},
		{"C:", "C:."}, {"C:foo", "C:."}, {"C:\\file", "C:\\"}, {"//host/share", "\\\\host\\share"},
		{"//host/share/leaf", "\\\\host\\share\\"}, {"//host/share/folder/leaf", "\\\\host\\share\\folder"},
		{"/??/UNC/host/share", "\\??\\UNC\\host\\share"}, {"//?/C:/foo", "\\\\?\\C:\\"},
		{"//host/../leaf", "\\"}, {"a/../c:/q", ".\\c:"}, {"\xff/\x00", "\xff"},
		{"./foo:bar/../leaf", ".\\."},
	}
	var rows bytes.Buffer
	var edges []string
	windows := filepath.Separator == '\\'
	for _, c := range fixed {
		fmt.Fprintln(&rows, encode(c.path, true))
		edges = append(edges, c.path)
		if windows && filepath.Dir(c.path) != c.windows {
			log.Fatalf("native Windows fixture %q differs from explicit result %q", c.path, c.windows)
		}
	}
	for i, row := range check.LeanRows(executor, len(fixed), &rows) {
		if want := wire(fixed[i].windows); row != want {
			log.Fatalf("simulated Windows path %q: Lean %s, want %s", fixed[i].path, row, want)
		}
	}
	fmt.Printf("filepath.Dir: %d fixed Windows profile cases agree.\n", len(fixed))
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}), string([]byte{'/', byte(b), '\\', '/'}))
	}
	check.CompareSamples(executor, edges, sample, func(path string) string { return encode(path, windows) },
		func(path string) string { return wire(filepath.Dir(path)) },
		func(input string) string {
			if windows {
				return wire(filepath.Dir(input))
			}
			return wire(path.Dir(input))
		})
}

func sample(r *rand.Rand, i int) string {
	raw := make([]byte, r.Intn(80))
	r.Read(raw)
	for j := 0; j < len(raw); j += 5 {
		raw[j] = "/\\:."[r.Intn(4)]
	}
	prefixes := []string{"", "/", "A:", "a:/", "//host/share/", "/??/UNC/host/share/", "//./", "\\\\"}
	suffixes := []string{"", "", "\\", "/"}
	return prefixes[i%len(prefixes)] + string(raw) + suffixes[r.Intn(len(suffixes))]
}

func encode(path string, windows bool) string {
	row := []byte("0")
	if windows {
		row[0] = '1'
	}
	for i := range len(path) {
		row = strconv.AppendUint(append(row, ' '), uint64(path[i]), 10)
	}
	return string(row)
}

func wire(s string) string {
	fields := make([]string, len(s))
	for i := range len(s) {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return "bytes:" + strings.Join(fields, ",")
}
