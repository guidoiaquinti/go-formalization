package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type step struct{ Mode, Count, Writes, Value int }
type input struct {
	Tag, Len, Off, Cap, Flag int
	NilReader                bool
	Steps                    []step
}
type fixture struct {
	Input input
	Want  string
}

func main() {
	cases := []fixture{
		{input{0, 0, 0, 0, 0, false, []step{{1, 0, 0, 0}}}, "nilReceiver"},
		{input{1, 0, 0, 0, 3, true, nil}, "nilReader|0|slice|0|0||0.512|"},
		{input{1, 0, 0, 0, 3, false, []step{{0, 0, 1, 7}, {1, 3, 2, 9}}},
			"ok|3|slice|0|0|9.2,0.1|0.509|512:0.512;512:7.1,0.511"},
		{input{2, 1000, 980, 1200, 3, false, []step{{1, 1, 2, 7}}},
			"ok|1|slice|0|0|22.20,7.1|7.1,11.958,22.20,33.200|1180:11.960,22.20,33.200"},
		{input{2, 200, 100, 256, -1, false, []step{{2, 2, 3, 8}}},
			"readError|2|slice|0|0|22.100,8.2|8.1,0.665|668:0.668"},
		{input{2, 5, 5, 520, 2, false, []step{{1, 0, 0, 0}}},
			"ok|0|slice|0|0||11.5,33.515|520:11.5,33.515"},
		{input{1, 0, 0, 0, 3, false, []step{{1, -1, 2, 7}}},
			"negativeRead|0|slice|0|0||7.2,0.510|512:0.512"},
		{input{1, 0, 0, 0, 3, false, []step{{1, 513, 2, 7}}},
			"bounds|0|slice|0|0||7.2,0.510|512:0.512"},
		{input{2, 2, 1, 515, 3, false, []step{{3, 0, 2, 7}}},
			"readerPanic|0|slice|1|0|11.1,22.1|7.2,33.511|513:33.513"},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		i := c.Input
		nilReader := 0
		if i.NilReader {
			nilReader = 1
		}
		fmt.Fprint(&rows, i.Tag, " ", i.Len, " ", i.Off, " ", i.Cap, " ", i.Flag, " ", nilReader, " ", 2147483647)
		for _, step := range i.Steps {
			fmt.Fprint(&rows, " ", step.Mode, " ", step.Count, " ", step.Writes, " ", step.Value)
		}
		fmt.Fprintln(&rows)
	}
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(rows.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-buffer-read-from")
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
	checkWrappedEOF()
	fmt.Printf("bytes.Buffer.ReadFrom reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-read-from-control-")
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
	for _, name := range []string{"bridge_test.go", "buffer_read_from_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/bytes/buffer-read-from/testdata", name))
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
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanReadFrom$", "-v", "bytes")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_BUFFER_READ_FROM_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanReadFrom")) {
		log.Fatalf("private native control: %v\n%s", err, out)
	}
}

func encodeCases(cases []fixture) ([]byte, error) {
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	for _, c := range cases {
		if err := encoder.Encode(c); err != nil {
			return nil, err
		}
	}
	return data.Bytes(), nil
}

type readerFunc func([]byte) (int, error)

func (read readerFunc) Read(p []byte) (int, error) { return read(p) }

func checkWrappedEOF() {
	var destination bytes.Buffer
	wrapped := fmt.Errorf("wrapped EOF: %w", io.EOF)
	reader := readerFunc(func(p []byte) (int, error) {
		p[0] = 'Z'
		return 1, wrapped
	})
	n, err := destination.ReadFrom(reader)
	if n != 1 || err != wrapped || destination.String() != "Z" {
		log.Fatalf("wrapped EOF changed count, error identity or backing: %d %v %q", n, err, destination.String())
	}
}
