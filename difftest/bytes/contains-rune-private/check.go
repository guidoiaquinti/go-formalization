package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text string
	seek rune
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	seeks := []rune{math.MinInt32, -1, 0, 'a', 0x7f, 0x80, 0xa9, 0xe9, 0x7ff,
		0x800, 0xd7ff, 0xd800, 0xdfff, 0xe000, utf8.RuneError, 0xffff,
		0x10000, '🚀', utf8.MaxRune, utf8.MaxRune + 1, math.MaxInt32}
	texts := []string{"", "a\x00\x7f", "\u0080©é", "\u07ff\u0800\ud7ff\ue000", "\uffff\U00010000\U0010ffff",
		"é🚀", "�", "\xff", "\xe2\x82", "\xed\xa0\x80", "\xc0\xaf", "\xf4\x90\x80\x80",
		strings.Repeat("©", 128), strings.Repeat("©", 128) + "é"}
	var edges []input
	for _, text := range texts {
		for _, seek := range seeks {
			edges = append(edges, input{text, seek})
		}
	}
	for i := range 256 {
		raw := string([]byte{byte(i)})
		for _, seek := range seeks {
			edges = append(edges, input{raw, seek}, input{raw + "é\xff🚀" + raw, seek})
		}
		edges = append(edges, input{raw, rune(i)})
	}
	dir, err := os.MkdirTemp("", "golean-bytes-contains-rune-private-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BYTES_CONTAINS_RUNE_PRIVATE_CASES", path); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fixture := bufio.NewWriter(file)
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			cases[i] = sample(r, i)
		}
		if err := compare(fixture, fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	const data = "difftest/bytes/contains-rune-private/testdata/"
	return check.TestInPackage("bytes", "TestGoLeanContainsRunePrivate", data+"bridge_test.go",
		data+"contains_rune_private_test.go")
}

func encode(c input) string {
	row := strconv.AppendInt(nil, int64(c.seek), 10)
	for i := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
	}
	return string(row)
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		fmt.Fprintln(&modelInput, encode(c))
		wants[i] = strconv.FormatBool(expected(c))
		if _, err := fmt.Fprintf(fixture, "%x\t%d\t%s\n", c.text, c.seek, wants[i]); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-bytes-contains-rune-private", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s containsRune(%x, %d): Lean=%s want=%s", label, cases[i].text,
				cases[i].seek, got, wants[i])
		}
	}
	fmt.Printf("bytes.containsRune %s: %d Lean rows agree.\n", label, len(cases))
	return nil
}

func expected(c input) bool {
	raw := []byte(c.text)
	for len(raw) != 0 {
		r, width := utf8.DecodeRune(raw)
		if r == c.seek {
			return true
		}
		raw = raw[width:]
	}
	return false
}

func sample(r *rand.Rand, i int) input {
	tokens := []string{"a", "\x00", "é", "©", "🚀", "�", "\xff", "\xe2\x82", "\xed\xa0\x80"}
	var text strings.Builder
	for range r.Intn(48) {
		if r.Intn(2) == 0 {
			text.WriteByte(byte(r.Intn(256)))
		} else {
			text.WriteString(tokens[r.Intn(len(tokens))])
		}
	}
	seek := rune(r.Uint32())
	switch i % 5 {
	case 0:
		seek = rune(r.Intn(128))
	case 1:
		seek = utf8.RuneError
	case 2:
		seek = rune(r.Intn(utf8.MaxRune + 1))
		text.WriteRune(seek)
	case 3:
		seek = rune(r.Intn(utf8.MaxRune + 1))
	}
	return input{text.String(), seek}
}
