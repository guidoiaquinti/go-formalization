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
	text, cutset, suffix string
	spare, wantCap       int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []fixture{
		{"", "é", "", -1, 0},
		{"", "", "", 0, 0},
		{"", "年", "", 5, 0},
		{"aaaa", "a", "", 3, 0},
		{"é🚀xé", "é🚀", "xé", 4, 7},
		{"\xff�x\x80", "�", "x\x80", 2, 4},
		{"\xffx\xfe", "a", "\xffx\xfe", 5, 8},
		{"éa", "\xc3", "éa", 2, 5},
		{"é", "", "é", 3, 5},
		{"\xe2\x82x\xed\xa0", "\xff", "x\xed\xa0", 3, 6},
		{"\xff", "�", "", 3, 0},
	}
	var input, data bytes.Buffer
	var nonnilEmpty, rawByteSet, clampCap int
	for _, c := range cases {
		row := fmt.Appendf(nil, "%d %d", len(c.cutset), max(0, c.spare))
		for _, raw := range []string{c.cutset, c.text} {
			for i := range len(raw) {
				row = strconv.AppendUint(append(row, ' '), uint64(raw[i]), 10)
			}
		}
		fmt.Fprintln(&input, string(row))
		fmt.Fprintf(&data, "%x\t%x\t%d\t%x\t%d\n", c.cutset, c.text, c.spare, c.suffix, c.wantCap)

		if c.suffix == "" {
			nonnilEmpty++
		} else if c.wantCap != len(c.suffix) {
			clampCap++
		}
		start := 0
		for start < len(c.text) && strings.IndexByte(c.cutset, c.text[start]) >= 0 {
			start++
		}
		if c.text[start:] != c.suffix {
			rawByteSet++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-left-unicode", len(cases), &input) {
		if want := observe(cases[i]); row != want {
			return fmt.Errorf("case %d: Lean=%s want=%s", i, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-trim-left-unicode-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_TRIM_LEFT_UNICODE_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/bytes/trim-left-unicode/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanTrimLeftUnicode",
		fixtures+"bridge_test.go", fixtures+"trim_left_unicode_test.go"); err != nil {
		return err
	}
	if nonnilEmpty != 5 || rawByteSet != 4 || clampCap != 6 {
		return fmt.Errorf("nonnil-empty failed %d, raw-byte cutset failed %d, clamped cap failed %d", nonnilEmpty, rawByteSet, clampCap)
	}
	fmt.Printf("trimLeftUnicode passed %d fixtures; nonnil-empty failed %d, raw-byte cutset failed %d, clamped cap failed %d\n",
		len(cases), nonnilEmpty, rawByteSet, clampCap)
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
