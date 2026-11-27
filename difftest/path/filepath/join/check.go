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

const executor = "godiff-path-filepath-join"

func main() {
	fixed := []struct {
		elements []string
		windows  string
	}{
		{nil, ""},
		{[]string{"", ""}, ""},
		{[]string{"c:", "f"}, "c:f"},
		{[]string{"c:", "\\f"}, "c:\\f"},
		{[]string{"c:", ""}, "c:."},
		{[]string{"\\", "\\host", "share"}, "\\host\\share"},
		{[]string{"\\\\", "host", "share"}, "\\\\host\\share"},
		{[]string{"\\", "??", "c:\\x"}, "\\.\\??\\c:\\x"},
		{[]string{"a", "/../c:"}, ".\\c:"},
		{[]string{"C:\\a", "", "..", "b"}, "C:\\b"},
		{[]string{"//host/share", "..", "q"}, "\\\\host\\share\\q"},
		{[]string{"foo:bar"}, "foo:bar"},
		{[]string{"/", "/"}, "\\"},
		{[]string{"//", ""}, "\\\\"},
		{[]string{"a", "", ".."}, "."},
		{[]string{"", "\xff", "..", "\x00"}, "\x00"},
	}
	var rows bytes.Buffer
	var edges [][]string
	windows := filepath.Separator == '\\'
	for _, c := range fixed {
		fmt.Fprintln(&rows, encode(c.elements, true))
		edges = append(edges, c.elements)
		if windows && filepath.Join(c.elements...) != c.windows {
			log.Fatalf("native Windows fixture %q differs from explicit result %q", c.elements, c.windows)
		}
	}
	for i, row := range check.LeanRows(executor, len(fixed), &rows) {
		if want := wire(fixed[i].windows); row != want {
			log.Fatalf("simulated Windows elements %q: Lean %s, want %s", fixed[i].elements, row, want)
		}
	}
	fmt.Printf("filepath.Join: %d fixed Windows profile cases agree.\n", len(fixed))
	for b := range 256 {
		raw := string([]byte{byte(b)})
		edges = append(edges, []string{raw}, []string{"", "c:", raw, "..", "leaf"})
	}
	check.CompareSamples(executor, edges, sample, func(parts []string) string { return encode(parts, windows) },
		func(parts []string) string { return wire(filepath.Join(parts...)) },
		func(parts []string) string {
			if windows {
				return wire(filepath.Join(parts...))
			}
			return wire(path.Join(parts...))
		})
}

func sample(r *rand.Rand, _ int) []string {
	values := []string{"", ".", "..", "a", "/", "\\", "//", "c:", "d:/", "?", "??", "c:leaf", "\xff"}
	parts := make([]string, r.Intn(12))
	for i := range parts {
		if r.Intn(3) == 0 {
			raw := make([]byte, r.Intn(24))
			r.Read(raw)
			for j := 0; j < len(raw); j += 5 {
				raw[j] = "/\\:."[r.Intn(4)]
			}
			parts[i] = string(raw)
		} else {
			parts[i] = values[r.Intn(len(values))]
		}
	}
	return parts
}

func encode(parts []string, windows bool) string {
	profile := 0
	if windows {
		profile = 1
	}
	row := fmt.Appendf(nil, "%d %d", profile, len(parts))
	for _, part := range parts {
		row = fmt.Appendf(row, " %d", len(part))
		for i := range len(part) {
			row = strconv.AppendUint(append(row, ' '), uint64(part[i]), 10)
		}
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
