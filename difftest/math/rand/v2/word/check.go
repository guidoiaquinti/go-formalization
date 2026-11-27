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
	Mode int
	Word uint64
	Want string
}

func main() {
	file, err := os.CreateTemp("", "golean-rand-v2-word")
	must(err)
	defer os.Remove(file.Name())
	output := bufio.NewWriter(file)
	enc := json.NewEncoder(output)

	var fixed []fixture
	for mode := range 4 {
		for bit := range 64 {
			single := uint64(1) << bit
			for _, word := range []uint64{single, single - 1, single + 1, ^single, 0, ^uint64(0)} {
				fixed = append(fixed, fixture{mode, word, ""})
			}
		}
	}
	compare(fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				mode := 0
				if i%8 == 0 {
					mode = rng.Intn(4)
				}
				rows[i] = fixture{mode, rng.Uint64(), ""}
			}
			compare(rows, enc)
		}
	}
	must(output.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_RAND_WORD_CASES", file.Name()))
	const data = "difftest/math/rand/v2/word/testdata/"
	must(check.TestInPackage("math/rand/v2", "TestGoLeanRandWord", data+"bridge_test.go", data+"word_test.go"))
	fmt.Printf("%d fixed cases and 20,000 per seed (17, 2718, 65537) passed.\n", len(fixed))
}

func compare(cases []fixture, out *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d\n", strconv.IntSize, c.Mode, c.Word)
	}
	for i, row := range check.LeanRows("godiff-rand-v2-word", len(cases), &input) {
		c := cases[i]
		c.Want = row
		must(out.Encode(c))
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
