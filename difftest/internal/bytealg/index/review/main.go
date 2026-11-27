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

type fixture struct {
	Input, Needle                   []byte
	InputOffset, NeedleOffset, Want int
}

func main() {
	long := bytes.Repeat([]byte{0xef}, 31)
	if runtime.GOARCH == "arm64" {
		long = append(long, 0xff)
	}
	cases := []fixture{
		{[]byte("aaaa"), []byte("aa"), 3, 15, 0},
		{[]byte("abcab"), []byte("bc"), 15, 31, 1},
		{[]byte{0xff, 0x80, 0, 0xff, 0x80}, []byte{0xff, 0x80}, 1, 7, 0},
		{[]byte{195, 169, 128}, []byte{169, 128}, 31, 3, 1},
		{[]byte("abc"), []byte("ca"), 8, 17, -1},
		{append([]byte{1}, long...), long, 31, 15, 1},
		{bytes.Repeat([]byte{2}, len(long)+17), long, 17, 30, -1},
		{bytes.Clone(long), long, 11, 29, 0},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, len(c.Needle))
		for _, block := range [][]byte{c.Needle, c.Input} {
			for _, b := range block {
				fmt.Fprint(&rows, " ", b)
			}
		}
		fmt.Fprintln(&rows)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-internal-bytealg-index")
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
		want := fmt.Sprintf("%d,%d", cases[i].Want, cases[i].Want)
		if got != want {
			log.Fatalf("fixture %d Lean=%q want %q", i, got, want)
		}
	}
	reverseFailures, finalFailures, beyondFailures := 0, 0, 0
	for _, c := range cases {
		if bytes.LastIndex(c.Input, c.Needle) != c.Want {
			reverseFailures++
		}
		if bytes.Index(c.Input[:len(c.Input)-1], c.Needle) != c.Want {
			finalFailures++
		}
		extended := append(bytes.Clone(c.Input), c.Needle...)
		if bytes.Index(extended, c.Needle) != c.Want {
			beyondFailures++
		}
	}
	if reverseFailures != 2 || finalFailures != 3 || beyondFailures != 2 {
		log.Fatalf("mutation counts reverse=%d final=%d beyond=%d", reverseFailures, finalFailures, beyondFailures)
	}
	if runtime.GOARCH == "arm64" || runtime.GOARCH == "amd64" {
		checkNative(cases)
	}
	fmt.Printf("internal/bytealg index reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-bytealg-index-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases.json")
	data, err := encodeCases(cases)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(caseFile, data, 0600); err != nil {
		log.Fatal(err)
	}
	replacements := make(map[string]string)
	for _, name := range []string{"bridge_test.go", "index_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/internal/bytealg/index/testdata", name))
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "internal", "bytealg", "golean_control_"+name)] = path
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanBytealgIndex$", "-v", "internal/bytealg")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_BYTEALG_INDEX_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanBytealgIndex")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
}

func encodeCases(cases []fixture) ([]byte, error) {
	var result bytes.Buffer
	for _, c := range cases {
		if err := json.NewEncoder(&result).Encode(c); err != nil {
			return nil, err
		}
	}
	return result.Bytes(), nil
}
