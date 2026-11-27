package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const executor = "godiff-path-filepath-split"

func main() {
	// Fixed Windows results use the original input bytes, even when the
	// volume classification accepts mixed separator spellings.
	windowsFixtures := []struct{ path, dir, file string }{
		{"", "", ""}, {"a", "", "a"}, {"/", "/", ""}, {"a/b", "a/", "b"},
		{"a\\b", "a\\", "b"}, {"C:", "C:", ""}, {"C:foo", "C:", "foo"},
		{"C:\\file", "C:\\", "file"}, {"//host/share", "//host/share", ""},
		{"//host/share/leaf", "//host/share/", "leaf"},
		{"/??/UNC/host/share", "/??/UNC/host/share", ""},
		{"//?/C:/foo", "//?/C:/", "foo"}, {"//host/../leaf", "//host/../", "leaf"},
		{"//host/..", "//host/", ".."}, {"a//", "a//", ""}, {"\xff/\x00", "\xff/", "\x00"},
	}
	var rows bytes.Buffer
	var edges []string
	windows := filepath.Separator == '\\'
	for _, c := range windowsFixtures {
		fmt.Fprintln(&rows, encode(c.path, true))
		edges = append(edges, c.path)
		if windows && native(c.path) != wire(c.dir, c.file) {
			log.Fatalf("native Windows fixture %q differs from expected components", c.path)
		}
	}
	for i, row := range check.LeanRows(executor, len(windowsFixtures), &rows) {
		c := windowsFixtures[i]
		if want := wire(c.dir, c.file); row != want {
			log.Fatalf("simulated Windows path %q: Lean %s, want %s", c.path, row, want)
		}
	}
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}), string([]byte{'/', byte(b), '/', 0xff}))
	}
	check.CompareSamples(executor, edges, sample, func(path string) string { return encode(path, windows) },
		native, func(path string) string {
			if windows {
				return native(path)
			}
			i := strings.LastIndexByte(path, '/') + 1
			return wire(path[:i], path[i:])
		})
}

func sample(r *rand.Rand, i int) string {
	raw := make([]byte, r.Intn(96))
	r.Read(raw)
	for j := 0; j < len(raw); j += 5 {
		raw[j] = "/\\:."[r.Intn(4)]
	}
	prefixes := []string{"", "/", "A:", "a:/", "//host/share/", "/??/UNC/host/share/", "//./", "\\\\"}
	return prefixes[i%len(prefixes)] + string(raw)
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

func native(path string) string {
	dir, file := filepath.Split(path)
	if dir+file != path {
		log.Fatalf("filepath.Split(%q) did not preserve the input", path)
	}
	return wire(dir, file)
}

func wire(dir, file string) string {
	encode := func(s string) string {
		fields := make([]string, len(s))
		for i := range len(s) {
			fields[i] = strconv.Itoa(int(s[i]))
		}
		return strings.Join(fields, ",")
	}
	return encode(dir) + "|" + encode(file)
}
