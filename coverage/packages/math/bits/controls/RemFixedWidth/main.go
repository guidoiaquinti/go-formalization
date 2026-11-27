package main

import (
	"fmt"
	"log"
	"math/bits"
	"runtime"
)

type calculate func(hi, lo, y uint64) uint64

type result struct {
	value     uint64
	panicText string
}

func observe(f calculate, hi, lo, y uint64) (out result) {
	defer func() {
		if p := recover(); p != nil {
			if failure, ok := p.(runtime.Error); ok {
				out.panicText = failure.Error()
			} else {
				out.panicText = fmt.Sprintf("non-runtime panic %T", p)
			}
		}
	}()
	return result{value: f(hi, lo, y)}
}

func main() {
	run(32, func(hi, lo, y uint64) uint64 {
		return uint64(bits.Rem32(uint32(hi), uint32(lo), uint32(y)))
	}, func(hi, lo, y uint64) uint64 {
		_, rem := bits.Div32(uint32(hi), uint32(lo), uint32(y))
		return uint64(rem)
	})
	run(64, bits.Rem64, func(hi, lo, y uint64) uint64 {
		_, rem := bits.Div64(hi, lo, y)
		return rem
	})
}

func run(width uint, rem, rawDiv calculate) {
	maxWord := (uint64(1) << width) - 1
	zero := result{panicText: "runtime error: integer divide by zero"}
	cases := []struct {
		hi, lo, y uint64
		want      result
	}{
		{0, 0, 0, zero},
		{maxWord, maxWord, 0, zero},
		{0, 7, 3, result{value: 1}},
		{2, 0, 3, result{value: 2}},
		{3, 2, 3, result{value: 2}},
		{4, 0, 3, result{value: 1}},
		{maxWord, maxWord, 2, result{value: 1}},
		{maxWord, 0, maxWord, result{}},
	}
	var rawFailures, loFailures, zeroFailures int
	for _, c := range cases {
		got := observe(rem, c.hi, c.lo, c.y)
		if got != c.want {
			log.Fatalf("Rem%d(%d,%d,%d)=%+v want %+v", width, c.hi, c.lo, c.y, got, c.want)
		}
		if observe(rawDiv, c.hi, c.lo, c.y) != c.want {
			rawFailures++
		}
		if observe(func(_, lo, y uint64) uint64 { return lo % y }, c.hi, c.lo, c.y) != c.want {
			loFailures++
		}
		zeroFallback := func(hi, lo, y uint64) uint64 {
			if y == 0 {
				return 0
			}
			return rem(hi, lo, y)
		}
		if observe(zeroFallback, c.hi, c.lo, c.y) != c.want {
			zeroFailures++
		}
	}
	if rawFailures != 4 || loFailures != 2 || zeroFailures != 2 {
		log.Fatalf("Rem%d mutant failures raw=%d low=%d zero=%d", width, rawFailures, loFailures, zeroFailures)
	}
	fmt.Printf("Rem%d: %d fixtures pass, mutant failures raw division=%d low-only=%d zero fallback=%d\n",
		width, len(cases), rawFailures, loFailures, zeroFailures)
}
