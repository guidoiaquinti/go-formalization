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
	Bytes               []byte
	Target              byte
	Offset, First, Last int
}

func main() {
	vector := bytes.Repeat([]byte{9}, 129)
	vector[15], vector[16], vector[63], vector[64], vector[128] = 255, 255, 255, 255, 255
	cases := []fixture{
		{nil, 0, -1, -1, -1},
		{[]byte{}, 255, 31, -1, -1},
		{[]byte{0, 255, 0, 254, 0}, 0, 7, 0, 4},
		{[]byte{192, 175, 255, 239, 191, 189, 255}, 255, 15, 2, 6},
		{[]byte{195, 169, 140, 169}, 169, 5, 1, 3},
		{[]byte{1, 2, 3}, 4, 23, -1, -1},
		{vector, 255, 1, 15, 128},
		{vector[17:128], 255, 13, 46, 47},
		{bytes.Repeat([]byte{3}, 257), 4, 29, -1, -1},
		{[]byte{2, 3}, 3, 3, 1, 1},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, c.Target)
		for _, b := range c.Bytes {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-internal-bytealg-byte-index")
	cmd.Stdin = &rows
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	results := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(results) != len(cases) {
		log.Fatalf("got %d Lean fixtures, want %d", len(results), len(cases))
	}
	firstLastDifferences, missing := 0, 0
	for i, got := range results {
		c := cases[i]
		want := fmt.Sprintf("%d,%d,%d,%d", c.First, c.First, c.Last, c.Last)
		if got != want {
			log.Fatalf("fixture %d Lean=%q want %q", i, got, want)
		}
		if c.First != c.Last {
			firstLastDifferences++
		}
		if c.First < 0 {
			missing++
		}
	}
	if firstLastDifferences != 5 || missing != 4 {
		log.Fatalf("mutation count first/last=%d missing=%d", firstLastDifferences, missing)
	}
	checkNative(cases)
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
	for _, name := range []string{"bridge_test.go", "byte_index_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/internal/bytealg/byte-index/testdata", name))
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
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanBytealgByteIndex$", "-v", "internal/bytealg")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_BYTEALG_BYTE_INDEX_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanBytealgByteIndex")) {
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
