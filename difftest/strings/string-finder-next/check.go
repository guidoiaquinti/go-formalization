package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ pattern, text string }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-string-finder-next-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := os.Setenv("GOLEAN_STRINGS_STRING_FINDER_NEXT_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)

	patterns := []string{"", "aa", "aba", "mississi", "abcxxxabc", "\x00a", "\x80",
		"\xff\x00", "é", "🚀", "\xff🚀", "\xe2\x82", strings.Repeat("aba", 40)}
	texts := []string{"", "aa", "aba", "abababa", "mississipimississippi", "abcxxxabc",
		"\x00a", "\x80", "\xff\x00", "é", "🚀", "\xff🚀", "\xe2\x82", strings.Repeat("aba", 41)}
	var edges []input
	for _, pattern := range patterns {
		for _, text := range texts {
			edges = append(edges, input{pattern, text}, input{pattern, text + pattern + pattern})
		}
	}
	// Exhaust binary patterns through length four and binary texts through length seven.
	for pattern := range 31 {
		for text := range 255 {
			edges = append(edges, input{binaryBytes(pattern), binaryBytes(text)})
		}
	}
	for b := range 256 {
		byteString := string([]byte{byte(b)})
		edges = append(edges, input{byteString, "\xff\x00" + byteString},
			input{byteString + "\xff", "\xff\x00" + byteString})
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			alphabet := [...]int{2, 7, 256}
			size := alphabet[i%len(alphabet)]
			pattern, text := raw(rng, 40, size), raw(rng, 96, size)
			switch i % 5 {
			case 0:
				text += pattern
			case 1:
				pos := rng.Intn(len(text) + 1)
				text = text[:pos] + pattern + text[pos:] + pattern
			case 2:
				text = pattern + text
			case 3:
				text = pattern
			}
			cases[i] = input{pattern, text}
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
	const fixtures = "difftest/strings/string-finder-next/testdata/"
	return check.TestInPackage("strings", "TestGoLeanStringFinderNext",
		fixtures+"bridge_test.go", fixtures+"string_finder_next_test.go")
}

// Encode each binary string uniquely using a leading one as its length marker.
func binaryBytes(index int) string {
	var value []byte
	for bits := index + 1; bits > 1; bits >>= 1 {
		value = append(value, byte(bits&1)*255)
	}
	return string(value)
}

func raw(rng *rand.Rand, limit, alphabet int) string {
	buf := make([]byte, rng.Intn(limit))
	for i := range buf {
		buf[i] = byte(rng.Intn(alphabet))
	}
	return string(buf)
}

func expected(in input) int {
	n := len(in.pattern)
	for i := 0; i+n <= len(in.text); i++ {
		if in.text[i:i+n] == in.pattern {
			return i
		}
	}
	return -1
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, in := range cases {
		row := strconv.AppendInt(nil, int64(len(in.pattern)), 10)
		for _, raw := range []string{in.pattern, in.text} {
			for j := range len(raw) {
				row = strconv.AppendUint(append(row, ' '), uint64(raw[j]), 10)
			}
		}
		fmt.Fprintln(&modelInput, string(row))
		wants[i] = strconv.Itoa(expected(in))
		if _, err := fmt.Fprintf(fixture, "%x\t%x\t%s\n", in.pattern, in.text, wants[i]); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-string-finder-next", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s stringFinder.next(%x,%x): Lean=%s want=%s", label,
				cases[i].pattern, cases[i].text, got, wants[i])
		}
	}
	fmt.Printf("strings.stringFinder.next %s: %d Lean inputs agree.\n", label, len(cases))
	return nil
}
