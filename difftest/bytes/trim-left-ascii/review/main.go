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
	text, members string
	spare         int
	nilTable      bool
	want          string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []fixture{
		{"", "", -1, true, "nil"},
		{"", "", 5, true, "nil"},
		{"x", "", 2, true, "panic"},
		{"", "", 3, false, "nil"},
		{"aaaa", "a", 4, false, "nil"},
		{"\xff\x80x\xff", "\xff\x80", 5, false, "7:120,255"},
		{"éa", "\xc3", 2, false, "4:169,97"},
		{"\x00\x00a", "\x00", 1, false, "2:97"},
		{"axa", "a", 4, false, "6:120,97"},
		{"\x80a", "", 3, false, "5:128,97"},
	}
	var input, data bytes.Buffer
	var eagerNil, asciiOnly, clampCap int
	for _, c := range cases {
		count, set := len(c.members), fmt.Sprintf("%x", c.members)
		if c.nilTable {
			count, set = -1, "-"
			if c.want != "panic" {
				eagerNil++
			}
		}
		row := fmt.Appendf(nil, "%d %d", count, max(0, c.spare))
		for _, raw := range []string{c.members, c.text} {
			for i := range len(raw) {
				row = strconv.AppendUint(append(row, ' '), uint64(raw[i]), 10)
			}
		}
		fmt.Fprintln(&input, string(row))
		fmt.Fprintf(&data, "%s\t%d\t%x\t%s\n", set, c.spare, c.text, c.want)
		if c.nilTable {
			continue
		}
		start := 0
		for start < len(c.text) && c.text[start] < 128 && strings.IndexByte(c.members, c.text[start]) >= 0 {
			start++
		}
		if observe(c.text[start:], len(c.text)+c.spare-start) != c.want {
			asciiOnly++
		}
		for i := range len(c.text) {
			if strings.IndexByte(c.members, c.text[i]) < 0 {
				if observe(c.text[i:], len(c.text)-i) != c.want {
					clampCap++
				}
				break
			}
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-left-ascii", len(cases), &input) {
		if row != cases[i].want {
			return fmt.Errorf("case %d: Lean=%s want=%s", i, row, cases[i].want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-trim-left-ascii-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_TRIM_LEFT_ASCII_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/bytes/trim-left-ascii/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanTrimLeftASCII",
		fixtures+"bridge_test.go", fixtures+"trim_left_ascii_test.go"); err != nil {
		return err
	}
	if eagerNil != 2 || asciiOnly != 2 || clampCap != 5 {
		return fmt.Errorf("eager nil panic failed %d, ASCII-only failed %d, clamped cap failed %d", eagerNil, asciiOnly, clampCap)
	}
	fmt.Printf("trimLeftASCII passed %d fixtures; eager nil panic failed %d, ASCII-only failed %d, clamped cap failed %d\n",
		len(cases), eagerNil, asciiOnly, clampCap)
	return nil
}

func observe(text string, cap int) string {
	if text == "" {
		return "nil"
	}
	row := fmt.Appendf(nil, "%d:", cap)
	for i := range len(text) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(text[i]), 10)
	}
	return string(row)
}
