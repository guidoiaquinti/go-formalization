package main

import (
	"fmt"
	"math/bits"
)

func check(name string, width int, count func(uint64) int) {
	wordMask := uint64(1)<<width - 1
	cases := []struct {
		word uint64
		want int
	}{
		{0, 0},
		{1, 1},
		{1 << (width - 1), 1},
		{1<<(width-1) | 1, 2},
		{wordMask, width},
		{0xaaaaaaaaaaaaaaaa & wordMask, width / 2},
	}
	var real, half, byteCount int
	for _, c := range cases {
		if count(c.word) != c.want {
			real++
		}
		if count(c.word&(uint64(1)<<(width/2)-1)) != c.want {
			half++
		}
		nonzeroBytes := 0
		for w := c.word; w != 0; w >>= 8 {
			if uint8(w) != 0 {
				nonzeroBytes++
			}
		}
		if nonzeroBytes != c.want {
			byteCount++
		}
	}
	fmt.Printf("%s: Go=%d/6 low-half=%d/6 nonzero-byte=%d/6\n", name, real, half, byteCount)
	if real != 0 || half != 4 || byteCount != 2 {
		panic("unexpected population-count control outcomes")
	}
}

func main() {
	if bits.UintSize != 64 {
		panic("OnesCount evidence requires a 64-bit native uint")
	}
	check("OnesCount16", 16, func(x uint64) int { return bits.OnesCount16(uint16(x)) })
	check("OnesCount32", 32, func(x uint64) int { return bits.OnesCount32(uint32(x)) })
	check("OnesCount64", 64, bits.OnesCount64)
	check("OnesCount", 64, func(x uint64) int { return bits.OnesCount(uint(x)) })
}
