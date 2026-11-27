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

func main() {
	cases := []struct {
		header string
		fields [4]string
		want   string
	}{
		{"0 0 0 0", [4]string{"orig", "a/../b", "other", ""}, "ok:97,47,46,46,47,98"},
		{"1 7 9 0", [4]string{"orig", "w", "n", ""}, "err:7:119"},
		{"1 0 9 0", [4]string{"orig", "w", "n", ""}, "err:9:110"},
		{"1 0 0 0", [4]string{"orig", "w", "n", ""}, "ok:110"},
		{"2 0 0 7", [4]string{"a", "walk", "norm", "prefix not a directory suffix"}, "notDirectory"},
		{"2 0 0 7", [4]string{"a", "walk", "norm", "Not a directory"}, "err:7:78,111,116,32,97,32,100,105,114,101,99,116,111,114,121"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.header)
		for _, f := range c.fields {
			fmt.Fprintf(&input, " %d", len(f))
			for i := range len(f) {
				fmt.Fprintf(&input, " %d", f[i])
			}
		}
		fmt.Fprintln(&input)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-path-filepath-eval-symlinks")
	cmd.Stdin = &input
	out, err := cmd.Output()
	must(err)
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		panic("unexpected model row count")
	}
	for i, c := range cases {
		if rows[i] != c.want {
			panic(fmt.Sprintf("row %d got %s want %s", i, rows[i], c.want))
		}
	}
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		return
	}
	home, err := os.Getwd()
	must(err)
	dir, err := os.MkdirTemp("", "golean-review-eval-")
	must(err)
	defer os.RemoveAll(dir)
	must(os.Chdir(dir))
	defer func() { must(os.Chdir(home)) }()
	root, err := os.Getwd()
	must(err)
	must(os.Mkdir("child", 0700))
	must(os.WriteFile("target", nil, 0600))
	must(os.Symlink("../target", "child/relative"))
	must(os.Symlink("child/relative", "chain"))
	must(os.Symlink(filepath.Join(root, "target"), "absolute"))
	must(os.Symlink("cycle-b", "cycle-a"))
	must(os.Symlink("cycle-a", "cycle-b"))
	must(os.Symlink("absent", "broken"))
	for _, c := range []struct{ path, want string }{
		{"child/relative", "target"}, {"chain", "target"}, {"absolute", filepath.Join(root, "target")},
	} {
		got, err := filepath.EvalSymlinks(c.path)
		if err != nil || got != c.want {
			panic(fmt.Sprintf("%s got %q %v want %q", c.path, got, err, c.want))
		}
	}
	if got, err := filepath.EvalSymlinks("cycle-a"); got != "" || err == nil || err.Error() != "EvalSymlinks: too many links" {
		panic(fmt.Sprintf("cycle got %q %v", got, err))
	}
	if got, err := filepath.EvalSymlinks("broken"); got != "" || !os.IsNotExist(err) {
		panic(fmt.Sprintf("broken got %q %v", got, err))
	}
	fmt.Printf("EvalSymlinks passed %d independent model rows and five native symlink fixtures\n", len(cases))
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
