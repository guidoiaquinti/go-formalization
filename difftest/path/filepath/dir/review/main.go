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

func wire(s string) string {
	fields := make([]string, len(s))
	for i := range len(s) {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return "bytes:" + strings.Join(fields, ",")
}

func main() {
	cases := []struct{ input, unix, windows string }{
		{"", ".", "."},
		{"a/../leaf", ".", "."},
		{"a/../leaf/", "leaf", "leaf"},
		{"a/./leaf", "a", "a"},
		{"foo:bar/leaf", "foo:bar", "foo:bar"},
		{"./foo:bar/leaf", "foo:bar", `.\foo:bar`},
		{"./foo:bar/../leaf", ".", `.\.`},
		{"foo:bar/../leaf", ".", "."},
		{`\.\??\..\leaf`, ".", `\.\`},
		{"C:", ".", "C:."},
		{"C:/leaf", "C:", `C:\`},
		{"//h/share", "/h", `\\h\share`},
		{"//h/share/./leaf", "/h/share", `\\h\share\`},
		{"x\xff/leaf", "x\xff", "x\xff"},
	}
	var input bytes.Buffer
	var wantRows []string
	for i, c := range cases {
		nativeWant := c.unix
		if runtime.GOOS == "windows" {
			nativeWant = c.windows
		}
		if got := filepath.Dir(c.input); got != nativeWant {
			panic(fmt.Sprintf("native row %d: got %q want %q", i, got, nativeWant))
		}
		for profile, want := range []string{c.unix, c.windows} {
			fmt.Fprintf(&input, "%d", profile)
			for _, b := range []byte(c.input) {
				fmt.Fprintf(&input, " %d", b)
			}
			fmt.Fprintln(&input)
			wantRows = append(wantRows, wire(want))
		}
	}
	binary := "lean/.lake/build/bin/godiff-path-filepath-dir"
	if len(os.Args) == 2 {
		binary = os.Args[1]
	}
	cmd := exec.Command(binary)
	cmd.Stdin = &input
	raw, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	actual := strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n")
	if len(actual) != len(wantRows) {
		panic(fmt.Sprintf("Lean produced %d rows", len(actual)))
	}
	for i, want := range wantRows {
		if actual[i] != want {
			panic(fmt.Sprintf("Lean row %d: got %q want %q", i, actual[i], want))
		}
	}
	// Cleaning the whole path first removes the trailing slash and changes which
	// component is treated as the file.
	terminal := cases[2]
	if filepath.Dir(filepath.Clean(terminal.input)) == terminal.unix {
		panic("clean-before-split mutant survived")
	}
	unc := cases[11]
	if strings.ReplaceAll(unc.input, "/", "\\")+"." == unc.windows {
		panic("UNC dot-retention mutant survived")
	}
	fmt.Printf("filepath.Dir passed %d pinned native and %d Lean profile witnesses\n", len(cases), len(wantRows))
}
