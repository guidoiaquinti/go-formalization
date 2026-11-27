package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Left, Right                   []byte
	Size, LeftOffset, RightOffset int
	Alias, Want                   bool
}

var widths = []int{3, 5, 7, 9, 15, 17, 31, 33, 63, 65, 127, 129, 257}

func main() {
	f, err := os.CreateTemp("", "golean-memequal-varlen-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	var edges []fixture
	for _, size := range widths {
		for offset := range 64 {
			left := bytes.Repeat([]byte{0x80}, size+1)
			inside, outside := bytes.Clone(left), bytes.Clone(left)
			inside[size-1] = 0xff
			outside[size] = 0xff
			edges = append(edges, fixture{Left: left, Right: inside, Size: size, LeftOffset: offset, RightOffset: 63 - offset},
				fixture{Left: left, Right: outside, Size: size, LeftOffset: offset, RightOffset: 63 - offset},
				fixture{Left: left, Right: left, Size: size, LeftOffset: offset, RightOffset: offset, Alias: true})
		}
	}
	compare("boundaries", edges, rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				size := widths[rng.Intn(len(widths))]
				left := make([]byte, size+4)
				rng.Read(left)
				right := bytes.Clone(left)
				switch i % 5 {
				case 0:
					right[rng.Intn(size)] ^= 255
				case 1:
					right[size+rng.Intn(4)] ^= 255
				case 2:
					rng.Read(right)
				}
				cases[i] = fixture{left, right, size, rng.Intn(64), rng.Intn(64), i%5 == 4, false}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_BYTEALG_MEMEQUAL_VARLEN_CASES", f.Name()))
	const file = "difftest/internal/bytealg/memequal-varlen/testdata/memequal_varlen_test.go"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanMemequalVarlen", file))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var input bytes.Buffer
	for i := range cases {
		c := &cases[i]
		fmt.Fprintf(&input, "%d %d", c.Size, len(c.Left))
		for _, raw := range [][]byte{c.Left, c.Right} {
			for _, b := range raw {
				fmt.Fprintf(&input, " %d", b)
			}
		}
		fmt.Fprintln(&input)
		c.Want = slices.Equal(c.Left[:c.Size], c.Right[:c.Size])
	}
	for i, got := range check.LeanRows("godiff-internal-bytealg-memequal-varlen", len(cases), &input) {
		c := cases[i]
		if got != strconv.FormatBool(c.Want) {
			log.Fatalf("%s row %d width=%d: Lean=%s want=%t", label, i, c.Size, got, c.Want)
		}
		must(rows.Encode(c))
	}
	fmt.Printf("memequal_varlen %s: %d Lean rows agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
