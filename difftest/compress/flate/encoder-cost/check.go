package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Present bool
	Words   []uint32
	Freq    []uint16
	Data    []byte
	Want    string
}

func expected(c fixture, width int) string {
	wordAt := func(index int) (uint32, string) {
		if !c.Present {
			return 0, "nil"
		}
		if index >= len(c.Words) {
			return 0, "bounds"
		}
		return c.Words[index], ""
	}
	encode := func(value uint64) string {
		if width == 32 {
			value = uint64(uint32(value))
		}
		return strconv.FormatUint(value, 10)
	}
	sum := func(failMissing bool) string {
		var total uint64
		for i, f := range c.Freq {
			if f == 0 {
				continue
			}
			word, errorText := wordAt(i)
			if errorText != "" {
				return errorText
			}
			if failMissing && word == 0 {
				return "2147483647"
			}
			total += uint64(f) * uint64(word%256)
		}
		return encode(total)
	}
	raw := func() string {
		var total uint64
		for _, symbol := range c.Data {
			word, errorText := wordAt(int(symbol))
			if errorText != "" {
				return errorText
			}
			total += uint64(max(1, word%256))
		}
		return encode(total)
	}
	return strings.Join([]string{sum(false), sum(true), raw()}, " ")
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-cost")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		for _, width := range []int{32, 64} {
			var input bytes.Buffer
			for _, c := range cases {
				present := 0
				if c.Present {
					present = 1
				}
				fmt.Fprint(&input, width, " ", present, " ", len(c.Words), " ", len(c.Freq))
				for _, w := range c.Words {
					fmt.Fprint(&input, " ", w)
				}
				for _, freq := range c.Freq {
					fmt.Fprint(&input, " ", freq)
				}
				for _, byteValue := range c.Data {
					fmt.Fprint(&input, " ", byteValue)
				}
				fmt.Fprintln(&input)
			}
			for i, row := range check.LeanRows("godiff-flate-encoder-cost", len(cases), &input) {
				c := cases[i]
				c.Want = expected(c, width)
				if row != c.Want {
					log.Fatalf("%s width %d row %d: Lean=%q want=%q", label, width, i, row, c.Want)
				}
				if width == strconv.IntSize {
					must(enc.Encode(c))
				}
			}
			fmt.Printf("%s width %d: %d Lean rows agree\n", label, width, len(cases))
		}
	}
	fixed := []fixture{{}, {Freq: []uint16{1}, Data: []byte{0}}, {Freq: []uint16{0, 0}},
		{Present: true, Freq: []uint16{0, 1}, Data: []byte{0}},
		{Present: true, Words: []uint32{0, 256, 2}, Freq: []uint16{0, 1, 2}, Data: []byte{0, 1, 2}},
		{Present: true, Words: []uint32{0}, Freq: []uint16{1, 2}},
		{Present: true, Words: []uint32{256, 1}, Freq: []uint16{1, 2}}}
	for _, size := range []int{128, 129, 130, 257, 258, 259, 512} {
		fixed = append(fixed, fixture{Present: true, Words: slices.Repeat([]uint32{0xffffffff}, size),
			Freq: slices.Repeat([]uint16{65535}, size), Data: []byte{0, 127}})
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{Present: r.Intn(7) != 0, Words: make([]uint32, r.Intn(30)),
				Freq: make([]uint16, r.Intn(35)), Data: make([]byte, r.Intn(20))}
			for j := range c.Words {
				if r.Intn(5) != 0 {
					c.Words[j] = r.Uint32()
				}
			}
			for j := range c.Freq {
				if r.Intn(3) != 0 {
					c.Freq[j] = uint16(r.Uint32())
				}
			}
			for j := range c.Data {
				c.Data[j] = byte(r.Intn(36))
			}
			cases[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_COST_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanEncoderCost", "difftest/compress/flate/encoder-cost/testdata/cost_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
