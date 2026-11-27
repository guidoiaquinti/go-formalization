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

func main() {
	fixtures := []struct{ raw, unix, windows string }{
		{"", "", ""}, {".profile", ".profile", ".profile"}, {"parent.dot/leaf", "", ""},
		{"parent.dot\\leaf", ".dot\\leaf", ""}, {"a.txt\\", ".txt\\", ""},
		{"foo.tar.gz", ".gz", ".gz"}, {"one/.", ".", "."}, {"bad.\xff\x00", ".\xff\x00", ".\xff\x00"},
		{"a.x//", "", ""}, {"C:\\a.x\\b", ".x\\b", ""}, {"//a////b.c/", "", ""},
	}
	var rows bytes.Buffer
	unixOnly, ignoreSeparators := 0, 0
	for _, f := range fixtures {
		fmt.Fprint(&rows, "0")
		for i := range len(f.raw) {
			fmt.Fprint(&rows, " ", f.raw[i])
		}
		fmt.Fprintln(&rows)
		want := f.unix
		if filepath.Separator == '\\' {
			want = f.windows
		}
		if got := filepath.Ext(f.raw); got != want {
			log.Fatalf("native Ext(%x)=%x, want %x", f.raw, got, want)
		}
		if f.unix != f.windows {
			unixOnly++
		}
		global := ""
		if i := strings.LastIndexByte(f.raw, '.'); i >= 0 {
			global = f.raw[i:]
		}
		if global != f.unix {
			ignoreSeparators++
		}
	}
	for i, got := range check.LeanRows("godiff-path-filepath-ext", len(fixtures), &rows) {
		f := fixtures[i]
		if want := wire(f.unix) + "|" + wire(f.windows); got != want {
			log.Fatalf("Lean profiles for %x gave %q, want %q", f.raw, got, want)
		}
	}
	if unixOnly != 3 || ignoreSeparators != 3 {
		log.Fatalf("mutation counts unix-only=%d global-dot=%d, want 3 and 3", unixOnly, ignoreSeparators)
	}
	fmt.Printf("Ext passed %d native and Lean profile fixtures, Unix-only and global-dot mutants failed %d/%d\n", len(fixtures), unixOnly, ignoreSeparators)
}
