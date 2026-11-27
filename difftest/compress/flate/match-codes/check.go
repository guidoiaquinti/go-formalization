package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	offset uint32
	length byte
}

func distanceCode(distance uint32) uint32 {
	base, width := uint32(0), uint32(1)
	for symbol := uint32(0); ; symbol++ {
		if distance < base+width {
			return symbol
		}
		base += width
		if symbol >= 3 && symbol%2 == 1 {
			width *= 2
		}
	}
}

func lengthIndex(length byte) int {
	base, width := uint32(0), uint32(1)
	for symbol := 0; symbol < 28; symbol++ {
		if length < 255 && uint32(length) < base+width {
			return symbol
		}
		base += width
		if symbol >= 7 && symbol%4 == 3 {
			width *= 2
		}
	}
	return 28
}

func main() {
	file, err := os.CreateTemp("", "golean-flate-match-codes")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	native := bufio.NewWriter(file)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.offset, c.length)
		}
		for i, row := range check.LeanRows("godiff-flate-match-codes", len(cases), &input) {
			c := cases[i]
			var distance uint32
			if c.offset < 32768 {
				distance = distanceCode(c.offset)
			} else {
				residue := c.offset % 32768
				distance = distanceCode(residue)
				if residue < 256 {
					distance = 14 + residue/128
				}
			}
			want := fmt.Sprintf("%d %d", distance, lengthIndex(c.length))
			if row != want {
				log.Fatalf("%s offset %d length %d: Lean=%q want=%q", label, c.offset, c.length, row, want)
			}
			_, err := fmt.Fprintln(native, c.offset, c.length, want)
			must(err)
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	edges := make([]fixture, 65536)
	for i := range edges {
		edges[i] = fixture{uint32(i), byte(i)}
	}
	compare("exhaustive low 16-bit offsets and all length bytes", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{r.Uint32(), byte(r.Uint32())}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(native.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_FLATE_MATCH_ROWS", file.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanMatchCodes", "difftest/compress/flate/match-codes/testdata/match_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
