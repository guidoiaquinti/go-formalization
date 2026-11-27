package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	word        uint64
	bits, count byte
}

var primes = [...]uint64{506832829, 2654435761, 889523592379, 227718039650203, 58295818150454627, 0xcf1bbcdcb7a56463}

func oracle(word uint64, bits, count byte) uint32 {
	if count < 3 || count > 8 {
		count = 4
	}
	width := uint(64)
	if count <= 4 {
		width = 32
	}
	if bits == 0 || uint(bits) > width {
		return 0
	}
	var value, p, modulus big.Int
	value.SetUint64(word)
	p.SetUint64(primes[count-3])
	value.Mul(&value, &p)
	value.Lsh(&value, width-8*uint(count))
	modulus.SetUint64(1)
	modulus.Lsh(&modulus, width)
	value.Mod(&value, &modulus)
	value.Rsh(&value, width-uint(bits))
	return uint32(value.Uint64())
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-hash")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	native := bufio.NewWriter(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.word, c.bits, c.count)
		}
		for i, row := range check.LeanRows("godiff-flate-hash-words", len(cases), &input) {
			c := cases[i]
			want := fmt.Sprintf("%d %d", oracle(c.word, c.bits, 4), oracle(c.word, c.bits, c.count))
			if row != want {
				log.Fatalf("%s row %d %+v: Lean=%q want=%q", label, i, c, row, want)
			}
			_, err := fmt.Fprintln(native, c.word, c.bits, c.count, want)
			must(err)
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	fixed := make([]fixture, 0, 65536)
	for count := range 256 {
		for bits := range 256 {
			fixed = append(fixed, fixture{0xcafebabe01234567, byte(bits), byte(count)})
		}
	}
	for _, word := range []uint64{0, 1, 0xff, 0xffffff, 0xffffffff, 1 << 32, 1 << 63, math.MaxUint64} {
		for _, bits := range []byte{0, 1, 15, 16, 31, 32, 33, 63, 64, 65, 255} {
			for count := byte(2); count <= 9; count++ {
				fixed = append(fixed, fixture{word, bits, count})
			}
		}
	}
	compare("all width and count bytes and word boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			bits, count := byte(r.Uint32()), byte(r.Uint32())
			if i%2 == 0 {
				bits, count = byte(r.Intn(33)), byte(3+r.Intn(6))
			}
			cases[i] = fixture{r.Uint64(), bits, count}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(native.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_HASH_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanHashWords", "difftest/compress/flate/hash-words/testdata/hash_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
