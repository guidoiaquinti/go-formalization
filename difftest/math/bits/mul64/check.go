package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"math/bits"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct{ x, y uint64 }

// Sum shifted copies of x so the reference product is independent of bits.Mul64.
func product(x, y uint64) (uint64, uint64) {
	var n big.Int
	term := new(big.Int).SetUint64(x)
	for y != 0 {
		if y&1 != 0 {
			n.Add(&n, term)
		}
		term.Lsh(term, 1)
		y >>= 1
	}
	var mask, low big.Int
	mask.SetUint64(^uint64(0))
	lo := low.And(&n, &mask).Uint64()
	return n.Rsh(&n, 64).Uint64(), lo
}

func compare(label string, cases []pair) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.x, c.y)
	}
	rows := check.LeanRows("godiff-mul64", len(cases), &input)
	for i, c := range cases {
		hi, lo := bits.Mul64(c.x, c.y)
		wantHi, wantLo := product(c.x, c.y)
		if hi != wantHi || lo != wantLo || rows[i] != fmt.Sprintf("%d %d", wantHi, wantLo) {
			log.Fatalf("Mul64 %s x=%d y=%d: Go=(%d %d), Lean=%q, expected=(%d %d)",
				label, c.x, c.y, hi, lo, rows[i], wantHi, wantLo)
		}
	}
	fmt.Printf("Mul64 %s: %d products agree.\n", label, len(cases))
}

func main() {
	const limb = 1 << 32
	edges := []uint64{0, 1, 2, limb - 1, limb, limb + 1, 1<<63 - 1, 1 << 63, 1<<63 + 1, 1<<64 - limb, 1<<64 - 2, 1<<64 - 1}
	var boundaries []pair
	for _, x := range edges {
		for _, y := range edges {
			boundaries = append(boundaries, pair{x, y})
		}
	}
	compare("boundaries", boundaries)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]pair, 20_000)
		for i := range cases {
			cases[i] = pair{rng.Uint64(), rng.Uint64()}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
}
