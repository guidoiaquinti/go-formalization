package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func wire(kind, s string) string {
	fields := make([]string, len(s))
	for i := range len(s) {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return kind + strings.Join(fields, ",")
}

func main() {
	cases := []struct {
		profile, failed                int
		path, query, value, kind, want string
	}{
		{0, 1, "/a/./b/..", "", "getwd denied", "ok:", "/a"},
		{0, 1, "child", "", "getwd denied", "err:", "getwd denied"},
		{0, 0, "../\xff", "", "/working/inside", "ok:", "/working/\xff"},
		{2, 1, "#dev/../x", "", "getwd denied", "ok:", "x"},
		{2, 1, "relative", "", "getwd denied", "err:", "getwd denied"},
		{1, 0, "", ".", `C:\working\.`, "ok:", `C:\working`},
		{1, 1, "", ".", "full path failure", "err:", "full path failure"},
		{1, 0, `a/../b`, `a/../b`, `D:/one/../b`, "ok:", `D:\b`},
		{1, 1, "a\x00b", "a\x00b", "invalid raw path", "err:", "invalid raw path"},
		{1, 1, `C:\absolute`, `C:\absolute`, "full path still required", "err:", "full path still required"},
	}
	var model bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&model, "%d %d %d %d", c.profile, c.failed, len(c.path), len(c.query))
		for _, part := range []string{c.path, c.query, c.value} {
			for i := range len(part) {
				fmt.Fprintf(&model, " %d", part[i])
			}
		}
		fmt.Fprintln(&model)
	}
	for i, got := range check.LeanRows("godiff-path-filepath-abs-private", len(cases), &model) {
		if want := wire(cases[i].kind, cases[i].want); got != want {
			log.Fatalf("profile row %d: got %q want %q", i, got, want)
		}
	}
	if runtime.GOOS != "windows" && runtime.GOOS != "plan9" {
		cwd, err := os.MkdirTemp("", "golean-review-cwd-")
		must(err)
		defer os.RemoveAll(cwd)
		canonical, err := filepath.EvalSymlinks(cwd)
		must(err)
		f, err := os.CreateTemp("", "golean-review-private-abs-")
		must(err)
		defer os.Remove(f.Name())
		printLine := func(s string) { _, err := fmt.Fprintln(f, "x"+hex.EncodeToString([]byte(s))); must(err) }
		printLine(canonical)
		native := []struct{ path, want string }{
			{"", canonical}, {"./x/../a", canonical + "/a"}, {"/a/../b", "/b"}, {"\xff", canonical + "/\xff"}, {"#dev/./x", canonical + "/#dev/x"},
		}
		for _, c := range native {
			printLine(c.path)
			printLine("ok:" + c.want)
		}
		must(f.Close())
		must(os.Setenv("GOLEAN_PRIVATE_ABS_CASES", f.Name()))
		const data = "difftest/path/filepath/abs-private/testdata/"
		must(check.TestInPackage("path/filepath", "TestGoLeanPrivateAbs", data+"bridge_test.go", data+"abs_test.go"))
		fmt.Printf("private filepath.abs passed %d native fixed cases\n", len(native))
	}
	fmt.Printf("private filepath.abs passed %d independent Lean profile observations\n", len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
