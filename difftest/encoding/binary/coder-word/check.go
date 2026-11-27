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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Width, Order, Length, Offset int
	Word                         uint64
	Backing                      []byte
	Want                         string
}

func main() {
	file, err := os.CreateTemp("", "golean-binary-coder-word")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	enc := json.NewEncoder(buffer)
	for _, width := range []int{2, 4, 8} {
		var fixed []fixture
		for capacity := range 19 {
			backing := make([]byte, capacity)
			for i := range backing {
				backing[i] = byte(239 + i)
			}
			for _, length := range []int{0, capacity / 2, capacity} {
				for _, offset := range []int{math.MinInt, -1, 0, 1, capacity - width, capacity - width + 1, math.MaxInt} {
					for order := range 5 {
						fixed = append(fixed, fixture{width, order, length, offset, 0x92c4031875d6bafe, backing, ""})
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
					capacity := rng.Intn(32)
					backing := make([]byte, capacity)
					rng.Read(backing)
					cases[i] = fixture{width, rng.Intn(5), rng.Intn(capacity + 1), rng.Intn(capacity+5) - 2, rng.Uint64(), backing, ""}
				}
				compare(cases, enc)
			}
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BINARY_CODER_WORD_CASES", file.Name()))
	const testdata = "difftest/encoding/binary/coder-word/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanCoderWord", testdata+"bridge_test.go", testdata+"coder_word_test.go"))
}

func compare(cases []fixture, out *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d", c.Width, c.Order, c.Length, c.Offset, c.Word)
		for _, b := range c.Backing {
			fmt.Fprintf(&wire, " %d", b)
		}
		fmt.Fprintln(&wire)
	}
	for i, r := range check.LeanRows("godiff-binary-coder-word", len(cases), &wire) {
		c := cases[i]
		c.Want = r
		must(out.Encode(c))
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
