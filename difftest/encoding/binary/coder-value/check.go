package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, Shape, Count, Order, Length, Offset int
	Words                                     []uint64
	Backing                                   []byte
	Want                                      string
}

func generate(r *rand.Rand, kind, shape, count, order, length, capacity, offset int) fixture {
	width := [...]int{1, 1, 1, 2, 2, 4, 4, 8, 8, 4, 8, 4, 8, 0}[kind]
	parts := 1
	if kind == 11 || kind == 12 {
		parts = 2
	}
	n := count
	if shape == 0 {
		n = 1
	} else if shape == 3 {
		n++
	}
	if width == 0 {
		n = 0
	}
	words := make([]uint64, n*parts)
	for i := range words {
		words[i] = r.Uint64() & (^uint64(0) >> (64 - 8*width))
		if kind == 0 {
			words[i] &= 1
		}
		if kind >= 9 {
			// Keep the reflected conversions in the finite subset.
			words[i] &^= uint64(1) << (8*width - 2)
		}
	}
	backing := make([]byte, capacity)
	r.Read(backing)
	if kind >= 9 && kind <= 12 {
		for i := range backing {
			backing[i] &^= 0x40
		}
	}
	return fixture{kind, shape, count, order, length, offset, words, backing, ""}
}

func main() {
	file, err := os.CreateTemp("", "golean-binary-coder-value")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	enc := json.NewEncoder(buffer)
	rng := rand.New(rand.NewSource(17))
	var fixed []fixture
	for kind := range 14 {
		for shape := range 4 {
			for order := range 2 {
				for count := range 4 {
					for _, length := range []int{0, 1, 3, 8, 12, 19, 33} {
						fixed = append(fixed, generate(rng, kind, shape, count, order, length, 33, (count+1)%4))
					}
				}
			}
		}
	}
	compare(fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				capacity := rng.Intn(90)
				cases[i] = generate(rng, rng.Intn(14), rng.Intn(4), rng.Intn(7), rng.Intn(2), rng.Intn(capacity+1), capacity, rng.Intn(capacity+4))
			}
			compare(cases, enc)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BINARY_CODER_VALUE_CASES", file.Name()))
	const data = "difftest/encoding/binary/coder-value/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanCoderValue", data+"bridge_test.go", data+"coder_value_test.go"))
	fmt.Printf("%d boundary cases and 20,000 per seed (17, 2718, 65537) passed.\n", len(fixed))
}

func compare(cases []fixture, enc *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d %d", c.Kind, c.Shape, c.Count, c.Order, c.Offset, c.Length)
		for _, v := range c.Words {
			fmt.Fprintf(&wire, " %d", v)
		}
		for _, b := range c.Backing {
			fmt.Fprintf(&wire, " %d", b)
		}
		fmt.Fprintln(&wire)
	}
	for i, row := range check.LeanRows("godiff-binary-coder-value", len(cases), &wire) {
		c := cases[i]
		c.Want = row
		must(enc.Encode(c))
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
