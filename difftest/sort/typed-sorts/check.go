package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind   int
	Ints   []int
	Floats []uint64
	Words  [][]byte
	Want   string
}

var special = []uint64{
	0, 1, 1 << 63, 1<<63 | 1, 0x0010000000000000, 0x7fefffffffffffff,
	0x7ff0000000000000, 0xfff0000000000000, 0x7ff0000000000001,
	0x7ff8000000000000, 0xfff0000000000001, math.MaxUint64,
}

func sample(r *rand.Rand, kind, width, pattern int) fixture {
	c := fixture{Kind: kind}
	switch kind {
	case 0:
		c.Ints = make([]int, width)
		for j := range c.Ints {
			c.Ints[j] = int(r.Uint64())
			switch pattern % 4 {
			case 0:
				c.Ints[j] = j
			case 1:
				c.Ints[j] = width - j
			case 2:
				c.Ints[j] = r.Intn(7) - 3
			}
		}
	case 1, 2, 5:
		c.Floats = make([]uint64, width)
		for j := range c.Floats {
			switch pattern % 5 {
			case 0:
				c.Floats[j] = special[r.Intn(len(special))]
			case 1:
				c.Floats[j] = math.Float64bits(float64(j - width/2))
			case 2:
				c.Floats[j] = math.Float64bits(float64(width - j))
			case 3:
				c.Floats[j] = r.Uint64()
				if j%2 == 0 {
					c.Floats[j] |= 0x7ff0000000000000
				}
			default:
				c.Floats[j] = r.Uint64()
			}
		}
	case 3, 4:
		c.Words = make([][]byte, width)
		for j := range c.Words {
			switch pattern % 4 {
			case 0:
				c.Words[j] = []byte{0, byte(j)}
			case 1:
				c.Words[j] = []byte{255, byte(width - j)}
			case 2:
				c.Words[j] = []byte{128, byte(r.Intn(5)), 0}
			default:
				c.Words[j] = make([]byte, r.Intn(12))
				r.Read(c.Words[j])
			}
		}
	}
	return c
}

func boundaries(kind int) []fixture {
	if kind == 5 {
		var fixed []fixture
		for exponent := range 2048 {
			for _, sign := range []uint64{0, 1 << 63} {
				for _, mantissa := range []uint64{0, 1, 1 << 51, (1 << 52) - 1} {
					fixed = append(fixed, fixture{Kind: kind, Floats: []uint64{sign | uint64(exponent)<<52 | mantissa}})
				}
			}
		}
		return fixed
	}
	rng := rand.New(rand.NewSource(1))
	fixed := []fixture{{Kind: kind}}
	for _, width := range []int{0, 1, 2, 11, 12, 13, 14, 47, 48, 49, 50, 51, 63, 64, 65, 100, 101, 127, 128, 129} {
		for pattern := range 5 {
			fixed = append(fixed, sample(rng, kind, width, pattern))
		}
	}
	switch kind {
	case 0:
		fixed = append(fixed, fixture{Kind: kind, Ints: []int{math.MaxInt, math.MinInt, 0}})
	case 1, 2:
		fixed = append(fixed, fixture{Kind: kind, Floats: special})
	case 3, 4:
		fixed = append(fixed, fixture{Kind: kind, Words: [][]byte{{255}, {0, 255}, {0}, {128}, nil, {0, 0}, {239, 191, 189}}})
	}
	return fixed
}

func main() {
	file, err := os.CreateTemp("", "golean-sort-typed-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	rows := json.NewEncoder(buffer)
	for kind := range 6 {
		compare(fmt.Sprintf("kind=%d fixed", kind), boundaries(kind), rows)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			for block := range 4 {
				cases := make([]fixture, 5000)
				for i := range cases {
					width := rng.Intn(33)
					if i%113 == 0 {
						width = 33 + rng.Intn(64)
					}
					if kind == 5 {
						width = 1
					}
					cases[i] = sample(rng, kind, width, i%7)
				}
				compare(fmt.Sprintf("kind=%d seed=%d block=%d", kind, seed, block), cases, rows)
			}
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_TYPED_SORT_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanTypedSorts", "difftest/sort/typed-sorts/testdata/typed_sorts_test.go"))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d", c.Kind, strconv.IntSize)
		for _, v := range c.Ints {
			fmt.Fprintf(&wire, " %d", v)
		}
		for _, v := range c.Floats {
			fmt.Fprintf(&wire, " %d", v)
		}
		for _, s := range c.Words {
			fmt.Fprintf(&wire, " %d", len(s))
			for _, b := range s {
				fmt.Fprintf(&wire, " %d", b)
			}
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-sort-typed", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("sort typed %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
