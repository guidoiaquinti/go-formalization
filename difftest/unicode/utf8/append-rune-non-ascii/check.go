package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
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
		for _, dst := range [][]byte{nil, {0xff, 0, 'x'}} {
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
			dst := make([]byte, rng.Intn(17))
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
		want := []byte(string(c.dst) + suffix)
		fmt.Fprintf(&fixtureInput, "%d x%x %x\n", c.r, c.dst, want)
		var row []byte
		for j, b := range want {
			if j != 0 {
				row = append(row, ' ')
			}
			row = strconv.AppendUint(row, uint64(b), 10)
		}
		wants[i] = string(row)
	}
	for i, row := range check.LeanRows("godiff-utf8-append-rune-non-ascii", len(cases), &modelInput) {
		if row != wants[i] {
			return fmt.Errorf("append %#x to %x: Lean=%q want=%q", cases[i].r, cases[i].dst, row, wants[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-utf8-append-non-ascii-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_UTF8_APPEND_NON_ASCII_CASES", file); err != nil {
		return err
	}
	const fixtures = "difftest/unicode/utf8/append-rune-non-ascii/testdata/"
	if err := check.TestInPackage("unicode/utf8", "TestGoLeanAppendRuneNonASCII", fixtures+"bridge_test.go", fixtures+"append_test.go"); err != nil {
		return err
	}
	fmt.Printf("unicode/utf8.appendRuneNonASCII: %d boundaries and %d samples agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}
