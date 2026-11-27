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
	Length, Mode, FailAt int
	Words                []uint64
	Want                 string
}

func main() {
	file, err := os.CreateTemp("", "golean-rand-v2-global-shuffle")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	enc := json.NewEncoder(buffer)
	var fixed []fixture
	for n := -1; n <= 34; n++ {
		for mode := range 3 {
			for _, words := range [][]uint64{nil, {0}, {0, 1}, {0, 0, ^uint64(0)}, make([]uint64, 35)} {
				fixed = append(fixed, fixture{n, mode, 1, words, ""})
			}
		}
	}
	compare(fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				words := make([]uint64, rng.Intn(40))
				for j := range words {
					if j%5 != 0 {
						words[j] = rng.Uint64()
					}
				}
				cases[i] = fixture{rng.Intn(32) - 2, rng.Intn(3), rng.Intn(8), words, ""}
			}
			compare(cases, enc)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_RAND_GLOBAL_SHUFFLE_CASES", file.Name()))
	must(check.TestInPackage("math/rand/v2", "TestGoLeanGlobalShuffle",
		"difftest/math/rand/v2/global-word/testdata/bridge_test.go",
		"difftest/math/rand/v2/global-shuffle/testdata/global_shuffle_test.go"))
	fmt.Printf("%d fixed cases and 20,000 per seed (17, 2718, 65537) passed.\n", len(fixed))
}

func compare(cases []fixture, enc *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d", strconv.IntSize, c.Length, c.Mode, c.FailAt)
		for _, v := range c.Words {
			fmt.Fprintf(&input, " %d", v)
		}
		fmt.Fprintln(&input)
	}
	for i, row := range check.LeanRows("godiff-rand-v2-global-shuffle", len(cases), &input) {
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
