package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	text, cutset, prefix string
	capacity             int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []fixture{
		{"", "", "", -1},
		{"", "�", "", 0},
		{"", "a", "", 5},
		{"aaaa", "a", "", 8},
		{"éxé🚀", "é🚀", "éx", 13},
		{"\xffx�\x80", "�", "\xffx", 8},
		{"\xffx\xfe", "a", "\xffx\xfe", 8},
		{"aé", "\xa9", "aé", 5},
		{"é", "", "é", 5},
		{"\xed\xa0x\xe2\x82", "\xff", "\xed\xa0x", 8},
		{"\xff", "�", "", 4},
	}
	var input, data bytes.Buffer
	var nilOnEmpty, rawByteSet, clampCap int
	for _, c := range cases {
		row := fmt.Appendf(nil, "%d %d", len(c.cutset), c.capacity)
		for _, raw := range []string{c.cutset, c.text} {
			for i := range len(raw) {
				row = strconv.AppendUint(append(row, ' '), uint64(raw[i]), 10)
			}
		}
		fmt.Fprintln(&input, string(row))
		fmt.Fprintf(&data, "%x\t%x\t%d\t%x\n", c.cutset, c.text, c.capacity, c.prefix)

		if c.capacity >= 0 && c.prefix == "" {
			nilOnEmpty++
		}
		if c.capacity >= 0 && c.capacity != len(c.prefix) {
			clampCap++
		}
		end := len(c.text)
		for end > 0 && strings.IndexByte(c.cutset, c.text[end-1]) >= 0 {
			end--
		}
		if c.text[:end] != c.prefix {
			rawByteSet++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-right-unicode", len(cases), &input) {
		if want := observe(cases[i]); row != want {
			return fmt.Errorf("case %d: Lean=%s want=%s", i, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-trim-right-unicode-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_TRIM_RIGHT_UNICODE_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/bytes/trim-right-unicode/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanTrimRightUnicode",
		fixtures+"bridge_test.go", fixtures+"trim_right_unicode_test.go"); err != nil {
		return err
	}
	if nilOnEmpty != 4 || rawByteSet != 4 || clampCap != 9 {
		return fmt.Errorf("nil-on-empty failed %d, raw-byte cutset failed %d, clamped cap failed %d", nilOnEmpty, rawByteSet, clampCap)
	}
	fmt.Printf("trimRightUnicode passed %d fixtures; nil-on-empty failed %d, raw-byte cutset failed %d, clamped cap failed %d\n",
		len(cases), nilOnEmpty, rawByteSet, clampCap)
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
