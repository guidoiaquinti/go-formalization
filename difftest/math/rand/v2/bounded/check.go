package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Bound uint64
	Words []uint64
	Want  string
}

func main() {
	file, err := os.CreateTemp("", "golean-rand-v2-bounded")
	must(err)
	defer os.Remove(file.Name())
	out := bufio.NewWriter(file)
	enc := json.NewEncoder(out)
	var fixed []fixture
	for bit := range 64 {
		value := uint64(1) << bit
		for _, bound := range []uint64{value - 1, value, value + 1, ^value} {
			for _, words := range [][]uint64{nil, {0}, {1}, {0, 0, 0, ^uint64(0)}, {^uint64(0)}, {value - 1, value, ^uint64(0)}} {
				fixed = append(fixed, fixture{bound, words, ""})
			}
		}
	}
	compare(fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				bound := rng.Uint64() >> uint(i%64)
				words := make([]uint64, rng.Intn(8))
				for j := range words {
					words[j] = rng.Uint64()
				}
				if len(words) > 0 && i%3 == 0 {
					words[0] = 0
				}
				rows[i] = fixture{bound, words, ""}
			}
			compare(rows, enc)
		}
	}
	must(out.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_RAND_BOUNDED_CASES", file.Name()))
	const dir = "difftest/math/rand/v2/bounded/testdata/"
	must(check.TestInPackage("math/rand/v2", "TestGoLeanRandBounded", dir+"bridge_test.go", dir+"bounded_test.go"))
	fmt.Printf("%d boundaries and 20,000 inputs per seed (17, 2718, 65537) passed.\n", len(fixed))
}

func compare(cases []fixture, enc *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d", strconv.IntSize, c.Bound)
		for _, w := range c.Words {
			fmt.Fprintf(&input, " %d", w)
		}
		fmt.Fprintln(&input)
	}
	for i, result := range check.LeanRows("godiff-rand-v2-bounded", len(cases), &input) {
		c := cases[i]
		c.Want = result
		must(enc.Encode(c))
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
