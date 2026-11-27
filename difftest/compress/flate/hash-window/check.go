package main

import (
	"bytes"
	"encoding/binary"
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
	Size, DstSize int
	Source        []byte
	Dst           []uint32
	Want          string
}

func expected(c fixture) string {
	hash := func(window []byte) uint64 {
		word := uint64(binary.LittleEndian.Uint32(window))
		return (word * 2654435761) % (1 << 32) / (1 << 15)
	}
	single := "bounds"
	if len(c.Source) >= 4 {
		single = strconv.FormatUint(hash(c.Source), 10)
	}
	count := max(0, c.Size-3)
	cells := slices.Clone(c.Dst)
	for i := range min(count, c.DstSize) {
		cells[i] = uint32(hash(c.Source[i:]))
	}
	wire := make([]string, len(cells))
	for i, x := range cells {
		wire[i] = strconv.FormatUint(uint64(x), 10)
	}
	return fmt.Sprintf("%s|%t|%s", single, count > c.DstSize, strings.Join(wire, ","))
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-hash-window")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, c.Size, " ", len(c.Source), " ", c.DstSize)
			for _, b := range c.Source {
				fmt.Fprint(&input, " ", b)
			}
			for _, word := range c.Dst {
				fmt.Fprint(&input, " ", word)
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-flate-hash-window", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c)
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	for srcCap := range 13 {
		for srcLen := 0; srcLen <= srcCap; srcLen++ {
			for dstCap := range 11 {
				for dstLen := 0; dstLen <= dstCap; dstLen++ {
					c := fixture{Size: srcLen, DstSize: dstLen, Source: make([]byte, srcCap), Dst: make([]uint32, dstCap)}
					for i := range c.Source {
						c.Source[i] = byte((i + 1) * 33)
					}
					for i := range c.Dst {
						c.Dst[i] = uint32(i) + 0xffff0000
					}
					fixed = append(fixed, c)
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{Source: make([]byte, r.Intn(65)), Dst: make([]uint32, r.Intn(65))}
			c.Size = r.Intn(len(c.Source) + 1)
			c.DstSize = r.Intn(len(c.Dst) + 1)
			r.Read(c.Source)
			for j := range c.Dst {
				c.Dst[j] = r.Uint32()
			}
			cases[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_HASH_WINDOW_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanHashWindow", "difftest/compress/flate/hash-window/testdata/window_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
