package main

import (
	"bytes"
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

func main() {
	cases := []struct {
		profile, failed           int
		path, cwd, status, result string
	}{
		{0, 1, "/./q/../\xff", "\x00\xff", "ok:", "/\xff"},
		{0, 1, "../q\xff", "e\x00\xff", "err:", "e\x00\xff"},
		{0, 0, "../\xff", "/base/inner", "ok:", "/base/\xff"},
		{1, 1, `C:\a\..\p`, "\x00", "ok:", `C:\p`},
		{1, 1, `C:child`, "\xff", "err:", "\xff"},
		{2, 1, "#c/child/../\xff", "\x00", "ok:", "#c/\xff"},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, c.profile, " ", c.failed, " ", len(c.path))
		for _, b := range []byte(c.path + c.cwd) {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
	}
	model := exec.Command("lean/.lake/build/bin/godiff-path-filepath-unix-abs")
	model.Stdin = &rows
	out, err := model.Output()
	if err != nil {
		log.Fatal(err)
	}
	results := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(results) != len(cases) {
		log.Fatalf("unexpected row count %d", len(results))
	}
	for i, c := range cases {
		parts := make([]string, len(c.result))
		for j := range len(c.result) {
			parts[j] = strconv.Itoa(int(c.result[j]))
		}
		if want := c.status + strings.Join(parts, ","); results[i] != want {
			log.Fatalf("row %d: Lean=%q want %q", i, results[i], want)
		}
	}
	fmt.Printf("unixAbs reviewer control: %d literal Lean fixtures passed.\n", len(cases))
	if runtime.GOOS == "darwin" {
		testFile := "difftest/path/filepath/unix-abs/review/testdata/deep_cwd_test.go"
		if len(os.Args) == 2 {
			testFile = os.Args[1]
		}
		checkDeepCwd(testFile)
	}
}

func checkDeepCwd(testFile string) {
	scratch, err := os.MkdirTemp("", "golean-unix-abs-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(scratch)
	replacements := make(map[string]string)
	for _, source := range []string{"difftest/path/filepath/unix-abs/testdata/bridge_test.go", testFile} {
		path, err := filepath.Abs(source)
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "path", "filepath", "golean_control_"+filepath.Base(path))] = path
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(scratch, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanUnixAbsDeepCwd$", "-v", "path/filepath")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(scratch, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanUnixAbsDeepCwd")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
	fmt.Print(string(out))
}
