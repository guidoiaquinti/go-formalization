package main

import (
	"fmt"
	"log"
	"math/bits"
)

type op func(uint64, uint64, uint64) (uint64, uint64)
type fixture struct{ x, y, in, word, out uint64 }

func main() {
	add := func(x, y, c uint64) (uint64, uint64) {
		a, b := bits.Add(uint(x), uint(y), uint(c))
		return uint64(a), uint64(b)
	}
	sub := func(x, y, c uint64) (uint64, uint64) {
		a, b := bits.Sub(uint(x), uint(y), uint(c))
		return uint64(a), uint64(b)
	}
	add32 := func(x, y, c uint64) (uint64, uint64) {
		a, b := bits.Add32(uint32(x), uint32(y), uint32(c))
		return uint64(a), uint64(b)
	}
	sub32 := func(x, y, c uint64) (uint64, uint64) {
		a, b := bits.Sub32(uint32(x), uint32(y), uint32(c))
		return uint64(a), uint64(b)
	}
	check("Add", bits.UintSize, add, true)
	check("Add32", 32, add32, true)
	check("Add64", 64, bits.Add64, true)
	check("Sub", bits.UintSize, sub, false)
	check("Sub32", 32, sub32, false)
	check("Sub64", 64, bits.Sub64, false)
}

func check(name string, width int, native op, add bool) {
	high := uint64(1) << (width - 1)
	maxWord := high | (high - 1)
	withoutFlag := func(x, y, _ uint64) (uint64, uint64) { return native(x, y, 0) }
	if add {
		tests := []fixture{
			{0, 0, 0, 0, 0}, {0, 0, 1, 1, 0}, {maxWord, 0, 1, 0, 1},
			{maxWord, maxWord, 1, maxWord, 1}, {maxWord, maxWord, 0, maxWord - 1, 1},
			{high, high, 0, 0, 1}, {high, 0, 0, high, 0}, {1, maxWord, 0, 0, 1},
		}
		run(name, native, tests, 0)
		run(name+" ignores carry-in", withoutFlag, tests, 3)
		run(name+" only checks both high bits", func(x, y, c uint64) (uint64, uint64) {
			word, _ := native(x, y, c)
			return word, (x & y) >> (width - 1)
		}, tests, 2)
		return
	}
	tests := []fixture{
		{0, 0, 0, 0, 0}, {0, 0, 1, maxWord, 1}, {0, maxWord, 1, 0, 1},
		{maxWord, maxWord, 1, maxWord, 1}, {maxWord, 0, 1, maxWord - 1, 0},
		{0, 1, 0, maxWord, 1}, {high, high, 0, 0, 0}, {high, 1, 0, high - 1, 0},
	}
	run(name, native, tests, 0)
	run(name+" ignores borrow-in", withoutFlag, tests, 4)
	run(name+" wraps the subtrahend before testing borrow", func(x, y, b uint64) (uint64, uint64) {
		word, _ := native(x, y, b)
		if x < (y+b)&maxWord {
			return word, 1
		}
		return word, 0
	}, tests, 2)
}

func run(name string, fn op, tests []fixture, want int) {
	failures := 0
	for _, t := range tests {
		word, out := fn(t.x, t.y, t.in)
		if word != t.word || out != t.out {
			failures++
			if want == 0 {
				log.Fatalf("%s(%#x,%#x,%d)=(%#x,%d), want (%#x,%d)", name, t.x, t.y, t.in, word, out, t.word, t.out)
			}
		}
	}
	if failures != want {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, failures, want)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, failures, len(tests))
}
