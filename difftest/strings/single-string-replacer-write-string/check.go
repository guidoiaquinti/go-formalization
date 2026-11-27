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

type reply struct{ limit, failure int }

type input struct {
	pattern, value, text string
	stringWriter         bool
	script               []reply
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-single-replacer-write-string-")
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
	if err := os.Setenv("GOLEAN_STRINGS_SINGLE_REPLACER_WRITE_STRING_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)
	patterns := []string{"a", "aa", "aba", "\xff\xff", "\xff", "\x00", "\xe2\x82", "🚀"}
	values := []string{"", "a", "aba", "\x00", "\xff\x00", "\xe2\x82", "替"}
	texts := []string{"", "ababa", "aaaaa", "é🚀é🚀", "\xff\xff\xff", "\x00a\x00", "\xffa\xe2\x82", "\x7f\xff", "\xff"}
	scripts := [][]reply{nil, {{0, 0}}, {{-1, 1}}, {{1, 2}}, {{-1, 0}, {1, 1}},
		{{0, 0}, {0, 0}, {-1, 2}}, {{0, 0}, {1, 0}, {1, 0}, {0, 1}},
		{{-1, 0}, {-1, 0}, {-1, 0}, {-1, 0}, {1, 2}}}
	var edges []input
	for _, pattern := range patterns {
		for _, value := range values {
			for _, text := range texts {
				for _, script := range scripts {
					for _, stringWriter := range []bool{false, true} {
						edges = append(edges, input{pattern, value, text, stringWriter, script},
							input{pattern, value, text + pattern + pattern + "tail", stringWriter, script})
					}
				}
			}
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
	const fixtures = "difftest/strings/single-string-replacer-write-string/testdata/"
	return check.TestInPackage("strings", "TestGoLeanSingleStringReplacerWriteString",
		fixtures+"bridge_test.go", fixtures+"single_string_replacer_write_string_test.go")
}

func raw(rng *rand.Rand, count, alphabet int) string {
	value := make([]byte, count)
	for i := range value {
		value[i] = byte(rng.Intn(alphabet))
	}
	return string(value)
}

func sample(rng *rand.Rand, i int) input {
	alphabet := [...]int{2, 7, 256}
	size := alphabet[i%len(alphabet)]
	pattern := raw(rng, 1+rng.Intn(15), size)
	value := raw(rng, rng.Intn(20), 256)
	text := raw(rng, rng.Intn(36), size)
	switch i % 4 {
	case 0:
		text = pattern + text + pattern
	case 1:
		position := rng.Intn(len(text) + 1)
		text = text[:position] + pattern + pattern + text[position:]
	case 2:
		text += pattern
	}
	script := make([]reply, rng.Intn(12))
	for j := range script {
		script[j] = reply{rng.Intn(8) - 1, 0}
		if i%5 != 0 && rng.Intn(4) == 0 {
			script[j].failure = 1 + rng.Intn(2)
		}
	}
	return input{pattern, value, text, i%2 == 0, script}
}

func expected(c input) string {
	var segments []string
	start := 0
	for i := 0; i+len(c.pattern) <= len(c.text); {
		if c.text[i:i+len(c.pattern)] == c.pattern {
			segments = append(segments, c.text[start:i], c.value)
			i += len(c.pattern)
			start = i
		} else {
			i++
		}
	}
	segments = append(segments, c.text[start:])
	mode := byte('B')
	if c.stringWriter {
		mode = 'S'
	}
	var calls []string
	n, failure := 0, 0
	for i, segment := range segments {
		event := []byte{mode}
		for j := range len(segment) {
			event = strconv.AppendUint(append(event, ','), uint64(segment[j]), 10)
		}
		calls = append(calls, string(event))
		accepted := len(segment)
		if i < len(c.script) {
			if c.script[i].limit >= 0 {
				accepted = min(accepted, c.script[i].limit)
			}
			failure = c.script[i].failure
		}
		n += accepted
		if failure != 0 {
			break
		}
	}
	return fmt.Sprintf("%d:%d:%s", n, failure, strings.Join(calls, ";"))
}

func encode(c input) string {
	mode := 0
	if c.stringWriter {
		mode = 1
	}
	row := fmt.Appendf(nil, "%d %d %d %d", mode, len(c.script), len(c.pattern), len(c.value))
	for _, reply := range c.script {
		row = fmt.Appendf(row, " %d %d", reply.limit, reply.failure)
	}
	for _, raw := range []string{c.pattern, c.value, c.text} {
		for j := range len(raw) {
			row = strconv.AppendUint(append(row, ' '), uint64(raw[j]), 10)
		}
	}
	return string(row)
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		row := encode(c)
		fmt.Fprintln(&modelInput, row)
		wants[i] = expected(c)
		if _, err := fmt.Fprintf(fixture, "%s\t%s\n", row, wants[i]); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-single-string-replacer-write-string", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s WriteString input %s: Lean=%s want=%s", label, encode(cases[i]), got, wants[i])
		}
	}
	fmt.Printf("strings.singleStringReplacer.WriteString %s: %d Lean traces agree.\n", label, len(cases))
	return nil
}
