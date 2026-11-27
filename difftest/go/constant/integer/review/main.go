package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type fixture struct {
	Tag     int
	Integer string
	Bytes   []byte
	Want    string
}

func main() {
	const maxU = "18446744073709551615"
	const minS = "-9223372036854775808"
	const sign = "9223372036854775808"
	high := []byte{0, 0, 0, 0, 0, 0, 0, 128}
	cases := []fixture{
		{0, "-1", []byte{1, 2, 0}, "small:-1|large:" + maxU + "|small:-1|small:513|-1:true|" + maxU + ":false|1|1"},
		{0, sign, high, "small:" + minS + "|large:" + sign + "|large:" + sign + "|large:" + sign + "|" + minS + ":false|" + sign + ":true|64|0,0,0,0,0,0,0,128"},
		{0, "-18446744073709551621", nil, "small:-5|large:18446744073709551611|large:-18446744073709551621|small:0|-5:false|5:false|65|5,0,0,0,0,0,0,0,1"},
		{0, "18446744073709551621", nil, "small:5|small:5|large:18446744073709551621|small:0|5:false|5:false|65|5,0,0,0,0,0,0,0,1"},
		{0, minS, []byte{0, 1}, "small:" + minS + "|large:" + sign + "|small:" + minS + "|small:256|" + minS + ":true|" + sign + ":false|64|0,0,0,0,0,0,0,128"},
		{1, "17", []byte{1, 2, 0}, "small:17|small:17|small:17|small:513|513:true|513:true|10|1,2"},
		{1, "0", []byte{0, 0}, "small:0|small:0|small:0|small:0|0:true|0:true|0|"},
		{2, "0", nil, "small:0|small:0|small:0|small:0|0:false|0:false|0|panic"},
		{3, "123", nil, "small:123|small:123|small:123|small:0|panic|panic|panic|panic"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.Tag, " ", c.Integer)
		for _, b := range c.Bytes {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-go-constant-integer")
	cmd.Stdin = &input
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	got := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(got) != len(cases) {
		log.Fatalf("got %d Lean rows, want %d", len(got), len(cases))
	}
	var rejected [8]int
	for i, c := range cases {
		if got[i] != c.Want {
			log.Fatalf("row %d: Lean %q want %q", i, got[i], c.Want)
		}
		parts := strings.Split(got[i], "|")
		if len(parts) != 8 {
			log.Fatalf("invalid observation at row %d", i)
		}
		for k, wrong := range mutants(c, parts) {
			if parts[k] != wrong {
				rejected[k]++
			}
		}
	}
	wantRejected := [8]int{4, 4, 3, 4, 3, 1, 3, 5}
	if rejected != wantRejected {
		log.Fatalf("mutation counts %v want %v", rejected, wantRejected)
	}
	checkNative(cases)
	fmt.Printf("go/constant integer review: %d pinned native and Lean rows, mutant rejection counts %v\n", len(cases), rejected)
}

func mutants(c fixture, fields []string) []string {
	x, ok := new(big.Int).SetString(c.Integer, 10)
	if !ok {
		log.Fatal("invalid fixture integer")
	}
	low := x.Int64()
	out := slices.Clone(fields)
	if low < 0 {
		out[0] = "large:" + strconv.FormatUint(uint64(low), 10)
	}
	out[1] = "small:" + strconv.FormatInt(low, 10)
	out[2] = "small:" + strconv.FormatInt(low, 10)
	bigEndian := new(big.Int).SetBytes(c.Bytes)
	repr := "large:"
	if bigEndian.IsInt64() {
		repr = "small:"
	}
	out[3] = repr + bigEndian.String()
	if c.Tag == 0 && !x.IsInt64() {
		out[4] = "0:false"
		if x.Sign() < 0 {
			out[5] = strconv.FormatUint(uint64(low), 10) + ":false"
		}
	}
	if c.Tag == 0 && x.Sign() < 0 {
		out[6] = strconv.Itoa(x.BitLen() + 1)
	}
	if fields[7] != "" && fields[7] != "panic" {
		cells := strings.Split(fields[7], ",")
		slices.Reverse(cells)
		out[7] = strings.Join(cells, ",")
	}
	return out
}

func checkNative(cases []fixture) {
	dir, err := os.MkdirTemp("", "golean-constant-integer-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	var rows bytes.Buffer
	encoder := json.NewEncoder(&rows)
	for _, c := range cases {
		if err := encoder.Encode(c); err != nil {
			log.Fatal(err)
		}
	}
	casesPath := filepath.Join(dir, "rows.jsonl")
	if err := os.WriteFile(casesPath, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	fixturePath, err := filepath.Abs("difftest/go/constant/integer/testdata/integer_test.go")
	if err != nil {
		log.Fatal(err)
	}
	overlay, err := json.Marshal(struct{ Replace map[string]string }{
		map[string]string{filepath.Join(runtime.GOROOT(), "src", "go", "constant", "golean_integer_review_test.go"): fixturePath}})
	if err != nil {
		log.Fatal(err)
	}
	overlayPath := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(overlayPath, overlay, 0600); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "-v", "-overlay", overlayPath,
		"-run", "^TestGoLeanConstantInteger$", "go/constant")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOROOT="+runtime.GOROOT(), "GODEBUG=goindex=0",
		"GOMODCACHE="+filepath.Join(dir, "modcache"), "GOLEAN_CONSTANT_INTEGER_CASES="+casesPath)
	output, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(output, []byte("--- PASS: TestGoLeanConstantInteger")) {
		log.Fatalf("native controls: %v\n%s", err, output)
	}
}
