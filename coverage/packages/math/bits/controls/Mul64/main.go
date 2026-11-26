// Exercise the Mul64 reconstruction theorem and a lost carry.
package main

import (
	"fmt"
	"log"
	"math/big"
	"math/bits"
)

// Split into 16-bit digits so each partial product fits in a uint32.
func product(x, y uint64) *big.Int {
	var result, term big.Int
	for i := uint(0); i < 64; i += 16 {
		for j := uint(0); j < 64; j += 16 {
			a, b := uint16(x>>i), uint16(y>>j)
			term.SetUint64(uint64(uint32(a) * uint32(b)))
			result.Add(&result, term.Lsh(&term, i+j))
		}
	}
	return &result
}

func withoutFinalCarry(x, y uint64) (uint64, uint64) {
	x0, x1 := uint64(uint32(x)), x>>32
	y0, y1 := uint64(uint32(y)), y>>32
	t := x1*y0 + ((x0 * y0) >> 32)
	// Go adds ((t & mask32) + x0*y1) >> 32 to the high word here.
	return x1*y1 + (t >> 32), x * y
}

func main() {
	var words []uint64
	limbs := []uint64{0, 1, 1 << 31, 1<<32 - 1}
	for _, hi := range limbs {
		for _, lo := range limbs {
			words = append(words, hi<<32|lo)
		}
	}
	for _, control := range []struct {
		name         string
		mul          func(uint64, uint64) (uint64, uint64)
		wantFailures int
	}{
		{"bits.Mul64", bits.Mul64, 0},
		{"omit final carry", withoutFinalCarry, 97},
	} {
		failures := 0
		for _, x := range words {
			for _, y := range words {
				hi, lo := control.mul(x, y)
				got := new(big.Int).SetUint64(hi)
				got.Lsh(got, 64).Add(got, new(big.Int).SetUint64(lo))
				want := product(x, y)
				if got.Cmp(want) != 0 {
					if failures == 0 {
						fmt.Printf("%s first mismatch x=%#x y=%#x got=%#x want=%#x\n", control.name, x, y, got, want)
					}
					failures++
				}
			}
		}
		fmt.Printf("%s: %d/%d reconstruction failures\n", control.name, failures, len(words)*len(words))
		if failures != control.wantFailures {
			log.Fatal("unexpected result for ", control.name)
		}
	}
}
