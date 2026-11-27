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

const executor = "godiff-path-filepath-clean"

func main() {
	fixed := []struct{ path, windows string }{
		{"", "."}, {"c:", "c:."}, {"c:/", "c:\\"}, {"c:./a", "c:a"},
		{"c:/../../x", "c:\\x"}, {"a/../c:", ".\\c:"}, {"a/../../c:", "..\\c:"},
		{"foo:bar", "foo:bar"}, {"/a/../??/a", "\\.\\??\\a"}, {"\\a\\..\\??\\q", "\\.\\??\\q"},
		{".\\c:foo", ".\\c:foo"}, {"//host/share/foo/../leaf", "\\\\host\\share\\leaf"},
		{"//host/share", "\\\\host\\share"}, {"//host/../leaf", "\\leaf"},
		{"//?/UNC/host/share/../x", "\\\\?\\UNC\\host\\share\\x"},
		{"/??/UNC/host/share/../x", "\\??\\UNC\\host\\share\\x"},
		{"//?/c:/a/..", "\\\\?\\c:\\"}, {"///", "\\\\\\"}, {".\\..", ".."},
		{"\xff/../\x00", "\x00"},
		{"./foo:bar/..", ".\\."}, {"\\.\\??\\..", "\\.\\"},
		{"./foo:bar/../ab", ".\\ab"}, {"foo:bar/..", "."},
	}
	var rows bytes.Buffer
	var edges []string
	windows := filepath.Separator == '\\'
	for _, c := range fixed {
		fmt.Fprintln(&rows, encode(c.path, true))
		edges = append(edges, c.path)
		if windows && filepath.Clean(c.path) != c.windows {
			log.Fatalf("native Windows fixture %q differs from explicit result %q", c.path, c.windows)
		}
	}
	for i, row := range check.LeanRows(executor, len(fixed), &rows) {
		if want := wire(fixed[i].windows); row != want {
			log.Fatalf("simulated Windows path %q: Lean %s, want %s", fixed[i].path, row, want)
		}
	}
	fmt.Printf("filepath.Clean: %d fixed Windows cases agree.\n", len(fixed))
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}), "a/../"+string([]byte{byte(b), '\\', '/'}))
	}
	check.CompareSamples(executor, edges, sample, func(path string) string { return encode(path, windows) },
		func(path string) string { return wire(filepath.Clean(path)) },
		func(input string) string {
			if windows {
				return wire(filepath.Clean(input))
			}
			return wire(path.Clean(input))
		})
}

func sample(r *rand.Rand, i int) string {
	prefixes := []string{"", "/", "C:", "C:/", "//host/share/", "/??/UNC/host/share/", "//./", "\\"}
	parts := []string{"", ".", "..", "..", "c:", "??", "a", "\xff", "\x00", "😺"}
	var out strings.Builder
	out.WriteString(prefixes[i%len(prefixes)])
	for j := range r.Intn(12) {
		if j > 0 {
			out.WriteByte("/\\"[r.Intn(2)])
		}
		if r.Intn(4) == 0 {
			raw := make([]byte, r.Intn(8))
			r.Read(raw)
			out.Write(raw)
		} else {
			out.WriteString(parts[r.Intn(len(parts))])
		}
	}
	if i%3 == 0 {
		out.WriteByte("/\\"[r.Intn(2)])
	}
	return out.String()
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
