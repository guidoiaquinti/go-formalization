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
	Kind, Form, Order, Length int
	Words                     []uint64
	Backing                   []byte
	Want                      string
}

func makeFixture(kind, form, order, n, length int, random *rand.Rand) fixture {
	width := [...]int{1, 1, 1, 2, 2, 4, 4, 8, 8, 4, 8, 1}[kind]
	if form < 2 {
		n = 1
	}
	words := make([]uint64, n)
	for i := range words {
		words[i] = random.Uint64() & (^uint64(0) >> (64 - 8*width))
		if kind == 0 {
			words[i] &= 1
		}
	}
	backing := make([]byte, length+3)
	random.Read(backing)
	return fixture{kind, form, order, length, words, backing, ""}
}

func main() {
	file, err := os.CreateTemp("", "golean-binary-fast-codec")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	enc := json.NewEncoder(buffer)
	rng := rand.New(rand.NewSource(17))
	var fixed []fixture
	for kind := range 12 {
		for form := range 4 {
			for order := range 2 {
				for n := range 5 {
					for length := range 18 {
						fixed = append(fixed, makeFixture(kind, form, order, n, length, rng))
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
				cases[i] = makeFixture(rng.Intn(12), rng.Intn(4), rng.Intn(2), rng.Intn(12), rng.Intn(48), rng)
			}
			compare(cases, enc)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BINARY_FAST_CODEC_CASES", file.Name()))
	const data = "difftest/encoding/binary/fast-codec/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanFastCodec", data+"bridge_test.go", data+"fast_codec_test.go"))
	fmt.Printf("%d fixed cases and 20,000 per seed (17, 2718, 65537) agreed.\n", len(fixed))
}

func compare(cases []fixture, enc *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d %d", c.Kind, c.Form, c.Order, len(c.Words), c.Length)
		for _, v := range c.Words {
			fmt.Fprintf(&input, " %d", v)
		}
		for _, v := range c.Backing {
			fmt.Fprintf(&input, " %d", v)
		}
		fmt.Fprintln(&input)
	}
	for i, row := range check.LeanRows("godiff-binary-fast-codec", len(cases), &input) {
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
