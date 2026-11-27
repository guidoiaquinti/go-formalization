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
	Kind, Size, Offset, Length int
	Backing                    []byte
	Want                       string
}

func main() {
	file, err := os.CreateTemp("", "golean-binary-coder-skip")
	must(err)
	defer os.Remove(file.Name())
	out := bufio.NewWriter(file)
	enc := json.NewEncoder(out)
	var fixed []fixture
	for kind := range 8 {
		for capacity := range 15 {
			backing := make([]byte, capacity)
			for i := range backing {
				backing[i] = byte(i + 10)
			}
			for _, size := range []int{0, 1, 4, 9} {
				for _, offset := range []int{math.MinInt, -1, 0, 1, capacity, capacity + 1, math.MaxInt - 1, math.MaxInt} {
					fixed = append(fixed, fixture{kind, size, offset, capacity / 2, backing, ""})
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
				cells := make([]byte, rng.Intn(65))
				rng.Read(cells)
				cases[i] = fixture{rng.Intn(8), rng.Intn(30), rng.Intn(len(cells)+6) - 3, rng.Intn(len(cells) + 1), cells, ""}
			}
			compare(cases, enc)
		}
	}
	must(out.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BINARY_CODER_SKIP_CASES", file.Name()))
	const data = "difftest/encoding/binary/coder-skip/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanCoderSkip", data+"bridge_test.go", data+"coder_skip_test.go"))
	fmt.Printf("%d edge cases and 20,000 per seed (17, 2718, 65537) agreed.\n", len(fixed))
}

func compare(cases []fixture, enc *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d", strconv.IntSize, c.Kind, c.Size, c.Offset)
		for _, b := range c.Backing {
			fmt.Fprintf(&input, " %d", b)
		}
		fmt.Fprintln(&input)
	}
	for i, row := range check.LeanRows("godiff-binary-coder-skip", len(cases), &input) {
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
