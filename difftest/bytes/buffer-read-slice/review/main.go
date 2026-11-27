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
	Tag, Offset, Flag int
	Stored, Spare     []byte
	Count             int
}

type fixture struct {
	Input input
	Want  string
}

func main() {
	cases := []fixture{
		{input{Count: 0xff}, "nil-receiver"},
		{input{Tag: 1, Flag: 4}, "eof:nil:0:-1::"},
		{input{Tag: 2, Flag: 3}, "eof:slice:0:::0:-1::"},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: []byte{255, 0, 255, 64, 255}, Spare: []byte{7, 8}, Count: 255},
			"ok:slice:6:0,255:64,255,7,8:3:-1:255,0,255,64,255:7,8"},
		{input{Tag: 2, Offset: 1, Flag: 4, Stored: []byte{0, 17, 23}, Spare: []byte{0, 77}},
			"eof:slice:4:17,23:0,77:3:-1:0,17,23:0,77"},
		{input{Tag: 2, Offset: 3, Flag: 4, Stored: []byte{0, 17, 23}, Spare: []byte{0, 77}},
			"eof:slice:2::0,77:3:-1:0,17,23:0,77"},
		{input{Tag: 2, Offset: 1, Flag: 3, Stored: []byte{255, 7, 0}, Spare: []byte{4}},
			"ok:slice:3:7,0:4:3:-1:255,7,0:4"},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		i := c.Input
		fmt.Fprint(&rows, i.Tag, " ", i.Count)
		if i.Tag != 0 {
			fmt.Fprint(&rows, " ", len(i.Stored), " ", i.Offset, " ", i.Flag)
			for _, block := range [][]byte{i.Stored, i.Spare} {
				for _, b := range block {
					fmt.Fprint(&rows, " ", b)
				}
			}
		}
		fmt.Fprintln(&rows)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-buffer-read-slice")
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
	fmt.Printf("Buffer.readSlice reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-read-slice-control-")
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
	for _, name := range []string{"bridge_test.go", "read_slice_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/bytes/buffer-read-slice/testdata", name))
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
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanBufferReadSlice$", "-v", "bytes")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_BUFFER_READ_SLICE_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanBufferReadSlice")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
}
