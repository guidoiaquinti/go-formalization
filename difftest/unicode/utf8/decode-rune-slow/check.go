package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/utf8decode"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := [][]byte{nil, {}}
	for a := range 256 {
		cases = append(cases, []byte{byte(a)})
		for b := range 256 {
			cases = append(cases, []byte{byte(a), byte(b)})
		}
	}
	for _, word := range []string{"€", "🐈", "\xef\xbf\xbd", "\xe0\x9f\xbf", "\xed\xa0\x80", "\xf4\x90\x80\x80", "\xf1\x80\x80\xc0", "\xf1\x80\x80"} {
		cases = append(cases, []byte(word))
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			var data []byte
			switch i % 3 {
			case 0:
				data = make([]byte, rng.Intn(9))
				rng.Read(data)
			case 1:
				data = utf8.AppendRune(nil, rune(rng.Intn(0x110000)))
			case 2:
				data = utf8.AppendRune(nil, rune(rng.Intn(0x110000)))
				data = data[:rng.Intn(len(data))]
			}
			cases = append(cases, data)
		}
	}

	var modelInput, fixtureInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, data := range cases {
		modelInput.WriteByte('0')
		for _, b := range data {
			fmt.Fprint(&modelInput, " ", b)
		}
		modelInput.WriteByte('\n')
		r, width := utf8decode.Expected(data)
		wants[i] = fmt.Sprintf("%d %d", r, width)
		if data == nil {
			fmt.Fprintf(&fixtureInput, "nil %s\n", wants[i])
		} else {
			fmt.Fprintf(&fixtureInput, "x%x %s\n", data, wants[i])
		}
	}
	for i, row := range check.LeanRows("godiff-utf8-decode-rune-slow", len(cases), &modelInput) {
		if row != wants[i] {
			return fmt.Errorf("decode %x: Lean=%s want=%s", cases[i], row, wants[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-utf8-decode-slice-slow-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_UTF8_DECODE_SLICE_SLOW_CASES", file); err != nil {
		return err
	}
	const fixtures = "difftest/unicode/utf8/decode-rune-slow/testdata/"
	if err := check.TestInPackage("unicode/utf8", "TestGoLeanDecodeRuneSlow", fixtures+"bridge_test.go", fixtures+"decode_test.go"); err != nil {
		return err
	}
	fmt.Printf("unicode/utf8.decodeRuneSlow: %d boundaries and %d samples agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}
