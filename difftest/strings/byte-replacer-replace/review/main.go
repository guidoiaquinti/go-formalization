package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Table       *[256]byte
	Input, Want []byte
	Panic       bool
}

func table(changes map[byte]byte) *[256]byte {
	t := new([256]byte)
	for i := range t {
		t[i] = byte(i)
	}
	for from, to := range changes {
		t[from] = to
	}
	return t
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []fixture{
		{table(nil), []byte("a\x00\xff"), []byte("a\x00\xff"), false},
		{table(map[byte]byte{'a': 'b', 'b': 'c'}), []byte("ab"), []byte("bc"), false},
		{table(map[byte]byte{'x': 'y', 'y': 'x'}), []byte("xyx"), []byte("yxy"), false},
		{table(map[byte]byte{0xc3: 'M', 0xff: 0}), []byte("\xc3\xa9\xff"), []byte("M\xa9\x00"), false},
		{table(map[byte]byte{0: 'a', 'a': 0xff}), []byte("\x00a\xff"), []byte("a\xff\xff"), false},
		{nil, nil, nil, false},
		{nil, []byte{0}, nil, true},
		{table(map[byte]byte{0: 1}), nil, nil, false},
	}
	var rows bytes.Buffer
	cascade, ascii := 0, 0
	for _, c := range cases {
		if c.Table == nil {
			fmt.Fprint(&rows, 0)
		} else {
			fmt.Fprint(&rows, 1)
			for _, b := range c.Table {
				fmt.Fprint(&rows, " ", b)
			}
			twice, asciiOnly := make([]byte, len(c.Input)), bytes.Clone(c.Input)
			for i, b := range c.Input {
				twice[i] = c.Table[c.Table[b]]
				if b < 128 {
					asciiOnly[i] = c.Table[b]
				}
			}
			if !bytes.Equal(twice, c.Want) {
				cascade++
			}
			if !bytes.Equal(asciiOnly, c.Want) {
				ascii++
			}
		}
		for _, b := range c.Input {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
	}
	for i, result := range check.LeanRows("godiff-strings-byte-replacer-replace", len(cases), &rows) {
		want := "panic"
		if !cases[i].Panic {
			fields := make([]string, len(cases[i].Want))
			for j, b := range cases[i].Want {
				fields[j] = strconv.Itoa(int(b))
			}
			want = "bytes:" + strings.Join(fields, ",")
		}
		if result != want {
			return fmt.Errorf("fixture %d: Lean %q, want %q", i, result, want)
		}
	}
	if cascade != 3 || ascii != 1 {
		return fmt.Errorf("cascade and ASCII-only failures=%d/%d, want 3/1", cascade, ascii)
	}
	dir, err := os.MkdirTemp("", "golean-byte-replacer-control-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	fixturePath := filepath.Join(dir, "cases.json")
	data, err := json.Marshal(cases)
	if err != nil {
		return err
	}
	if err := os.WriteFile(fixturePath, data, 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_REVIEW_BYTE_REPLACER", fixturePath); err != nil {
		return err
	}
	const base = "difftest/strings/byte-replacer-replace/"
	if err := check.TestInPackage("strings", "TestGoLeanReviewByteReplacer",
		base+"testdata/bridge_test.go", base+"review/testdata/review_test.go"); err != nil {
		return err
	}
	fmt.Printf("byteReplacer.Replace passed %d Go and Lean fixtures; cascade and ASCII-only mutants failed %d and %d.\n",
		len(cases), cascade, ascii)
	return nil
}
