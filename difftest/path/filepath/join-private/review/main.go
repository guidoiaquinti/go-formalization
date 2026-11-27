package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
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
	var native bytes.Buffer
	for _, c := range cases {
		nativeWant := c.unix
		if runtime.GOOS == "windows" {
			nativeWant = c.windows
		}
		fmt.Fprintf(&native, "x%s", hex.EncodeToString([]byte(nativeWant)))
		for _, part := range c.parts {
			fmt.Fprintf(&native, " x%s", hex.EncodeToString([]byte(part)))
		}
		fmt.Fprintln(&native)
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
	f, err := os.CreateTemp("", "golean-review-private-join-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(native.Bytes()); err != nil {
		panic(err)
	}
	if err = f.Close(); err != nil {
		panic(err)
	}
	if err = os.Setenv("GOLEAN_PRIVATE_JOIN_CASES", f.Name()); err != nil {
		panic(err)
	}
	const data = "difftest/path/filepath/join-private/testdata/"
	if err = check.TestInPackage("path/filepath", "TestGoLeanPrivateJoin", data+"bridge_test.go", data+"join_test.go"); err != nil {
		panic(err)
	}
	binary := "lean/.lake/build/bin/godiff-path-filepath-join-private"
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
	fmt.Printf("private filepath.join passed %d pinned native and %d Lean profile witnesses\n", len(cases), len(wantRows))
}
