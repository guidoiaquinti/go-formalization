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
	Nil            bool
	Length, Offset int
	Argument       byte
	Backing        []byte
	Want           string
}

func main() {
	cells := []byte{90, 128, 255, 17}
	bounds := func(offset int) string {
		return strings.TrimSuffix(strings.Repeat(fmt.Sprintf("bounds:%d:90,128,255,17|", offset), 6), "|")
	}
	cases := []fixture{
		{Length: 2, Offset: 1, Argument: 255, Backing: cells, Want: "true:2:90,128,255,17|128:2:90,128,255,17|-128:2:90,128,255,17|ok:2:90,1,255,17|ok:2:90,255,255,17|ok:2:90,255,255,17"},
		{Length: 3, Offset: 2, Argument: 0, Backing: cells, Want: "true:3:90,128,255,17|255:3:90,128,255,17|-1:3:90,128,255,17|ok:3:90,128,0,17|ok:3:90,128,0,17|ok:3:90,128,0,17"},
		{Length: 1, Offset: 0, Argument: 128, Backing: []byte{0, 19}, Want: "false:1:0,19|0:1:0,19|0:1:0,19|ok:1:1,19|ok:1:128,19|ok:1:128,19"},
		{Length: 2, Offset: 2, Argument: 42, Backing: cells, Want: bounds(2)},
		{Length: 2, Offset: -1, Argument: 42, Backing: cells, Want: bounds(-1)},
		{Length: 0, Offset: 0, Argument: 42, Backing: cells, Want: bounds(0)},
		{Nil: true, Length: 2, Offset: 1, Argument: 42, Backing: cells, Want: "nil|nil|nil|nil|nil|nil"},
		{Length: 0, Offset: 0, Want: strings.TrimSuffix(strings.Repeat("bounds:0:|", 6), "|")},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		tag := 1
		if c.Nil {
			tag = 0
		}
		fmt.Fprint(&rows, tag, " ", c.Length, " ", c.Offset, " ", c.Argument)
		for _, b := range c.Backing {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
	}
	if len(os.Args) > 1 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(rows.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-binary-coder-byte")
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
	fmt.Printf("encoding/binary coder byte reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-coder-byte-control-")
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
	for _, name := range []string{"bridge_test.go", "coder_byte_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/encoding/binary/coder-byte/testdata", name))
		if err != nil {
			log.Fatal(err)
		}
		replacements[filepath.Join(runtime.GOROOT(), "src", "encoding", "binary", "golean_control_"+name)] = path
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(file, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanCoderByte$", "-v", "encoding/binary")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_BINARY_CODER_BYTE_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanCoderByte")) {
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
