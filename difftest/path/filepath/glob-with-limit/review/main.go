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
	"strconv"
	"strings"
)

type fixture struct {
	depth           int
	pattern, status string
	paths           []string
}

func main() {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		fmt.Println("globWithLimit reviewer fixtures require Unix path semantics")
		return
	}
	cases := []fixture{
		{9999, "./plain", "", []string{"./plain"}},
		{10000, "plain", "bad", nil},
		{10001, "plain", "", []string{"plain"}},
		{-1, "missing", "nil", nil},
		{9999, "alpha/*.txt", "", []string{"alpha/a.txt"}},
		{9999, "*/*.txt", "bad", nil},
		{9998, "*/*.txt", "", []string{"alpha/a.txt", "beta/b.txt"}},
		{9998, "*/missing", "nil", nil},
		{0, "plain/*", "nil", nil},
		{0, "*/[", "bad", nil},
		{0, "./alpha//a.txt", "", []string{"./alpha//a.txt"}},
		{0, "[ab]*/?.txt", "", []string{"alpha/a.txt", "beta/b.txt"}},
	}
	var input, native bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, "0 ", c.depth)
		for _, b := range []byte(c.pattern) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		fmt.Fprintf(&native, "%d x%s %s\n", c.depth, hex.EncodeToString([]byte(c.pattern)), c.expected())
	}
	model := exec.Command("lean/.lake/build/bin/godiff-path-filepath-glob-with-limit")
	model.Stdin = &input
	out, err := model.Output()
	if err != nil {
		log.Fatal(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		log.Fatalf("got %d Lean rows, want %d", len(rows), len(cases))
	}
	for i, c := range cases {
		if want := c.expected(); rows[i] != want {
			log.Fatalf("fixture %d: Lean=%q want %q", i, rows[i], want)
		}
	}
	checkNative(native.Bytes())
	fmt.Printf("globWithLimit reviewer control: %d literal private native and Lean cases passed.\n", len(cases))
}

func (c fixture) expected() string {
	if c.status != "" {
		return c.status
	}
	parts := make([]string, len(c.paths))
	for i, path := range c.paths {
		vals := make([]string, len(path))
		for j := range len(path) {
			vals[j] = strconv.Itoa(int(path[j]))
		}
		parts[i] = strings.Join(vals, ",")
	}
	return strconv.Itoa(len(parts)) + ":" + strings.Join(parts, ";")
}

func checkNative(fixture []byte) {
	scratch, err := os.MkdirTemp("", "golean-glob-limit-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(scratch)
	root := filepath.Join(scratch, "tree")
	for _, name := range []string{"alpha", "beta"} {
		if err := os.MkdirAll(filepath.Join(root, name), 0700); err != nil {
			log.Fatal(err)
		}
	}
	for _, name := range []string{"plain", "alpha/a.txt", "alpha/z.go", "beta/b.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0600); err != nil {
			log.Fatal(err)
		}
	}
	cases := filepath.Join(scratch, "cases")
	if err := os.WriteFile(cases, fixture, 0600); err != nil {
		log.Fatal(err)
	}
	data := "difftest/path/filepath/glob-with-limit/testdata"
	replacements := make(map[string]string)
	for _, name := range []string{"bridge_test.go", "glob_limit_test.go"} {
		path, err := filepath.Abs(filepath.Join(data, name))
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "path", "filepath", "golean_glob_control_"+name)] = path
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(scratch, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanGlobWithLimit$", "-v", "path/filepath")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(scratch, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_GLOB_LIMIT_DIR="+root, "GOLEAN_GLOB_LIMIT_CASES="+cases)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanGlobWithLimit")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
}
