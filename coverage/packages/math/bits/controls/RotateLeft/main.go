package main

import (
	"fmt"
	"log"
	"math"
	"math/bits"
)

type rotation func(uint64, int) uint64

func main() {
	for _, f := range []struct {
		name  string
		width int
		run   rotation
	}{
		{"RotateLeft8", 8, func(x uint64, k int) uint64 { return uint64(bits.RotateLeft8(uint8(x), k)) }},
		{"RotateLeft16", 16, func(x uint64, k int) uint64 { return uint64(bits.RotateLeft16(uint16(x), k)) }},
		{"RotateLeft32", 32, func(x uint64, k int) uint64 { return uint64(bits.RotateLeft32(uint32(x), k)) }},
		{"RotateLeft64", 64, bits.RotateLeft64},
		{"RotateLeft", bits.UintSize, func(x uint64, k int) uint64 { return uint64(bits.RotateLeft(uint(x), k)) }},
	} {
		check(f.name, f.width, f.run)
	}
}

func check(name string, width int, native rotation) {
	high := uint64(1) << (width - 1)
	mask := high | (high - 1)
	tests := []struct {
		x     uint64
		count int
		want  uint64
	}{
		{high, 1, 1}, {1, -1, high}, {1, 0, 1}, {high | 1, width, high | 1},
		{high | 1, -width, high | 1}, {1, math.MinInt, 1}, {high, math.MinInt + 1, 1},
		{1, math.MaxInt, high}, {mask, 3, mask},
	}
	for _, c := range []struct {
		label string
		run   rotation
		fail  int
	}{
		{"Go", native, 0},
		{"unsigned shift without wrap", func(x uint64, k int) uint64 { return (x << uint(k&(width-1))) & mask }, 3},
		{"rotate in the opposite direction", func(x uint64, k int) uint64 { return native(x, -k) }, 4},
		{"use magnitude of negative count", func(x uint64, k int) uint64 {
			if k < 0 {
				k = -k
			}
			return native(x, k)
		}, 2},
	} {
		mismatches := 0
		for _, t := range tests {
			if got := c.run(t.x, t.count); got != t.want {
				mismatches++
				if c.fail == 0 {
					log.Fatalf("%s(%#x,%d)=%#x, want %#x", name, t.x, t.count, got, t.want)
				}
			}
		}
		if mismatches != c.fail {
			log.Fatalf("%s %s failed %d fixtures, want %d", name, c.label, mismatches, c.fail)
		}
		fmt.Printf("%s %s disagreed on %d of %d fixtures\n", name, c.label, mismatches, len(tests))
	}
}
