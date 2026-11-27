package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
	"github.com/guidoiaquinti/go-formalization/difftest/path/filepath/internal/islocaldevice"
)

const executor = "godiff-path-filepath-localize"

func main() {
	profile := map[string]int{"windows": 1, "plan9": 2}[runtime.GOOS]
	fixed := []struct {
		path    string
		profile int
		device  bool
		result  string
	}{
		{".", 0, false, "."},
		{".", 1, false, "."},
		{".", 2, false, "."},
		{"a/b", 1, false, `a\b`},
		{"a//b", 1, false, ""},
		{"x/./y", 2, false, ""},
		{"a/../b", 0, false, ""},
		{"/leaf", 0, false, ""},
		{"tail/", 2, false, ""},
		{"#dev/file", 2, false, ""},
		{"x/#dev", 2, false, "x/#dev"},
		{"C:leaf", 1, false, ""},
		{"C:leaf", 0, false, "C:leaf"},
		{`a\b`, 1, false, ""},
		{`a\b`, 2, false, `a\b`},
		{"a/COM¹", 1, false, ""},
		{"NUL.txt/z", 1, false, `NUL.txt\z`},
		{"NUL.txt/z", 1, true, ""},
		{"a/\x00", 0, false, ""},
		{"x/\xff", 1, false, ""},
		{"a/水", 1, false, `a\水`},
	}
	var input bytes.Buffer
	var edges []string
	for _, fixture := range fixed {
		fmt.Fprintln(&input, encode(fixture.path, fixture.profile, func(string) bool { return fixture.device }))
		edges = append(edges, fixture.path)
	}
	for i, row := range check.LeanRows(executor, len(fixed), &input) {
		f := fixed[i]
		want := "invalid path"
		if f.result != "" {
			want = wire(f.result)
		}
		if row != want {
			log.Fatalf("Localize profile %d for %q: Lean %q, want %q", f.profile, f.path, row, want)
		}
	}
	fmt.Printf("filepath.Localize: %d fixed profile cases passed.\n", len(fixed))
	for i := range 256 {
		raw := string([]byte{byte(i)})
		edges = append(edges, raw, "a/"+raw)
	}
	check.CompareSamples(executor, edges, sample,
		func(s string) string { return encode(s, profile, islocaldevice.IsDOSName) }, native, native)
}

func native(path string) string {
	local, err := filepath.Localize(path)
	if err != nil {
		if local != "" {
			log.Fatalf("Localize(%q) produced %q with error %v", path, local, err)
		}
		return err.Error()
	}
	return wire(local)
}

func sample(r *rand.Rand, _ int) string {
	fragments := []string{"", ".", "..", "a", "CON.txt", "z:", `a\b`, "#dev", "\x00", "🍊", "ÿ"}
	parts := make([]string, 1+r.Intn(7))
	for i := range parts {
		if r.Intn(3) == 0 {
			raw := make([]byte, r.Intn(8))
			r.Read(raw)
			parts[i] = string(raw)
		} else {
			parts[i] = fragments[r.Intn(len(fragments))]
		}
	}
	return strings.Join(parts, "/")
}

func wire(path string) string {
	fields := make([]string, len(path))
	for i := range len(path) {
		fields[i] = strconv.Itoa(int(path[i]))
	}
	return "ok:" + strings.Join(fields, ",")
}
func encode(s string, profile int, device func(string) bool) string {
	var reserved []string
	if profile == 1 {
		for _, part := range strings.FieldsFunc(s, func(r rune) bool { return r == '/' || r == '\\' }) {
			if device(part) {
				reserved = append(reserved, part)
			}
		}
	}
	row := fmt.Appendf(nil, "%d %d", profile, len(reserved))
	for _, part := range reserved {
		row = fmt.Appendf(row, " %d", len(part))
		for i := range len(part) {
			row = strconv.AppendUint(append(row, ' '), uint64(part[i]), 10)
		}
	}
	for i := range len(s) {
		row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
	}
	return string(row)
}
