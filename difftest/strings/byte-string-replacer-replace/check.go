package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	table *[256][]byte
	text  string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-byte-string-replacer-replace-")
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
	if err := os.Setenv("GOLEAN_STRINGS_BYTE_STRING_REPLACER_REPLACE_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)
	transit := new([256][]byte)
	transit['a'], transit['b'] = []byte{'b', 0xff}, []byte{}
	edges := []input{{nil, ""}, {nil, "\xff"}, {nil, "a🚀\x00"}, {transit, "ab\xff"},
		{new([256][]byte), ""}, {new([256][]byte), "\x00🤀\xff"}}
	for _, keys := range []int{1, 2, 3, 17, 64, 128, 256} {
		table := new([256][]byte)
		for b := range keys {
			switch b % 4 {
			case 0:
				table[b] = []byte{}
			case 1:
				table[b] = []byte{byte(b)}
			case 2:
				table[b] = []byte{byte(b + 1), byte(b), 0xff}
			case 3:
				table[b] = []byte{0, 0x80, 0xe2, 0x82}
			}
		}
		for _, width := range []int{0, 1, keys, 8*keys - 1, 8 * keys, 8*keys + 1} {
			data := make([]byte, width)
			for i := range data {
				data[i] = byte(i % max(keys, 4))
			}
			edges = append(edges, input{table, string(data)}, input{table, strings.Repeat("\xff", width)})
		}
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		for _, replacement := range []string{"", "\x00", "🚀", "\xff\x80", raw} {
			table := new([256][]byte)
			table[b] = []byte(replacement)
			for _, width := range []int{1, 3, 7, 8, 9, 30} {
				edges = append(edges, input{table, strings.Repeat(raw, width)})
			}
		}
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]input, 5_000)
			for i := range cases {
				cases[i] = sample(rng, i)
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
	const fixtures = "difftest/strings/byte-string-replacer-replace/testdata/"
	return check.TestInPackage("strings", "TestGoLeanByteStringReplacerReplace",
		fixtures+"bridge_test.go", fixtures+"byte_string_replacer_replace_test.go")
}

func sample(rng *rand.Rand, i int) input {
	var table *[256][]byte
	alphabet := [...]int{2, 7, 31, 256}[i%4]
	keys := rng.Intn(min(alphabet, 16) + 1)
	if i%31 != 0 {
		table = new([256][]byte)
		for _, b := range rng.Perm(alphabet)[:keys] {
			raw := make([]byte, rng.Intn(8))
			rng.Read(raw)
			table[b] = raw
		}
	}
	size := rng.Intn(150)
	if i%3 == 0 {
		size = max(0, keys*8+rng.Intn(3)-1)
	}
	text := make([]byte, size)
	for j := range text {
		text[j] = byte(rng.Intn(alphabet))
	}
	return input{table, string(text)}
}

func encode(c input) string {
	row := []byte("0 0")
	if c.table != nil {
		var keys []int
		for i, entry := range c.table {
			if entry != nil {
				keys = append(keys, i)
			}
		}
		row = fmt.Appendf(nil, "1 %d", len(keys))
		for _, key := range keys {
			entry := c.table[key]
			row = fmt.Appendf(row, " %d %d", key, len(entry))
			for _, b := range entry {
				row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
			}
		}
	}
	for i := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
	}
	return string(row)
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		row := encode(c)
		fmt.Fprintln(&modelInput, row)
		oracle := "panic"
		if c.table != nil {
			result := make([]byte, 0, len(c.text))
			for j := range len(c.text) {
				byte := c.text[j]
				if entry := c.table[byte]; entry != nil {
					result = append(result, entry...)
				} else {
					result = append(result, byte)
				}
			}
			oracle = hex.EncodeToString(result)
			output := []byte("bytes:")
			for j, b := range result {
				if j != 0 {
					output = append(output, ',')
				}
				output = strconv.AppendUint(output, uint64(b), 10)
			}
			wants[i] = string(output)
		} else {
			wants[i] = "panic"
		}
		if _, err := fmt.Fprintf(fixture, "%s\t%s\n", row, oracle); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-byte-string-replacer-replace", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s byteStringReplacer.Replace(%s): Lean=%s want=%s", label, encode(cases[i]), got, wants[i])
		}
	}
	fmt.Printf("strings.byteStringReplacer.Replace %s: %d Lean results agree.\n", label, len(cases))
	return nil
}
