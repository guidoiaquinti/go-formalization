package main

import (
	"fmt"
	"math/bits"
)

func check(width int, reverseBytes, reverseBits func(uint64) uint64) {
	var actual, identity, reflected int
	for i := range width {
		input := uint64(1) << i
		want := uint64(1) << ((width/8-1-i/8)*8 + i%8)
		if reverseBytes(input) != want {
			actual++
		}
		if input != want {
			identity++
		}
		if reverseBits(input) != want {
			reflected++
		}
	}
	fmt.Printf("ReverseBytes%d single-bit inputs: Go=%d/%d identity=%d/%d bit-reversal=%d/%d\n",
		width, actual, width, identity, width, reflected, width)
	if actual != 0 || identity != width || reflected != width {
		panic("unexpected byte-reversal control counts")
	}
}

func main() {
	check(32,
		func(x uint64) uint64 { return uint64(bits.ReverseBytes32(uint32(x))) },
		func(x uint64) uint64 { return uint64(bits.Reverse32(uint32(x))) })
	check(64, bits.ReverseBytes64, bits.Reverse64)
}
