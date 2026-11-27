package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/utf16oracle"
)

type input struct {
	Words []uint16
	Dst   []rune
	Spare int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	bounds := [][]uint16{nil, {}, {0}, {0xffff}, {0xd800}, {0xdc00},
		{0xd800, 0xdc00}, {0xdbff, 0xdfff}, {0xd800, 0xd800, 0xdc00}, {0xdc00, 0xd800}}
	prefixes := [][]rune{nil, {}, {-1}, {0, 0xd800, 0x10ffff}}
	var cases []input
	for _, words := range bounds {
		for _, dst := range prefixes {
			for spare := range 3 {
				cases = append(cases, input{words, dst, spare})
			}
		}
	}
	edgeCount := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			cases = append(cases, sample(rng, i))
		}
	}

	dir, err := os.MkdirTemp("", "golean-utf16-decode-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	var stdin bytes.Buffer
	for _, c := range cases {
		length, count := len(c.Dst), len(c.Words)
		if c.Dst == nil {
			length = -1
		}
		if c.Words == nil {
			count = -1
		}
		fmt.Fprint(&stdin, c.Spare, " ", length, " ", count)
		for _, r := range c.Dst {
			fmt.Fprint(&stdin, " ", r)
		}
		for _, w := range c.Words {
			fmt.Fprint(&stdin, " ", w)
		}
		fmt.Fprintln(&stdin)
	}
	if err := os.WriteFile(filepath.Join(dir, "inputs"), stdin.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_UTF16_DECODE_DIR", dir); err != nil {
		return err
	}
	rows := check.LeanRows("godiff-utf16-decode-append", len(cases), &stdin)
	const fixture = "difftest/unicode/utf16/decode-append/testdata/decode_append_test.go"
	if err := check.TestInPackage("unicode/utf16", "TestGoLeanDecodeAppend", fixture); err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(dir, "outputs"))
	if err != nil {
		return err
	}
	actual := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(actual) != len(cases) {
		return fmt.Errorf("Go produced %d rows, want %d", len(actual), len(cases))
	}
	for i, c := range cases {
		want := append(slices.Clone(c.Dst), utf16oracle.Decode(c.Words)...)
		if expected := format(want); actual[i] != expected || rows[i] != expected {
			return fmt.Errorf("case %d %+v: Go=%v Lean=%q want=%v", i, c, actual[i], rows[i], want)
		}
	}
	fmt.Printf("unicode/utf16.decode: %d boundary and %d sampled cases agree, seeds %v.\n",
		edgeCount, len(cases)-edgeCount, seeds)
	return nil
}

func sample(rng *rand.Rand, i int) input {
	words := make([]uint16, rng.Intn(17))
	for j := range words {
		words[j] = uint16(rng.Uint32())
		switch i % 3 {
		case 1:
			words[j] = 0xd800 + uint16(rng.Intn(2048))
		case 2:
			words[j] = 0xd800 + uint16(1024*(j%2)+rng.Intn(1024))
		}
	}
	dst := make([]rune, rng.Intn(9))
	for j := range dst {
		dst[j] = rune(rng.Uint32())
	}
	return input{words, dst, rng.Intn(4)}
}

func format(runes []rune) string {
	text := []byte("runes")
	for _, r := range runes {
		text = strconv.AppendInt(append(text, ' '), int64(r), 10)
	}
	return string(text)
}
