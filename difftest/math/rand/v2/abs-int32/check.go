package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	f, err := os.CreateTemp("", "golean-rand-v2-abs")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	oracle := bufio.NewWriter(f)
	compare := func(label string, words []uint32) {
		var input bytes.Buffer
		for _, word := range words {
			fmt.Fprintln(&input, word)
		}
		for i, row := range check.LeanRows("godiff-rand-v2-abs-int32", len(words), &input) {
			word := words[i]
			x := int64(int32(word))
			if x < 0 {
				x = -x
			}
			magnitude, err := strconv.ParseUint(row, 10, 32)
			must(err)
			if magnitude != uint64(x) {
				log.Fatalf("word %#08x: Lean = %d, expected %d", word, magnitude, x)
			}
			recovered := uint32(magnitude)
			if int32(word) < 0 {
				recovered = -recovered
			}
			if recovered != word {
				log.Fatalf("word %#08x: sign recovery yielded %#08x", word, recovered)
			}
			_, err = fmt.Fprintln(oracle, word, x)
			must(err)
		}
		fmt.Printf("%s: %d Lean inputs agree with signed magnitude\n", label, len(words))
	}
	compare("boundaries", []uint32{0, 1, 2, math.MaxInt32 - 1, math.MaxInt32, 1 << 31, 1<<31 + 1, math.MaxUint32 - 1, math.MaxUint32})
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		words := make([]uint32, 20_000)
		for i := range words {
			words[i] = r.Uint32()
		}
		compare(fmt.Sprintf("seed %d", seed), words)
	}
	must(oracle.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_ABS_INT32_ROWS", f.Name()))
	must(check.TestInPackage("math/rand/v2", "TestGoLeanAbsInt32", "difftest/math/rand/v2/abs-int32/testdata/abs_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
