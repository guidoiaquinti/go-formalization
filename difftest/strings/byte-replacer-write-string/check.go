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
	table        *[256]byte
	text         string
	stringWriter bool
	script       []reply
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
	if err := os.Setenv("GOLEAN_STRINGS_BYTE_REPLACER_WRITE_STRING_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)
	identity := new([256]byte)
	for i := range identity {
		identity[i] = byte(i)
	}
	cycle, constant, mixed := new([256]byte), new([256]byte), new([256]byte)
	*mixed = *identity
	mixed['b'], mixed[0xff], mixed[0] = 0xff, 0, 'b'
	for i := range cycle {
		cycle[i] = byte(i + 1)
		constant[i] = 0xff
	}
	texts := []string{"", "a", "b", "ababbb", "aaaaa", "bba", "é🚀é🚀",
		"b\xff\xff\xff", "\x00a\x00", "\xffa\xe2\x82", string(identity[:])}
	scripts := [][]reply{nil, {{0, 0}}, {{-1, 1}}, {{1, 2}}, {{-1, 0}, {1, 1}},
		{{0, 0}, {0, 0}, {-1, 2}}, {{-1, 0}, {1, 0}, {1, 0}, {0, 1}},
		{{-1, 0}, {-1, 0}, {-1, 0}, {-1, 0}, {1, 2}}}
	var edges []input
	for _, table := range []*[256]byte{nil, identity, cycle, constant, mixed} {
		for _, text := range texts {
			for _, script := range scripts {
				for _, stringWriter := range []bool{false, true} {
					edges = append(edges, input{table, text, stringWriter, script})
				}
			}
		}
	}
	for shift := range 256 {
		table := new([256]byte)
		for b := range table {
			table[b] = byte(b + shift)
		}
		edges = append(edges, input{table, string(identity[:]), true, nil})
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]input, 5_000)
			for i := range cases {
				cases[i] = sample(rng, i, identity)
			}
			if err := compare(fixture, fmt.Sprintf("seed=%d block=%d", seed, block), cases); err != nil {
				return err
			}
		}
	}
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	const fixtures = "difftest/strings/byte-replacer-write-string/testdata/"
	return check.TestInPackage("strings", "TestGoLeanByteReplacerWriteString",
		fixtures+"bridge_test.go", fixtures+"byte_replacer_write_string_test.go")
}

func raw(rng *rand.Rand, count, alphabet int) string {
	value := make([]byte, count)
	for i := range value {
		value[i] = byte(rng.Intn(alphabet))
	}
	return string(value)
}

func sample(rng *rand.Rand, i int, identity *[256]byte) input {
	alphabet := [...]int{2, 7, 256}
	size := alphabet[i%len(alphabet)]
	text := raw(rng, rng.Intn(36), size)
	var table *[256]byte
	if i%23 != 0 {
		table = new([256]byte)
		*table = *identity
		for range [...]int{0, 1, 16, 256}[i%4] {
			table[rng.Intn(size)] = byte(rng.Uint32())
		}
	}
	script := make([]reply, rng.Intn(12))
	for j := range script {
		script[j] = reply{rng.Intn(8) - 1, 0}
		if i%5 != 0 && rng.Intn(4) == 0 {
			script[j].failure = 1 + rng.Intn(2)
		}
	}
	return input{table, text, i%2 == 0, script}
}

type event struct {
	mode byte
	raw  string
}

func expected(c input) string {
	if c.text == "" {
		return "0:0:"
	}
	if c.table == nil {
		return "panic"
	}
	var events []event
	unchanged := byte('B')
	if c.stringWriter {
		unchanged = 'S'
	}
	start := 0
	for i := range len(c.text) {
		original := c.text[i]
		if replacement := c.table[original]; replacement != original {
			if start < i {
				events = append(events, event{unchanged, c.text[start:i]})
			}
			events = append(events, event{'B', string([]byte{replacement})})
			start = i + 1
		}
	}
	if start != len(c.text) {
		events = append(events, event{unchanged, c.text[start:]})
	}
	var calls []string
	n, failure := 0, 0
	for i, event := range events {
		row := []byte{event.mode}
		for j := range len(event.raw) {
			row = strconv.AppendUint(append(row, ','), uint64(event.raw[j]), 10)
		}
		calls = append(calls, string(row))
		accepted := len(event.raw)
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
	mode, hasTable := 0, 0
	if c.stringWriter {
		mode = 1
	}
	if c.table != nil {
		hasTable = 1
	}
	row := fmt.Appendf(nil, "%d %d %d", mode, hasTable, len(c.script))
	for _, reply := range c.script {
		row = fmt.Appendf(row, " %d %d", reply.limit, reply.failure)
	}
	if c.table != nil {
		for _, value := range c.table {
			row = strconv.AppendUint(append(row, ' '), uint64(value), 10)
		}
	}
	for j := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[j]), 10)
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
	for i, got := range check.LeanRows("godiff-strings-byte-replacer-write-string", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s WriteString input %s: Lean=%s want=%s", label, encode(cases[i]), got, wants[i])
		}
	}
	fmt.Printf("strings.byteReplacer.WriteString %s: %d Lean traces agree.\n", label, len(cases))
	return nil
}
