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

type input struct{ pattern, value, text string }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-single-replacer-replace-")
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
	if err := os.Setenv("GOLEAN_STRINGS_SINGLE_REPLACER_REPLACE_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)
	patterns := []string{"a", "aa", "aba", "\xff\xff", "\x00", "\xff", "\xe2\x82", "🚀", "🚀é🚀"}
	values := []string{"", "a", "aba", "\x00", "\xff\x00", "\xe2\x82", "替"}
	texts := []string{"", "ababa", "aaaaa", "é🚀é🚀", "\xff\xff\xff", "\x00a\x00",
		"\xffa\xe2\x82", strings.Repeat("aaab", 64)}
	var edges []input
	for _, pattern := range patterns {
		for _, value := range values {
			for _, text := range texts {
				edges = append(edges, input{pattern, value, text}, input{pattern, value, text + pattern + pattern})
			}
		}
	}
	// Exhaust nonempty binary patterns through length four and texts through length seven.
	for pattern := 1; pattern < 31; pattern++ {
		for text := range 255 {
			edges = append(edges, input{binaryBytes(pattern), "\xff\x00", binaryBytes(text)})
		}
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
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
	const fixtures = "difftest/strings/single-string-replacer-replace/testdata/"
	return check.TestInPackage("strings", "TestGoLeanSingleStringReplacerReplace",
		fixtures+"bridge_test.go", fixtures+"single_string_replacer_replace_test.go")
}

func binaryBytes(index int) string {
	var raw []byte
	for bits := index + 1; bits > 1; bits >>= 1 {
		raw = append(raw, byte(bits&1)*255)
	}
	return string(raw)
}

func raw(rng *rand.Rand, count, alphabet int) string {
	value := make([]byte, count)
	for i := range value {
		value[i] = byte(rng.Intn(alphabet))
	}
	return string(value)
}

func sample(rng *rand.Rand, index int) input {
	alphabet := [...]int{2, 7, 256}
	size := alphabet[index%len(alphabet)]
	pattern := raw(rng, 1+rng.Intn(21), size)
	value := raw(rng, rng.Intn(25), 256)
	text := raw(rng, rng.Intn(60), size)
	switch index % 5 {
	case 0:
		text = pattern + text + pattern
	case 1:
		position := rng.Intn(len(text) + 1)
		text = text[:position] + pattern + pattern + text[position:]
	case 2:
		text = pattern
	case 3:
		value = pattern + value
	}
	return input{pattern, value, text}
}

func expected(c input) string {
	var out strings.Builder
	for i := 0; i < len(c.text); {
		if n := len(c.pattern); i+n <= len(c.text) && c.text[i:i+n] == c.pattern {
			out.WriteString(c.value)
			i += n
		} else {
			out.WriteByte(c.text[i])
			i++
		}
	}
	return out.String()
}

func describe(raw string) string {
	row := []byte("bytes:")
	for i := range len(raw) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(raw[i]), 10)
	}
	return string(row)
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		row := strconv.AppendInt(nil, int64(len(c.pattern)), 10)
		row = strconv.AppendInt(append(row, ' '), int64(len(c.value)), 10)
		for _, raw := range []string{c.pattern, c.value, c.text} {
			for j := range len(raw) {
				row = strconv.AppendUint(append(row, ' '), uint64(raw[j]), 10)
			}
		}
		fmt.Fprintln(&modelInput, string(row))
		want := expected(c)
		wants[i] = describe(want)
		if _, err := fmt.Fprintf(fixture, "%x\t%x\t%x\t%x\n", c.pattern, c.value, c.text, want); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-single-string-replacer-replace", len(cases), &modelInput) {
		if got != wants[i] {
			c := cases[i]
			return fmt.Errorf("%s Replace(%x,%x,%x): Lean=%s want=%s",
				label, c.pattern, c.value, c.text, got, wants[i])
		}
	}
	fmt.Printf("strings.singleStringReplacer.Replace %s: %d Lean inputs agree.\n", label, len(cases))
	return nil
}
