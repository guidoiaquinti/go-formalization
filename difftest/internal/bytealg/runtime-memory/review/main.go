package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type fixture struct {
	Left, Right                           []byte
	Length, Size, LeftOffset, RightOffset int
	Alias                                 bool
	Want                                  string
}

func main() {
	input := bytes.Repeat([]byte{0xa9}, 65)
	other := bytes.Clone(input)
	other[64] = 0
	cases := []fixture{
		{Length: 0, Want: "0:0:true|true"},
		{Length: 1, Right: []byte{55}, Size: 0, Want: "1:8:true|true"},
		{Length: 9, Left: []byte{0, 255, 64}, Right: []byte{0, 255, 99}, Size: 2, LeftOffset: 7, RightOffset: 15, Want: "9:16:true|true"},
		{Length: 257, Left: []byte{0, 255, 64}, Right: []byte{0, 255, 99}, Size: 3, Want: "257:288:true|false"},
		{Length: 512, Left: input, Right: other, Size: 64, LeftOffset: 31, RightOffset: 63, Want: "512:512:true|true"},
		{Length: 481, Left: input, Right: other, Size: 65, LeftOffset: 63, RightOffset: 31, Want: "481:512:true|false"},
		{Length: 8, Left: input, Right: input, Size: 65, LeftOffset: 19, Alias: true, Want: "8:8:true|true"},
		{Length: math.MinInt, Left: []byte{1}, Right: []byte{2}, Size: 1, Want: "panic|false"},
	}
	if strconv.IntSize == 64 {
		cases = append(cases, fixture{Length: math.MaxInt, Want: "panic|true"})
	}
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, c.Length, " ", c.Size, " ", len(c.Left))
		for _, block := range [][]byte{c.Left, c.Right} {
			for _, b := range block {
				fmt.Fprint(&rows, " ", b)
			}
		}
		fmt.Fprintln(&rows)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-internal-bytealg-runtime-memory")
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
	fullViewFailures, shortWindowFailures, capacityFailures := 0, 0, 0
	for _, c := range cases {
		equal := strings.HasSuffix(c.Want, "|true")
		if bytes.Equal(c.Left, c.Right) != equal {
			fullViewFailures++
		}
		if c.Size > 0 && bytes.Equal(c.Left[:c.Size-1], c.Right[:c.Size-1]) != equal {
			shortWindowFailures++
		}
		if c.Length >= 0 && c.Length <= 512 &&
			!strings.HasPrefix(c.Want, fmt.Sprintf("%d:%d:", c.Length, c.Length)) {
			capacityFailures++
		}
	}
	if fullViewFailures != 3 || shortWindowFailures != 3 || capacityFailures != 4 {
		log.Fatalf("mutation counts full=%d short=%d exact-capacity=%d", fullViewFailures, shortWindowFailures, capacityFailures)
	}
	checkNative(cases)
	fmt.Printf("internal/bytealg memory reviewer control: %d pinned private native and Lean fixtures passed.\n", len(cases))
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-bytealg-memory-control-")
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
	for _, name := range []string{"bridge_test.go", "runtime_memory_test.go"} {
		path, err := filepath.Abs(filepath.Join("difftest/internal/bytealg/runtime-memory/testdata", name))
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
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-overlay", file, "-run", "^TestGoLeanBytealgRuntimeMemory$", "-v", "internal/bytealg")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH,
		"GOLEAN_BYTEALG_RUNTIME_MEMORY_CASES="+caseFile)
	out, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte("--- PASS: TestGoLeanBytealgRuntimeMemory")) {
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
