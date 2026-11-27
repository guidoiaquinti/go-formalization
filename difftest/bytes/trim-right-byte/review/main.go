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
	input, prefix string
	cut           byte
	capacity      int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []fixture{
		{"", "", 0, -1},
		{"", "", 255, 0},
		{"", "", 128, 9},
		{"aaaa", "", 'a', 8},
		{"abaaa", "ab", 'a', 9},
		{"ax", "ax", 'a', 5},
		{"é", "\xc3", 0xa9, 5},
		{"\xffx\xff\xff", "\xffx", 255, 6},
		{"\x00x\x00\x00", "\x00x", 0, 4},
		{"a\xff", "a\xff", 254, 4},
	}
	var modelInput, fixtures bytes.Buffer
	var oneRemoval, clampedSize, nilIfEmpty int
	for _, c := range cases {
		nilInput := 0
		if c.capacity < 0 {
			nilInput = 1
		}
		row := fmt.Appendf(nil, "%d %d %d", c.cut, max(0, c.capacity), nilInput)
		for i := range len(c.input) {
			row = strconv.AppendUint(append(row, ' '), uint64(c.input[i]), 10)
		}
		fmt.Fprintln(&modelInput, string(row))
		fmt.Fprintf(&fixtures, "%d\t%d\t%x\t%t\t%x\n", c.cut, c.capacity, c.input, nilInput == 1, c.prefix)
		end := len(c.input)
		if end > 0 && c.input[end-1] == c.cut {
			end--
		}
		if c.input[:end] != c.prefix {
			oneRemoval++
		}
		if c.capacity >= 0 && c.capacity != len(c.prefix) {
			clampedSize++
		}
		if c.capacity >= 0 && c.prefix == "" {
			nilIfEmpty++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-right-byte", len(cases), &modelInput) {
		if want := observe(cases[i]); row != want {
			return fmt.Errorf("case %d: Lean=%s want=%s", i, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-trim-right-byte-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixtures.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_TRIM_RIGHT_BYTE_CASES", path); err != nil {
		return err
	}
	const data = "difftest/bytes/trim-right-byte/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanTrimRightByte", data+"bridge_test.go", data+"trim_right_byte_test.go"); err != nil {
		return err
	}
	if oneRemoval != 4 || clampedSize != 8 || nilIfEmpty != 3 {
		return fmt.Errorf("single removal failed %d, clamped cap failed %d, nil on empty failed %d", oneRemoval, clampedSize, nilIfEmpty)
	}
	fmt.Printf("trimRightByte passed %d fixtures; single removal failed %d, clamped cap failed %d, nil on empty failed %d\n",
		len(cases), oneRemoval, clampedSize, nilIfEmpty)
	return nil
}

func observe(c fixture) string {
	if c.capacity < 0 {
		return "nil"
	}
	row := fmt.Appendf(nil, "%d:", c.capacity)
	for i := range len(c.prefix) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(c.prefix[i]), 10)
	}
	return string(row)
}
