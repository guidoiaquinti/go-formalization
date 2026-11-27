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
	"strings"
)

type input struct {
	NonNil              bool
	Raw, Spare, Sep     []byte
	Save, First, Second int
}
type fixture struct {
	Input input
	Want  string
}

func main() {
	cases := []fixture{
		{input{true, []byte("a\n\nb"), []byte("X"), []byte("\n"), 0, 1, 1},
			"ok/slice:2:97,10|ok/slice:1:10|ok/slice:1:98|ok/#ok/slice:1:97|ok/slice:1:97|ok/slice:1:97;slice:0:;slice:1:98|ok/slice:1:98"},
		{input{true, []byte("ab--"), []byte("c"), []byte("-"), -1, 0, 0},
			"ok/|ok/|ok/slice:4:97,98,45,45|ok/#ok/|ok/|panic/slice:1:97|panic/"},
		{input{true, []byte("a-"), []byte("XY"), []byte("-"), 3, 1, 2},
			"ok/slice:2:97,45|ok/|ok/|ok/#ok/slice:4:97,45,88,89|ok/slice:4:97,45,88,89;slice:0:|ok/slice:0:|ok/slice:0:"},
		{input{true, []byte{255, 195, 169}, nil, nil, 42, 1, 2},
			"ok/slice:3:255,195,169|ok/|ok/|ok/#ok/slice:1:255|ok/slice:1:255;slice:2:195,169|ok/slice:2:195,169|ok/"},
		{input{false, nil, nil, []byte("x"), 1, 1, 1},
			"ok/|ok/|ok/|ok/#ok/nil:0:|ok/nil:0:|ok/nil:0:|ok/nil:0:"},
		{input{true, []byte("\n"), []byte("X"), []byte("\n"), 1, 1, 1},
			"ok/slice:1:10|ok/|ok/|ok/#ok/slice:1:10|ok/slice:1:10|ok/slice:1:10;slice:0:|ok/slice:0:"},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		i := c.Input
		tag := 0
		if i.NonNil {
			tag = 1
		}
		fmt.Fprint(&rows, tag, " ", len(i.Raw), " ", len(i.Spare), " ", len(i.Sep), " ", i.Save, " ", i.First, " ", i.Second)
		for _, block := range [][]byte{i.Raw, i.Spare, i.Sep} {
			for _, b := range block {
				fmt.Fprint(&rows, " ", b)
			}
		}
		fmt.Fprintln(&rows)
	}
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(rows.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-configured-iter")
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
		if got != cases[i].Want {
			log.Fatalf("fixture %d Lean=%q want %q", i, got, cases[i].Want)
		}
	}
	checkNative(cases)
	fmt.Printf("bytes configured iterator reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-bytes-configured-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases.json")
	data, err := json.Marshal(cases)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(caseFile, data, 0600); err != nil {
		log.Fatal(err)
	}
	replacements := make(map[string]string)
	for _, name := range []string{"bridge_test.go", "configured_iter_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/bytes/configured-iter/testdata", name))
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "bytes", "golean_control_"+name)] = path
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanConfiguredIter$", "-v", "bytes")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_CONFIGURED_ITER_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanConfiguredIter")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
}
