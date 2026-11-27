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
		{"/../../x//./y", "/x/y", `\x\y`},
		{"../../x//./y", "../../x/y", `..\..\x\y`},
		{"foo:bar", "foo:bar", "foo:bar"},
		{"./foo:bar", "foo:bar", `.\foo:bar`},
		{"./foo:bar/..", ".", `.\.`},
		{"foo:bar/..", ".", "."},
		{`\.\??\..`, `\.\??\..`, `\.\`},
		{`\a\..\??\c:\x`, `\a\..\??\c:\x`, `\.\??\c:\x`},
		{"./foo:bar/../ab", "ab", `.\ab`},
		{"C:", "C:", "C:."},
		{"//host/share", "/host/share", `\\host\share`},
		{"//host/../x", "/x", `\x`},
		{"x\xff//../\x00", "\x00", "\x00"},
	}
	var input bytes.Buffer
	var wantRows []string
	for i, c := range cases {
		nativeWant := c.unix
		if runtime.GOOS == "windows" {
			nativeWant = c.windows
		}
		if got := filepath.Clean(c.input); got != nativeWant {
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
	binary := "lean/.lake/build/bin/godiff-path-filepath-clean"
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
	// A visible-prefix guard sees none of the stale bytes in these cases.
	staleGuardMutant := map[int]string{5: ".", 7: `\`, 9: "ab"}
	for index, mutant := range staleGuardMutant {
		if mutant == cases[index].windows {
			panic(fmt.Sprintf("visible-prefix guard mutant survived row %d", index))
		}
	}
	fmt.Printf("filepath.Clean passed %d pinned native and %d Lean profile witnesses\n", len(cases), len(wantRows))
}
