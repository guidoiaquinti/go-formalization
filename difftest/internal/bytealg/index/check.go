package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"slices"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Input, Needle                   []byte
	InputOffset, NeedleOffset, Want int
}

func main() {
	maxWidth := 31
	switch runtime.GOARCH {
	case "arm64":
		maxWidth = 32
	case "amd64":
	default:
		fmt.Printf("Architecture %s is outside the tested small-needle profiles.\n", runtime.GOARCH)
		return
	}
	file, err := os.CreateTemp("", "golean-bytealg-index-")
	must(err)
	defer os.Remove(file.Name())
	writer := bufio.NewWriter(file)
	fixtureRows := json.NewEncoder(writer)
	// Fill flanks with the needle's repeated body to distinguish early,
	// final and out-of-range alignments, and vary the compared addresses.
	var edges []fixture
	for width := 2; width <= maxWidth; width++ {
		needle := bytes.Repeat([]byte{0x80}, width)
		needle[width-1] = 0xff
		for offset := range 64 {
			prefix := bytes.Repeat([]byte{0x80}, offset+1)
			early := append(bytes.Clone(needle), prefix...)
			late := append(bytes.Clone(prefix), needle...)
			missing := append(bytes.Clone(prefix), bytes.Repeat([]byte{0x80}, width)...)
			repeated := append(bytes.Clone(late), needle...)
			for _, data := range [][]byte{early, late, missing, repeated} {
				edges = append(edges, fixture{data, needle, offset, 63 - offset, 0})
			}
		}
	}
	compare("boundaries", edges, fixtureRows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5_000)
			for i := range cases {
				needle := make([]byte, 2+rng.Intn(maxWidth-1))
				input := make([]byte, len(needle)+rng.Intn(129))
				rng.Read(needle)
				rng.Read(input)
				switch i % 5 {
				case 0:
					copy(input[:len(needle)], needle)
				case 1:
					copy(input[len(input)-len(needle):], needle)
				case 2:
					copy(input[rng.Intn(len(input)-len(needle)+1):], needle)
				case 3:
					input = append(bytes.Clone(needle), input...)
				}
				cases[i] = fixture{input, needle, rng.Intn(64), rng.Intn(64), 0}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, fixtureRows)
		}
	}
	must(writer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BYTEALG_INDEX_CASES", file.Name()))
	const data = "difftest/internal/bytealg/index/testdata/"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanBytealgIndex", data+"bridge_test.go", data+"index_test.go"))
}

func compare(label string, cases []fixture, fixtureRows *json.Encoder) {
	var input bytes.Buffer
	for i := range cases {
		c := &cases[i]
		fmt.Fprint(&input, len(c.Needle))
		for _, word := range [][]byte{c.Needle, c.Input} {
			for _, b := range word {
				fmt.Fprintf(&input, " %d", b)
			}
		}
		fmt.Fprintln(&input)
		c.Want = -1
		for j := 0; j+len(c.Needle) <= len(c.Input); j++ {
			if slices.Equal(c.Input[j:j+len(c.Needle)], c.Needle) {
				c.Want = j
				break
			}
		}
	}
	for i, got := range check.LeanRows("godiff-internal-bytealg-index", len(cases), &input) {
		c := cases[i]
		want := fmt.Sprintf("%d,%d", c.Want, c.Want)
		if got != want {
			log.Fatalf("%s row %d input=%x needle=%x: Lean=%s want=%s", label, i, c.Input, c.Needle, got, want)
		}
		must(fixtureRows.Encode(c))
	}
	fmt.Printf("internal/bytealg Index %s: %d Lean rows agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
