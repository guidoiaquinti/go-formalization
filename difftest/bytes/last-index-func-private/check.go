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
)

type input struct {
	s     string
	kind  int
	truth bool
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	texts := []string{"", "abc", "é🚀🚀", "abc🚀", "�x", "a\xff�", "a\xe2\x82", "\xf4\x90\x80\x80", "\x00"}
	for b := range 256 {
		texts = append(texts, string([]byte{byte(b)}), string([]byte{'a', byte(b), 0x80, 'z'}))
	}
	var edges []input
	for _, s := range texts {
		for kind := range 6 {
			edges = append(edges, input{s, kind, false}, input{s, kind, true})
		}
	}
	cases := append([]input(nil), edges...)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			cases = append(cases, sample(rng, 0))
		}
	}
	var leanInput, fixture bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		fmt.Fprintln(&leanInput, encodeInput(c))
		wants[i] = expected(c)
		fmt.Fprintf(&fixture, "%t %d x%x %s\n", c.truth, c.kind, c.s, wants[i])
	}
	for i, got := range check.LeanRows("godiff-bytes-last-index-func-private", len(cases), &leanInput) {
		if got != wants[i] {
			return fmt.Errorf("%+v: Lean=%q want=%q", cases[i], got, wants[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-last-index-func-private-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_LAST_INDEX_FUNC_CASES", path); err != nil {
		return err
	}
	const data = "difftest/bytes/last-index-func-private/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanLastIndexFunc", data+"bridge_test.go", data+"last_index_func_test.go"); err != nil {
		return err
	}
	fmt.Printf("bytes.lastIndexFunc: %d boundaries and %d samples agree, seeds %v.\n", len(edges), len(cases)-len(edges), seeds)
	return nil
}

func sample(rng *rand.Rand, _ int) input {
	chunks := []string{"é", "🚀", "�", "\xf0\x9f", "\xff", "\x80", "abc"}
	var s strings.Builder
	for range rng.Intn(49) {
		if rng.Intn(2) == 0 {
			s.WriteByte(byte(rng.Intn(256)))
		} else {
			s.WriteString(chunks[rng.Intn(len(chunks))])
		}
	}
	return input{s.String(), rng.Intn(6), rng.Intn(2) != 0}
}

func predicate(kind int) func(rune) bool {
	switch kind {
	case 0:
		return nil
	case 1:
		return func(rune) bool { return false }
	case 2:
		return func(rune) bool { return true }
	case 3:
		return func(r rune) bool { return r == '\uFFFD' }
	case 4:
		return func(r rune) bool { return r == '🚀' }
	default:
		return func(r rune) bool { return r%7 == 2 }
	}
}

func encodeInput(c input) string {
	truth := 0
	if c.truth {
		truth = 1
	}
	row := fmt.Appendf(nil, "%d %d", truth, c.kind)
	for i := range len(c.s) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.s[i]), 10)
	}
	return string(row)
}

func describe(index int, calls []rune) string {
	row := strconv.AppendInt(nil, int64(index), 10)
	row = append(row, '|')
	for i, r := range calls {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendInt(row, int64(r), 10)
	}
	return string(row)
}

func expected(c input) string {
	checkRune := predicate(c.kind)
	if checkRune == nil && c.s != "" {
		return "panic"
	}
	// Reverse the offsets obtained from forward string range.
	type decoded struct {
		offset int
		rune   rune
	}
	var runes []decoded
	for offset, r := range c.s {
		runes = append(runes, decoded{offset, r})
	}
	var calls []rune
	for _, entry := range slices.Backward(runes) {
		calls = append(calls, entry.rune)
		if checkRune(entry.rune) == c.truth {
			return describe(entry.offset, calls)
		}
	}
	return describe(-1, calls)
}
