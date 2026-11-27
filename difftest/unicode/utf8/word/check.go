package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func expected(width int, data []byte) string {
	if len(data) < width {
		return "bounds"
	}
	if width == 4 {
		return strconv.FormatUint(uint64(binary.LittleEndian.Uint32(data)), 10)
	}
	return strconv.FormatUint(binary.LittleEndian.Uint64(data), 10)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases [][]byte
	for n := range 13 {
		for _, b := range []byte{0, 1, 0x7f, 0x80, 0xff} {
			cases = append(cases, bytes.Repeat([]byte{b}, n))
		}
	}
	for lane := range 8 {
		for b := range 256 {
			data := make([]byte, 9)
			data[lane] = byte(b)
			data[8] = 0xff
			cases = append(cases, data)
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			data := make([]byte, rng.Intn(17))
			rng.Read(data)
			cases = append(cases, data)
		}
	}
	for _, width := range []int{4, 8} {
		var input bytes.Buffer
		for _, data := range cases {
			fmt.Fprint(&input, width)
			for _, b := range data {
				fmt.Fprint(&input, " ", b)
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-utf8-word", len(cases), &input) {
			if want := expected(width, cases[i]); row != want {
				return fmt.Errorf("width %d data %x: Lean=%s want=%s", width, cases[i], row, want)
			}
		}
	}

	dir, err := os.MkdirTemp("", "golean-utf8-word-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	var inputs bytes.Buffer
	for _, data := range cases {
		fmt.Fprintf(&inputs, "x%x %s\n", data, expected(strconv.IntSize/8, data))
	}
	if err := os.WriteFile(file, inputs.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_UTF8_WORD_CASES", file); err != nil {
		return err
	}
	const fixtures = "difftest/unicode/utf8/word/testdata/"
	if err := check.TestInPackage("unicode/utf8", "TestGoLeanWord", fixtures+"bridge_test.go", fixtures+"word_test.go"); err != nil {
		return err
	}
	fmt.Printf("unicode/utf8.word: %d boundaries and %d samples, seeds %v. Go word width %d, Lean widths 4 and 8.\n",
		boundaries, len(cases)-boundaries, seeds, strconv.IntSize/8)
	return nil
}
