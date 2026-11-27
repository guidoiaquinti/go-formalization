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
	Tag, Previous        int
	Cursor, Offset       int64
	Sharing, Size        int
	Raw, Spare, Separate []byte
}
type fixture struct {
	Input input
	Want  string
}

func main() {
	raw, spare := []byte{1, 2, 3, 4}, []byte{9}
	cases := []fixture{
		{input{Tag: 2, Previous: 7, Offset: 2, Sharing: 1, Size: 4, Raw: raw, Spare: spare},
			"4:ok:1,2,3,4:slice:4:-1:1,1,2,3:4|2:EOF:3,4,4,9:slice:0:7:1,3,4,4:9"},
		{input{Tag: 2, Previous: 7, Cursor: 2, Offset: 1, Sharing: 0, Size: 3, Raw: raw, Spare: spare},
			"2:ok:3,4,3:slice:4:-1:3,4,3,4:9|3:ok:2,3,4:slice:2:7:2,3,4,4:9"},
		{input{Tag: 2, Previous: 7, Cursor: 3, Offset: 3, Sharing: -1, Size: 2, Raw: raw, Spare: spare, Separate: []byte{88, 99}},
			"1:ok:4,99:slice:4:-1:1,2,3,4:9|1:EOF:4,99:slice:3:7:1,2,3,4:9"},
		{input{Tag: 2, Previous: 7, Cursor: 4, Offset: 4, Sharing: -1, Raw: raw, Spare: spare},
			"0:EOF::slice:4:7:1,2,3,4:9|0:EOF::slice:4:7:1,2,3,4:9"},
		{input{Tag: 2, Previous: 7, Cursor: 3, Offset: 3, Sharing: -1, Raw: raw, Spare: spare},
			"0:ok::slice:3:-1:1,2,3,4:9|0:ok::slice:3:7:1,2,3,4:9"},
		{input{Tag: 0, Offset: -1, Sharing: -1, Size: 1, Separate: []byte{8}},
			"panic|0:bytes.Reader.ReadAt: negative offset:8:receiver-nil"},
		{input{Tag: 0, Offset: 0, Sharing: -1}, "panic|panic"},
		{input{Tag: 1, Previous: 2, Cursor: 9223372036854775807, Offset: 9223372036854775807, Sharing: -1},
			"0:EOF::nil:9223372036854775807:2::|0:EOF::nil:9223372036854775807:2::"},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		i := c.Input
		fmt.Fprint(&rows, i.Tag, " ", len(i.Raw), " ", len(i.Spare), " ", i.Cursor, " ", i.Previous, " ", i.Offset, " ", i.Sharing, " ", i.Size)
		for _, block := range [][]byte{i.Raw, i.Spare, i.Separate} {
			for _, b := range block {
				fmt.Fprint(&rows, " ", b)
			}
		}
		fmt.Fprintln(&rows)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-reader-read")
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
	fmt.Printf("bytes.Reader copy reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-reader-copy-control-")
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
	for _, name := range []string{"bridge_test.go", "reader_read_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/bytes/reader-read/testdata", name))
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
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanReaderCopy$", "-v", "bytes")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_READER_COPY_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanReaderCopy")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
}
