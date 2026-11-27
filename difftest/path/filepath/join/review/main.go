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
	cases := []struct {
		parts         []string
		unix, windows string
	}{
		{nil, "", ""},
		{[]string{"", ""}, "", ""},
		{[]string{"", "a", "", "../x"}, "x", "x"},
		{[]string{"C:", "leaf"}, "C:/leaf", "C:leaf"},
		{[]string{"C:", `\leaf`}, `C:/\leaf`, `C:\leaf`},
		{[]string{"/", "/host/share"}, "/host/share", `\host\share`},
		{[]string{"//", "host", "share"}, "/host/share", `\\host\share`},
		{[]string{`\`, "??", "x"}, `\/??/x`, `\.\??\x`},
		{[]string{`\`, "??z", "x"}, `\/??z/x`, `\??z\x`},
		{[]string{"/", "??z", "x"}, "/??z/x", `\.\??z\x`},
		{[]string{"./foo:bar", ".."}, ".", `.\.`},
		{[]string{"foo:bar", ".."}, ".", "."},
		{[]string{"x\xff", "..", "\x00"}, "\x00", "\x00"},
	}
	var input bytes.Buffer
	var wantRows []string
	for i, c := range cases {
		nativeWant := c.unix
		if runtime.GOOS == "windows" {
			nativeWant = c.windows
		}
		if got := filepath.Join(c.parts...); got != nativeWant {
			panic(fmt.Sprintf("native row %d: got %q want %q", i, got, nativeWant))
		}
		for profile, want := range []string{c.unix, c.windows} {
			fmt.Fprintf(&input, "%d %d", profile, len(c.parts))
			for _, part := range c.parts {
				fmt.Fprintf(&input, " %d", len(part))
				for _, b := range []byte(part) {
					fmt.Fprintf(&input, " %d", b)
				}
			}
			fmt.Fprintln(&input)
			wantRows = append(wantRows, wire(want))
		}
	}
	binary := "lean/.lake/build/bin/godiff-path-filepath-join"
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
	drive := cases[3]
	if strings.Join(drive.parts, `\`) == drive.windows {
		panic("unconditional separator mutant survived")
	}
	deliberatelyNoGuard := cases[7]
	if deliberatelyNoGuard.parts[0]+strings.Join(deliberatelyNoGuard.parts[1:], `\`) == deliberatelyNoGuard.windows {
		panic("omitted root device guard survived")
	}
	fmt.Printf("filepath.Join passed %d pinned native and %d Lean profile witnesses\n", len(cases), len(wantRows))
}
