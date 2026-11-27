package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type fixture struct {
	p, w, n, s int
	fields     [4]string
	want       string
}

func main() {
	cases := []fixture{
		{0, 0, 5, 6, [4]string{"q", "X", "Y", "bad"}, "ok:88"},
		{0, 7, 0, 0, [4]string{"q", "\x00\xff", "", ""}, "err:7:0,255"},
		{1, 7, 5, 0, [4]string{"q", "W", "N", ""}, "err:7:87"},
		{1, 0, 5, 0, [4]string{"q", "W", "N", ""}, "err:5:78"},
		{1, 0, 0, 0, [4]string{"q", "W", "N", ""}, "ok:78"},
		{2, 9, 8, 7, [4]string{"", "W", "N", "not a directory"}, "ok:46"},
		{2, 9, 8, 0, [4]string{"a/../b/./c", "W", "N", "ignored"}, "ok:98,47,99"},
		{2, 0, 0, 4, [4]string{"a", "", "", "\x00not a directory\xff"}, "notDirectory"},
		{2, 0, 0, -7, [4]string{"a", "", "", "not A directory"}, "err:-7:110,111,116,32,65,32,100,105,114,101,99,116,111,114,121"},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, c.p, " ", c.w, " ", c.n, " ", c.s)
		for _, field := range c.fields {
			fmt.Fprint(&rows, " ", len(field))
			for _, value := range []byte(field) {
				fmt.Fprint(&rows, " ", value)
			}
		}
		fmt.Fprintln(&rows)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-path-filepath-eval-symlinks-private")
	cmd.Stdin = &rows
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	results := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(results) != len(cases) {
		log.Fatalf("got %d Lean fixtures, want %d", len(results), len(cases))
	}
	for i, got := range results {
		if got != cases[i].want {
			log.Fatalf("fixture %d Lean=%q want %q", i, got, cases[i].want)
		}
	}
	if runtime.GOOS != "windows" && runtime.GOOS != "plan9" {
		checkNative()
	}
	fmt.Printf("evalSymlinks reviewer control: %d literal modeled results and native Unix delegation fixtures passed.\n", len(cases))
}

func checkNative() {
	dir, err := os.MkdirTemp("", "golean-eval-private-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	tree := filepath.Join(dir, "tree")
	if err := os.MkdirAll(filepath.Join(tree, "a", "b"), 0700); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tree, "a", "b", "file"), nil, 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Symlink("a/b", filepath.Join(tree, "jump")); err != nil {
		log.Fatal(err)
	}
	if err := os.Symlink("missing", filepath.Join(tree, "broken")); err != nil {
		log.Fatal(err)
	}
	var data bytes.Buffer
	for _, path := range []string{"", "jump/file", "jump/../b/file", "broken", "jump/file/child", "jump/\xff"} {
		fmt.Fprintln(&data, hex.EncodeToString([]byte(path)))
	}
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, data.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	replacements := make(map[string]string)
	for _, name := range []string{"bridge_test.go", "eval_symlinks_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/path/filepath/eval-symlinks-private/testdata", name))
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "path", "filepath", "golean_control_"+name)] = path
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanPrivateEvalSymlinks$", "-v", "path/filepath")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_EVAL_SYMLINK_CASES="+caseFile, "GOLEAN_EVAL_SYMLINK_TREE="+tree)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanPrivateEvalSymlinks")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
}
