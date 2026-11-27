package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	word  uint16
	count byte
}

func expected(c fixture) uint16 {
	width := int(c.count) % 16
	if width == 0 {
		width = 16
	}
	var result uint16
	for i := range width {
		result = result*2 + (c.word>>i)%2
	}
	return result
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-reverse")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	native := bufio.NewWriter(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.word, c.count)
		}
		for i, got := range check.LeanRows("godiff-flate-reverse-bits", len(cases), &input) {
			c := cases[i]
			want := expected(c)
			if got != strconv.Itoa(int(want)) {
				log.Fatalf("%s row %d %+v: Lean=%q want=%d", label, i, c, got, want)
			}
			_, err := fmt.Fprintln(native, c.word, c.count, want)
			must(err)
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	var edges []fixture
	for n := range 256 {
		for _, word := range []uint16{0, 0xffff, 0xaaaa, 0x5555, 0xff00, 0x00ff} {
			edges = append(edges, fixture{word, byte(n)})
		}
		for bit := range 16 {
			word := uint16(1) << bit
			edges = append(edges, fixture{word, byte(n)}, fixture{^word, byte(n)})
		}
	}
	compare("all counts and bit-position boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{uint16(r.Uint32()), byte(r.Uint32())}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(native.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_REVERSE_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanReverseBits", "difftest/compress/flate/reverse-bits/testdata/reverse_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
