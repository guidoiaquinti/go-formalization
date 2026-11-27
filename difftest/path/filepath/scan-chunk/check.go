package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type output struct {
	star        bool
	chunk, rest string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []string{"", "*", "***", "a*b", "\\*a", "[*]*", "[[*]*", "\\[*", "a\\",
		"[*", "a]*", "\xff\x00*"}
	for b := range 256 {
		cases = append(cases, string([]byte{byte(b)}), string([]byte{'*', byte(b), '*', 0xff}))
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		r := rand.New(rand.NewSource(seed))
		for range 20_000 {
			raw := make([]byte, r.Intn(96))
			r.Read(raw)
			for j := 0; j < len(raw); j += 4 {
				raw[j] = "*[]\\"[r.Intn(4)]
			}
			cases = append(cases, string(raw))
		}
	}
	var lean, fixture bytes.Buffer
	for _, pattern := range cases {
		fmt.Fprint(&lean, 0)
		for i := range len(pattern) {
			fmt.Fprint(&lean, " ", pattern[i])
		}
		fmt.Fprintln(&lean)
		want := expected(pattern, runtime.GOOS == "windows")
		fmt.Fprintf(&fixture, "x%s %t x%s x%s\n", hex.EncodeToString([]byte(pattern)),
			want.star, hex.EncodeToString([]byte(want.chunk)), hex.EncodeToString([]byte(want.rest)))
	}
	for i, row := range check.LeanRows("godiff-path-filepath-scan-chunk", len(cases), &lean) {
		pattern := cases[i]
		want := wire(expected(pattern, false)) + "|" + wire(expected(pattern, true))
		if row != want {
			return fmt.Errorf("pattern %q: Lean %q, want %q", pattern, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-scan-chunk-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_FILEPATH_SCAN_CHUNK_CASES", path); err != nil {
		return err
	}
	const overlay = "difftest/path/filepath/scan-chunk/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanScanChunk", overlay+"bridge_test.go",
		overlay+"scan_chunk_test.go"); err != nil {
		return err
	}
	fmt.Printf("path/filepath.scanChunk: %d boundaries and %d seeded rows agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(pattern string, windows bool) output {
	body := strings.TrimLeft(pattern, "*")
	star := len(body) != len(pattern)
	control := make([]bool, len(body))
	if !windows {
		skip := false
		for i := range len(body) {
			control[i] = skip
			skip = !skip && body[i] == '\\'
		}
	}
	lastBracket := byte(']')
	for i := range len(body) {
		if control[i] {
			continue
		}
		switch body[i] {
		case '[', ']':
			lastBracket = body[i]
		case '*':
			if lastBracket == ']' {
				return output{star, body[:i], body[i:]}
			}
		}
	}
	return output{star, body, ""}
}

func wire(result output) string {
	encode := func(s string) string {
		fields := make([]string, len(s))
		for i := range len(s) {
			fields[i] = strconv.Itoa(int(s[i]))
		}
		return strings.Join(fields, ",")
	}
	return strconv.FormatBool(result.star) + ":" + encode(result.chunk) + ":" + encode(result.rest)
}
