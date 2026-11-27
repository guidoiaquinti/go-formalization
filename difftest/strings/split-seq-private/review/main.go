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
	Input, Sep          []byte
	Save, First, Second int
	Want                string
}

func main() {
	maxInt := int(^uint(0) >> 1)
	cases := []fixture{
		{[]byte("ab--cd--"), []byte("--"), 1, 1, 2, "o:f97,98,45|o:f97,98,45;f99,100,45|o:f99,100,45;f|o:f"},
		{[]byte("ab--"), []byte("-"), -1, 1, 2, "o:f97|p:f97|p:|p:"},
		{[]byte("ab"), []byte("b"), maxInt, 1, 0, "p:|o:|p:|p:"},
		{[]byte("a|b"), []byte("|"), 2, 2, 1, "o:f97,124,98;f98|o:f98|o:f98|o:f98"},
		{[]byte{255, 195, 169, 226, 130}, nil, -maxInt - 1, 1, 2, "o:f255|o:f255;f195,169|o:f195,169;f226;f130|o:"},
		{[]byte("z"), []byte("-"), -maxInt - 1, 1, 2, "o:f122|o:f122|o:f122|o:f122"},
		{nil, nil, maxInt, 2, 0, "o:|o:|o:|o:"},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, len(c.Sep), " ", c.Save, " ", c.First, " ", c.Second)
		for _, block := range [][]byte{c.Sep, c.Input} {
			for _, b := range block {
				fmt.Fprint(&rows, " ", b)
			}
		}
		fmt.Fprintln(&rows)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-strings-split-seq-private")
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
	fmt.Printf("strings.splitSeq reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-split-seq-private-control-")
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
	for _, name := range []string{"bridge_test.go", "split_seq_private_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/strings/split-seq-private/testdata", name))
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "strings", "golean_control_"+name)] = path
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanSplitSeqPrivate$", "-v", "strings")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_SPLIT_SEQ_PRIVATE_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanSplitSeqPrivate")) {
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
