package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	r   rune
	dst []byte
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	runes := []rune{math.MinInt32, -1, 128, 0x7ff, 0x800, 0xd7ff, 0xd800, 0xdfff, 0xe000, 0xfffd, 0xffff, 0x10000, 0x10ffff, 0x110000, math.MaxInt32}
	for r := range 128 {
		runes = append(runes, rune(r))
	}
	for _, r := range runes {
		for size := range 7 {
			dst := bytes.Repeat([]byte{0xaa}, size)
			cases = append(cases, input{r, dst})
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			r := rune(rng.Uint32())
			switch i % 4 {
			case 1:
				r = rune(rng.Intn(0x110000))
			case 2:
				r = rune(rng.Intn(0x800))
			case 3:
				r = rune(rng.Intn(128))
			}
			dst := make([]byte, rng.Intn(9))
			rng.Read(dst)
			cases = append(cases, input{r, dst})
		}
	}
	var modelInput, fixtureInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		fmt.Fprint(&modelInput, c.r)
		for _, b := range c.dst {
			fmt.Fprint(&modelInput, " ", b)
		}
		fmt.Fprintln(&modelInput)
		suffix := string(c.r)
		if 0 <= c.r && c.r < 128 {
			suffix = string([]byte{0xc0 | byte(c.r>>6), 0x80 | byte(c.r&63)})
		}
		want := slices.Clone(c.dst)
		status := "bounds"
		if len(want) >= len(suffix) {
			copy(want, suffix)
			status = strconv.Itoa(len(suffix))
		}
		fmt.Fprintf(&fixtureInput, "%d x%x %s x%x\n", c.r, c.dst, status, want)
		row := []byte(status)
		for _, b := range want {
			row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
		}
		wants[i] = string(row)
	}
	for i, row := range check.LeanRows("godiff-utf8-encode-rune-non-ascii", len(cases), &modelInput) {
		if row != wants[i] {
			return fmt.Errorf("encode %#x to %x: Lean=%q want=%q", cases[i].r, cases[i].dst, row, wants[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-utf8-encode-non-ascii-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_UTF8_ENCODE_NON_ASCII_CASES", file); err != nil {
		return err
	}
	const fixtures = "difftest/unicode/utf8/encode-rune-non-ascii/testdata/"
	if err := check.TestInPackage("unicode/utf8", "TestGoLeanEncodeRuneNonASCII", fixtures+"bridge_test.go", fixtures+"encode_test.go"); err != nil {
		return err
	}
	fmt.Printf("unicode/utf8.encodeRuneNonASCII: %d boundaries and %d samples agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}
