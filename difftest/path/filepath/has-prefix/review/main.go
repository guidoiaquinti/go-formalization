package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	cases := []struct {
		path, prefix  string
		unix, windows bool
	}{
		{"", "", true, true},
		{"", "x", false, false},
		{"abc", "ab", true, true},
		{"C:/X", "c:", false, true},
		{"A/leaf", `a\`, false, false},
		{"Éx", "\xc3", true, true},
		{"Éx", "é", false, true},
		{"\xffx", "\xfe", false, true},
		{"\xff", "�", false, true},
		{"ſx", "s", false, false},
		{"ςx", "Σ", false, false},
		{"Kx", "k", false, true},
		{"K", "K", false, true},
		{"a\x00", "a", true, true},
	}
	var rows bytes.Buffer
	var expected []bool
	for i, c := range cases {
		nativeWant := c.unix
		if runtime.GOOS == "windows" {
			nativeWant = c.windows
		}
		if got := filepath.HasPrefix(c.path, c.prefix); got != nativeWant {
			panic(fmt.Sprintf("native row %d: got %t want %t", i, got, nativeWant))
		}
		windowsGo := strings.HasPrefix(c.path, c.prefix) ||
			strings.HasPrefix(strings.ToLower(c.path), strings.ToLower(c.prefix))
		if windowsGo != c.windows {
			panic(fmt.Sprintf("portable Windows row %d: got %t want %t", i, windowsGo, c.windows))
		}
		for profile, want := range []bool{c.unix, c.windows} {
			fmt.Fprintf(&rows, "%d %d", profile, len(c.path))
			for _, b := range []byte(c.path + c.prefix) {
				fmt.Fprintf(&rows, " %d", b)
			}
			fmt.Fprintln(&rows)
			expected = append(expected, want)
		}
	}
	binary := "lean/.lake/build/bin/godiff-path-filepath-has-prefix"
	if len(os.Args) == 2 {
		binary = os.Args[1]
	}
	cmd := exec.Command(binary)
	cmd.Stdin = &rows
	raw, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n")
	if len(lines) != len(expected) {
		panic(fmt.Sprintf("Lean produced %d rows", len(lines)))
	}
	for i, want := range expected {
		if lines[i] != strconv.FormatBool(want) {
			panic(fmt.Sprintf("Lean row %d: got %q want %t", i, lines[i], want))
		}
	}
	partial := cases[5]
	if strings.HasPrefix(strings.ToLower(partial.path), strings.ToLower(partial.prefix)) == partial.windows {
		panic("lowercase-only mutant survived the partial UTF-8 prefix")
	}
	fold := cases[9]
	if strings.EqualFold("ſ", fold.prefix) == fold.windows {
		panic("EqualFold mutant survived")
	}
	fmt.Printf("filepath.HasPrefix passed %d native and portable Windows fixtures and %d Lean profile rows\n", len(cases), len(expected))
}
