package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/bits"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	f, err := os.CreateTemp("", "golean-strconv-integer-bits")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	native := bufio.NewWriter(f)
	compare := func(label string, width int, words []uint64) {
		var input bytes.Buffer
		for _, word := range words {
			fmt.Fprintln(&input, width, word)
		}
		for i, got := range check.LeanRows("godiff-strconv-integer-bits", len(words), &input) {
			word := words[i]
			want := strconv.FormatBool(bits.OnesCount64(word) <= 1)
			if width == 8 {
				// Bit 5 is high exactly in the alternating second blocks of 32.
				low := word
				if (word/32)%2 == 0 {
					low += 32
				}
				want = strconv.FormatUint(low, 10)
			}
			if got != want {
				log.Fatalf("%s width %d word %#x: Lean = %s, want %s", label, width, word, got, want)
			}
			if width == 8 || width == bits.UintSize {
				_, err := fmt.Fprintln(native, width, word, want)
				must(err)
			}
		}
		fmt.Printf("%s width %d: %d Lean rows agree\n", label, width, len(words))
	}
	all := make([]uint64, 256)
	for i := range all {
		all[i] = uint64(i)
	}
	compare("exhaustive", 8, all)
	for _, width := range []int{32, 64} {
		mask := ^uint64(0) >> (64 - width)
		edges := []uint64{0, mask}
		for bit := range width {
			x := uint64(1) << bit
			edges = append(edges, x, (x-1)&mask, (x+1)&mask, x|(uint64(1)<<((bit+1)%width)))
		}
		compare("boundaries", width, edges)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			words := make([]uint64, 20_000)
			for i := range words {
				x := rng.Uint64()
				if i%4 == 0 {
					x = uint64(1) << rng.Intn(width)
				}
				words[i] = x & mask
			}
			compare(fmt.Sprintf("seed %d", seed), width, words)
		}
	}
	must(native.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_STRCONV_BITS_ROWS", f.Name()))
	must(check.TestInPackage("internal/strconv", "TestGoLeanIntegerBits", "difftest/internal/strconv/integer-bits/testdata/export_test.go", "difftest/internal/strconv/integer-bits/testdata/integer_bits_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
