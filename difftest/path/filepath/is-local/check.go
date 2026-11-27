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

const executor = "godiff-path-filepath-is-local"

func main() {
	profile := 0
	switch runtime.GOOS {
	case "windows":
		profile = 1
	case "plan9":
		profile = 2
	}
	cases := []struct {
		path                                  string
		windowsOrdinary, windowsDevice, plan9 bool
	}{
		{"", false, false, false},
		{"../leaf", false, false, false},
		{"a/../leaf", true, true, true},
		{"a/../..", false, false, false},
		{"#device", true, true, false},
		{"a/#device", true, true, true},
		{"/leaf", false, false, false},
		{"\\leaf", false, false, true},
		{"C:leaf", false, false, true},
		{"a:b", false, false, true},
		{"cOn", false, false, true},
		{"CON.txt", true, false, true},
		{"AUX /../leaf", true, false, true},
		{"COM¹", false, false, true},
		{"COM0", true, true, true},
		{"LPT².txt", true, false, true},
		{"CONIN$", false, false, true},
		{"CONOUT$", false, false, true},
		{"CON.txt\x00", true, true, true},
		{"\xff/", true, true, true},
		{"a/../NUL/..", false, false, true},
	}
	var fixed bytes.Buffer
	edges := make([]string, 0, len(cases)+512)
	for _, c := range cases {
		fmt.Fprintln(&fixed, encode(c.path, 1, func(string) bool { return false }))
		fmt.Fprintln(&fixed, encode(c.path, 1, func(string) bool { return true }))
		fmt.Fprintln(&fixed, encode(c.path, 2, islocaldevice.IsDOSName))
		edges = append(edges, c.path)
	}
	for i, row := range check.LeanRows(executor, len(cases)*3, &fixed) {
		c := cases[i/3]
		want := []bool{c.windowsOrdinary, c.windowsDevice, c.plan9}[i%3]
		if row != strconv.FormatBool(want) {
			log.Fatalf("profile case %d for %q: Lean %s, want %t", i%3, c.path, row, want)
		}
	}
	fmt.Printf("filepath.IsLocal: %d Windows and Plan 9 profile observations passed.\n", len(cases)*3)
	for i := range 256 {
		raw := string([]byte{byte(i)})
		edges = append(edges, raw, "a/../"+raw)
	}
	native := func(s string) string { return strconv.FormatBool(filepath.IsLocal(s)) }
	input := func(s string) string { return encode(s, profile, islocaldevice.IsDOSName) }
	check.CompareSamples(executor, edges, sample, input, native, native)
}

func sample(r *rand.Rand, _ int) string {
	tokens := []string{"", ".", "..", "a", "c:", "/", "\\", "#", "NUL.txt", "LPT¹", "COM0", "CONIN$"}
	raw := make([]byte, 0, 32)
	for range r.Intn(12) {
		if r.Intn(3) != 0 {
			raw = append(raw, tokens[r.Intn(len(tokens))]...)
		} else {
			raw = append(raw, byte(r.Intn(256)))
		}
	}
	return string(raw)
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
