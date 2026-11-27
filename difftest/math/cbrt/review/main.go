package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const sign = uint64(1) << 63
	const inf = uint64(0x7ff0000000000000)
	const minNormal = uint64(0x0010000000000000)
	bits := math.Float64bits
	tests := []struct{ word, want uint64 }{
		{0, 0}, {sign, sign}, {inf, inf}, {sign | inf, sign | inf},
		{inf + 1, inf + 1}, {sign | inf | 0x8000000000005, sign | inf | 0x8000000000005},
		{1, bits(0x1p-358)}, {sign | 1, bits(-0x1p-358)},
		{8, bits(0x1p-357)}, {sign | 8, bits(-0x1p-357)},
		{bits(0.125), bits(0.5)}, {bits(-0.125), bits(-0.5)},
		{bits(27), bits(3)}, {bits(-27), bits(-3)},
		{bits(0x1p1023), bits(0x1p341)}, {bits(-0x1p1023), bits(-0x1p341)},
	}
	var inputs bytes.Buffer
	var lostSign, flushSubnormals int
	for _, c := range tests {
		fmt.Fprintln(&inputs, 0, c.word)
		fmt.Fprintln(&inputs, 1, c.word)
		actual := bits(math.Cbrt(math.Float64frombits(c.word)))
		if actual != c.want {
			log.Fatalf("Cbrt(%#x)=%#x want %#x", c.word, actual, c.want)
		}
		if actual&^sign != c.want {
			lostSign++
		}
		mutant := actual
		if c.word&^sign < minNormal {
			mutant = c.word & sign
		}
		if mutant != c.want {
			flushSubnormals++
		}
	}
	for i, row := range check.LeanRows("godiff-math-cbrt", 2*len(tests), &inputs) {
		if want := strconv.FormatUint(tests[i/2].want, 10); row != want {
			log.Fatalf("case %d profile %d: Lean=%s want=%s", i/2, i%2, row, want)
		}
	}
	if lostSign != 8 || flushSubnormals != 4 {
		log.Fatalf("lost sign failures=%d, flushed subnormal failures=%d", lostSign, flushSubnormals)
	}
	fmt.Printf("Cbrt passed %d fixtures; lost sign failed %d; flushed subnormals failed %d\n",
		len(tests), lostSign, flushSubnormals)
}
