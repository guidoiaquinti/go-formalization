package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func wire(path string) string {
	var fields []string
	for i := range len(path) {
		fields = append(fields, strconv.Itoa(int(path[i])))
	}
	return strings.Join(fields, ",")
}

func coalesce(path string) string {
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	return strings.ReplaceAll(path, "/", "\\")
}

func main() {
	fixtures := []struct{ raw, backslash string }{
		{"", ""}, {"/", "\\"}, {"\\", "\\"}, {"a\\b", "a\\b"},
		{"//a///b", "\\\\a\\\\\\b"},
		{"C:/a\\b/", "C:\\a\\b\\"},
		{"a\x00/\xff//", "a\x00\\\xff\\\\"},
	}
	var rows bytes.Buffer
	unchanged, collapsed := 0, 0
	for _, f := range fixtures {
		fmt.Fprint(&rows, "0")
		for i := range len(f.raw) {
			fmt.Fprint(&rows, " ", f.raw[i])
		}
		fmt.Fprintln(&rows)
		want := f.raw
		if filepath.Separator == '\\' {
			want = f.backslash
		}
		if got := filepath.FromSlash(f.raw); got != want {
			log.Fatalf("native FromSlash(%x) = %x, want %x", f.raw, got, want)
		}
		if f.raw != f.backslash {
			unchanged++
		}
		if coalesce(f.raw) != f.backslash {
			collapsed++
		}
	}
	for i, got := range check.LeanRows("godiff-path-filepath-from-slash", len(fixtures), &rows) {
		f := fixtures[i]
		if want := wire(f.raw) + "|" + wire(f.backslash); got != want {
			log.Fatalf("Lean profiles for %x gave %q, want %q", f.raw, got, want)
		}
	}
	if unchanged != 4 || collapsed != 2 {
		log.Fatalf("mutation disagreements unchanged=%d and coalesced=%d, want 4 and 2", unchanged, collapsed)
	}
	fmt.Printf("FromSlash passed %d native and Lean profile fixtures, unchanged/coalesced mutants failed %d/%d\n", len(fixtures), unchanged, collapsed)
}
