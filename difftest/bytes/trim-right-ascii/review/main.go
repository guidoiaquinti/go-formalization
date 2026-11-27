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
		{"", "", 0, true, "0:"},
		{"", "", 5, true, "5:"},
		{"x", "", 2, true, "panic"},
		{"aaaa", "a", 4, false, "8:"},
		{"\xffx\x80\xff", "\xff\x80", 5, false, "9:255,120"},
		{"aé", "\xa9", 2, false, "5:97,195"},
		{"\x00a\x00\x00", "\x00", 1, false, "5:0,97"},
		{"axa", "a", 4, false, "7:97,120"},
		{"a\x80", "", 3, false, "5:97,128"},
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
		capacity := len(c.text) + c.spare
		row := fmt.Appendf(nil, "%d %d", count, capacity)
		for _, raw := range []string{c.members, c.text} {
			for i := range len(raw) {
				row = strconv.AppendUint(append(row, ' '), uint64(raw[i]), 10)
			}
		}
		fmt.Fprintln(&input, string(row))
		fmt.Fprintf(&data, "%s\t%d\t%x\t%s\n", set, c.spare, c.text, c.want)
		if c.want == "panic" || c.want == "nil" {
			continue
		}
		end := len(c.text)
		for end > 0 && c.text[end-1] < 128 && strings.IndexByte(c.members, c.text[end-1]) >= 0 {
			end--
		}
		if observe(c.text[:end], capacity) != c.want {
			asciiOnly++
		}
		end = len(c.text)
		for end > 0 && strings.IndexByte(c.members, c.text[end-1]) >= 0 {
			end--
		}
		if observe(c.text[:end], end) != c.want {
			clampCap++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-right-ascii", len(cases), &input) {
		if row != cases[i].want {
			return fmt.Errorf("case %d: Lean=%s want=%s", i, row, cases[i].want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-trim-right-ascii-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_TRIM_RIGHT_ASCII_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/bytes/trim-right-ascii/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanTrimRightASCII",
		fixtures+"bridge_test.go", fixtures+"trim_right_ascii_test.go"); err != nil {
		return err
	}
	if eagerNil != 3 || asciiOnly != 2 || clampCap != 7 {
		return fmt.Errorf("eager nil panic failed %d, ASCII-only failed %d, clamped cap failed %d", eagerNil, asciiOnly, clampCap)
	}
	fmt.Printf("trimRightASCII passed %d fixtures; eager nil panic failed %d, ASCII-only failed %d, clamped cap failed %d\n",
		len(cases), eagerNil, asciiOnly, clampCap)
	return nil
}

func observe(text string, cap int) string {
	row := fmt.Appendf(nil, "%d:", cap)
	for i := range len(text) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(text[i]), 10)
	}
	return string(row)
}
