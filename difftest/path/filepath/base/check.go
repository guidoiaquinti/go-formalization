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

const executor = "godiff-path-filepath-base"

func main() {
	fixed := []struct{ path, windows string }{
		{"", "."}, {"a", "a"}, {"/", "\\"}, {"////", "\\"}, {"a//", "a"},
		{"C:", "\\"}, {"C:foo", "foo"}, {"C:\\x\\", "x"}, {"//host/share", "\\"},
		{"//host/share/leaf", "leaf"}, {"//host/..", ".."}, {"/??/UNC/host/share", "\\"},
		{"//?/C:/foo", "foo"}, {"a\\b", "b"}, {"a/.", "."}, {"\xff/\x00", "\x00"},
	}
	var rows bytes.Buffer
	var edges []string
	windows := filepath.Separator == '\\'
	for _, c := range fixed {
		fmt.Fprintln(&rows, encode(c.path, true))
		edges = append(edges, c.path)
		if windows && filepath.Base(c.path) != c.windows {
			log.Fatalf("native Windows fixture %q differs from explicit result %q", c.path, c.windows)
		}
	}
	for i, row := range check.LeanRows(executor, len(fixed), &rows) {
		if want := wire(fixed[i].windows); row != want {
			log.Fatalf("simulated Windows path %q: Lean %s, want %s", fixed[i].path, row, want)
		}
	}
	fmt.Printf("filepath.Base: %d fixed Windows profile cases agree.\n", len(fixed))
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}), string([]byte{'/', byte(b), '\\', '/'}))
	}
	check.CompareSamples(executor, edges, sample, func(path string) string { return encode(path, windows) },
		func(path string) string { return wire(filepath.Base(path)) },
		func(input string) string {
			if windows {
				return wire(filepath.Base(input))
			}
			return wire(path.Base(input))
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
