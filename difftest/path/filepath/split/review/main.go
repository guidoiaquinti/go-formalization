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

type pair struct{ dir, file string }

func wire(p pair) string {
	part := func(s string) string {
		fields := make([]string, len(s))
		for i := range len(s) {
			fields[i] = strconv.Itoa(int(s[i]))
		}
		return strings.Join(fields, ",")
	}
	return part(p.dir) + "|" + part(p.file)
}

func main() {
	cases := []struct {
		input         string
		unix, windows pair
	}{
		{"", pair{}, pair{}},
		{"a//.././b", pair{"a//.././", "b"}, pair{"a//.././", "b"}},
		{"a//", pair{"a//", ""}, pair{"a//", ""}},
		{"C:", pair{"", "C:"}, pair{"C:", ""}},
		{"/:name", pair{"/", ":name"}, pair{"/:", "name"}},
		{"Z:leaf", pair{"", "Z:leaf"}, pair{"Z:", "leaf"}},
		{"a\\..\\leaf", pair{"", "a\\..\\leaf"}, pair{"a\\..\\", "leaf"}},
		{"//h/share", pair{"//h/", "share"}, pair{"//h/share", ""}},
		{"/\\h/share", pair{"/\\h/", "share"}, pair{"/\\h/share", ""}},
		{"//h/..", pair{"//h/", ".."}, pair{"//h/", ".."}},
		{"//h/share/..", pair{"//h/share/", ".."}, pair{"//h/share/", ".."}},
		{"/??/uNc/h/share", pair{"/??/uNc/h/", "share"}, pair{"/??/uNc/h/share", ""}},
		{"//?/C:", pair{"//?/", "C:"}, pair{"//?/C:", ""}},
		{"//?/..", pair{"//?/", ".."}, pair{"//?/", ".."}},
		{"\xff/\x00\\\x80", pair{"\xff/", "\x00\\\x80"}, pair{"\xff/\x00\\", "\x80"}},
	}
	var rows bytes.Buffer
	var expected []pair
	for i, c := range cases {
		dir, file := filepath.Split(c.input)
		nativeWant := c.unix
		if runtime.GOOS == "windows" {
			nativeWant = c.windows
		}
		if (pair{dir, file}) != nativeWant {
			panic(fmt.Sprintf("native row %d: %q %q want %+v", i, dir, file, nativeWant))
		}
		for profile, want := range []pair{c.unix, c.windows} {
			if want.dir+want.file != c.input {
				panic(fmt.Sprintf("fixture %d does not reconstruct input", i))
			}
			fmt.Fprintf(&rows, "%d", profile)
			for _, b := range []byte(c.input) {
				fmt.Fprintf(&rows, " %d", b)
			}
			fmt.Fprintln(&rows)
			expected = append(expected, want)
		}
	}
	binary := "lean/.lake/build/bin/godiff-path-filepath-split"
	if len(os.Args) == 2 {
		binary = os.Args[1]
	}
	cmd := exec.Command(binary)
	cmd.Stdin = &rows
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(lines) != len(expected) {
		panic(fmt.Sprintf("Lean produced %d rows", len(lines)))
	}
	for i, want := range expected {
		if lines[i] != wire(want) {
			panic(fmt.Sprintf("Lean row %d: got %q want %q", i, lines[i], wire(want)))
		}
	}
	volumeBlind := 0
	for _, c := range cases {
		i := strings.LastIndexAny(c.input, "/\\") + 1
		if (pair{c.input[:i], c.input[i:]}) != c.windows {
			volumeBlind++
		}
	}
	if volumeBlind == 0 {
		panic("volume-blind split survived")
	}
	raw := cases[1]
	if filepath.Clean(raw.input) == raw.unix.dir+raw.unix.file {
		panic("clean-before-split mutant survived")
	}
	fmt.Printf("filepath.Split passed %d native and %d Lean profile fixtures; volume-blind mutant failed %d\n", len(cases), len(expected), volumeBlind)
}
