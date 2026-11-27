package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	input, suffix     string
	cut               byte
	inputCap, wantCap int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []fixture{
		{"", "", 0, -1, 0},
		{"", "", 255, 0, 0},
		{"", "", 128, 9, 0},
		{"aaaa", "", 'a', 8, 0},
		{"aaaba", "ba", 'a', 9, 6},
		{"xa", "xa", 'a', 5, 5},
		{"é", "\xa9", 0xc3, 5, 4},
		{"\xff\xffx\xff", "x\xff", 255, 6, 4},
		{"\x00\x00x\x00", "x\x00", 0, 4, 2},
		{"\xffa", "\xffa", 254, 4, 4},
	}
	var modelInput, fixtures bytes.Buffer
	var oneRemoval, clampedSize, nilIfEmpty int
	for _, c := range cases {
		row := fmt.Appendf(nil, "%d %d", c.cut, max(0, c.inputCap-len(c.input)))
		for i := range len(c.input) {
			row = strconv.AppendUint(append(row, ' '), uint64(c.input[i]), 10)
		}
		fmt.Fprintln(&modelInput, string(row))
		fmt.Fprintf(&fixtures, "%d\t%d\t%x\t%d\t%x\n", c.cut, c.inputCap, c.input, c.wantCap, c.suffix)
		start := 0
		if len(c.input) > 0 && c.input[0] == c.cut {
			start++
		}
		if c.input[start:] != c.suffix {
			oneRemoval++
		}
		if c.suffix != "" && c.wantCap != len(c.suffix) {
			clampedSize++
		}
		if c.suffix == "" {
			nilIfEmpty++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-left-byte", len(cases), &modelInput) {
		if want := observe(cases[i]); row != want {
			return fmt.Errorf("case %d: Lean=%s want=%s", i, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-trim-left-byte-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixtures.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_TRIM_LEFT_BYTE_CASES", path); err != nil {
		return err
	}
	const data = "difftest/bytes/trim-left-byte/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanTrimLeftByte", data+"bridge_test.go", data+"trim_left_byte_test.go"); err != nil {
		return err
	}
	if oneRemoval != 4 || clampedSize != 5 || nilIfEmpty != 4 {
		return fmt.Errorf("single removal failed %d, clamped cap failed %d, nonnil on empty failed %d", oneRemoval, clampedSize, nilIfEmpty)
	}
	fmt.Printf("trimLeftByte passed %d fixtures; single removal failed %d, clamped cap failed %d, nonnil on empty failed %d\n",
		len(cases), oneRemoval, clampedSize, nilIfEmpty)
	return nil
}

func observe(c fixture) string {
	if c.suffix == "" {
		return "nil"
	}
	row := fmt.Appendf(nil, "%d:", c.wantCap)
	for i := range len(c.suffix) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(c.suffix[i]), 10)
	}
	return string(row)
}
