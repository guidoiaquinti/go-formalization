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
	input, members string
	valid          bool
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ascii := make([]byte, 128)
	for i := range ascii {
		ascii[i] = byte(i)
	}
	cases := []fixture{
		{"", "", true},
		{"\x00\x01\x7f\x00", "\x00\x01\x7f", true},
		{"\x80a", "", false},
		{"a\xffb", "a", false},
		{"\x00zéx", "\x00z", false},
		{"a\x80a\x00", "a", false},
		{"\x7f\x7f", "\x7f", true},
		{"\x7f\x80\x00", "\x7f", false},
		{string(ascii), string(ascii), true},
		{string(ascii) + "\xffz", string(ascii), false},
		{"\xff\x80", "", false},
		{"\x00\x00", "\x00", true},
	}
	var input, data bytes.Buffer
	var cleared, continued, insertedHigh int
	want := make([]string, len(cases))
	for i, c := range cases {
		row := []byte{'0'}
		for j := range len(c.input) {
			row = strconv.AppendUint(append(row, ' '), uint64(c.input[j]), 10)
		}
		fmt.Fprintln(&input, string(row))
		want[i] = observe(c.members, c.valid)
		fmt.Fprintf(&data, "x%x %s\n", c.input, want[i])
		if c.valid {
			continue
		}
		if observe("", false) != want[i] {
			cleared++
		}
		asciiOnly := make([]byte, 0, len(c.input))
		for j := range len(c.input) {
			if c.input[j] < 128 {
				asciiOnly = append(asciiOnly, c.input[j])
			}
		}
		if observe(string(asciiOnly), false) != want[i] {
			continued++
		}
		for j := range len(c.input) {
			if c.input[j] >= 128 {
				if observe(c.members+string([]byte{c.input[j]}), false) != want[i] {
					insertedHigh++
				}
				break
			}
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-make-ascii-set", len(cases), &input) {
		if row != want[i] {
			return fmt.Errorf("case %d: Lean=%s want=%s", i, row, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-bytes-make-ascii-set-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_MAKE_ASCII_SET_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/bytes/make-ascii-set/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanMakeASCIISet",
		fixtures+"bridge_test.go", fixtures+"make_ascii_set_test.go"); err != nil {
		return err
	}
	if cleared != 5 || continued != 5 || insertedHigh != 7 {
		return fmt.Errorf("cleared on failure failed %d, continued after failure failed %d, inserted rejected byte failed %d", cleared, continued, insertedHigh)
	}
	fmt.Printf("makeASCIISet passed %d fixtures; cleared on failure failed %d, continued after failure failed %d, inserted rejected byte failed %d\n",
		len(cases), cleared, continued, insertedHigh)
	return nil
}

func observe(members string, valid bool) string {
	result := append(strconv.AppendBool(nil, valid), '|')
	for i := range 256 {
		digit := byte('0')
		if strings.IndexByte(members, byte(i)) >= 0 {
			digit = '1'
		}
		result = append(result, digit)
	}
	return string(result)
}
