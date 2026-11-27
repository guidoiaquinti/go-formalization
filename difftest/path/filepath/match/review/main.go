package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const match = "true:false"
const miss = "false:false"
const bad = "false:true"

func main() {
	fixtures := []struct{ pattern, name, unix, windows string }{
		{"", "", match, match},
		{"", "x", miss, miss},
		{"ab", "abab", miss, miss},
		{"*ab", "abab", match, match},
		{"a*[", "b", miss, miss},
		{"a[", "b", bad, bad},
		{"*a*[", "xx", miss, miss},
		{"*a*[", "xa", bad, bad},
		{"*", "x/y", miss, match},
		{"*", "x\\y", match, miss},
		{"*z", "x/z", miss, match},
		{"*z", "x\\z", match, miss},
		{"?", "/", miss, match},
		{"?", "\\", match, miss},
		{"[/]", "/", match, match},
		{"a\\*", "a*", match, miss},
		{"a\\*", "a\\suffix", miss, match},
		{"\\", "", bad, miss},
		{"?", "\xff", match, match},
		{"[�]", "\xff", match, match},
		{"[\xff]", "�", bad, bad},
		{"*[a/]*b", "a/b", miss, match},
	}
	var rows bytes.Buffer
	var expected []string
	for i, c := range fixtures {
		ok, err := filepath.Match(c.pattern, c.name)
		if err != nil && err != filepath.ErrBadPattern {
			panic(err)
		}
		nativeWant := c.unix
		if runtime.GOOS == "windows" {
			nativeWant = c.windows
		}
		if got := fmt.Sprintf("%t:%t", ok, err == filepath.ErrBadPattern); got != nativeWant {
			panic(fmt.Sprintf("native row %d: %q want %q", i, got, nativeWant))
		}
		for profile, want := range []string{c.unix, c.windows} {
			fmt.Fprintf(&rows, "%d %d", profile, len(c.pattern))
			for _, b := range []byte(c.pattern + c.name) {
				fmt.Fprintf(&rows, " %d", b)
			}
			fmt.Fprintln(&rows)
			expected = append(expected, want)
		}
	}
	binary := "lean/.lake/build/bin/godiff-path-filepath-match"
	if len(os.Args) == 2 {
		binary = os.Args[1]
	}
	cmd := exec.Command(binary)
	cmd.Stdin = &rows
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.TrimSuffix(string(stdout), "\n"), "\n")
	if len(lines) != len(expected) {
		panic(fmt.Sprintf("Lean produced %d rows", len(lines)))
	}
	for i, want := range expected {
		if lines[i] != want {
			panic(fmt.Sprintf("Lean row %d: %q want %q", i, lines[i], want))
		}
	}

	// These single-purpose mutants represent three plausible top-level omissions.
	prefix := fixtures[2]
	if strings.HasPrefix(prefix.name, prefix.pattern) == (prefix.unix == match) {
		panic("prefix-only acceptance survived")
	}
	eager := fixtures[4]
	if strings.Contains(eager.pattern, "[") && eager.unix == bad {
		panic("eager syntax validation survived")
	}
	star := fixtures[8]
	if star.pattern == "*" && star.unix == match {
		panic("separator-blind star survived")
	}
	greedy := fixtures[len(fixtures)-1]
	if ok, _ := filepath.Match("[a/]*b", greedy.name[1:]); runtime.GOOS != "windows" && (!ok || greedy.unix == match) {
		panic("global backtracking mutant survived")
	}
	fmt.Printf("filepath.Match passed %d pinned native and %d Lean profile fixtures\n", len(fixtures), len(expected))
}
