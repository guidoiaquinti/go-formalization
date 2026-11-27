package main

import (
	"fmt"
	"log"
	"math/bits"
	"runtime"
)

type divide func(hi, lo, y uint64) (uint64, uint64)

type outcome struct {
	q, r      uint64
	panicText string
}

const (
	divideByZero = "runtime error: integer divide by zero"
	overflow     = "runtime error: integer overflow"
)

func observe(fn divide, hi, lo, y uint64) (out outcome) {
	defer func() {
		switch p := recover().(type) {
		case nil:
		case runtime.Error:
			out.panicText = p.Error()
		default:
			panic(p)
		}
	}()
	out.q, out.r = fn(hi, lo, y)
	return out
}

func check(name string, width uint, fn divide) {
	top := uint64(1) << (width - 1)
	max := top | (top - 1)
	fixtures := []struct {
		hi, lo, y uint64
		want      outcome
	}{
		{0, 0, 0, outcome{panicText: divideByZero}},
		{7, 1, 0, outcome{panicText: divideByZero}},
		{7, 1, 7, outcome{panicText: overflow}},
		{8, 1, 7, outcome{panicText: overflow}},
		{0, 17, 5, outcome{q: 3, r: 2}},
		{1, 1, 2, outcome{q: top, r: 1}},
		{max - 1, max, max, outcome{q: max, r: max - 1}},
	}
	var actualErrors, wrongPriority, dropHigh, allowEqual int
	for _, f := range fixtures {
		got := observe(fn, f.hi, f.lo, f.y)
		if got != f.want {
			actualErrors++
			log.Printf("%s(%#x, %#x, %#x) = %+v, want %+v", name, f.hi, f.lo, f.y, got, f.want)
		}
		priority := got
		if f.y <= f.hi {
			priority = outcome{panicText: overflow}
		}
		if priority != f.want {
			wrongPriority++
		}
		if observe(fn, 0, f.lo, f.y) != f.want {
			dropHigh++
		}
		equal := got
		if f.y != 0 && f.y == f.hi {
			equal = outcome{r: f.lo % f.y}
		}
		if equal != f.want {
			allowEqual++
		}
	}
	fmt.Printf("%s mismatches out of %d fixtures: Go %d, overflow first %d, ignore high word %d, permit equality %d\n",
		name, len(fixtures), actualErrors, wrongPriority, dropHigh, allowEqual)
	if actualErrors != 0 || wrongPriority != 2 || dropHigh != 4 || allowEqual != 1 {
		log.Fatalf("%s comparison or mutation check failed", name)
	}
}

func main() {
	if bits.UintSize != 64 {
		log.Fatal("bits.Div control requires 64-bit uint")
	}
	check("Div", 64, func(hi, lo, y uint64) (uint64, uint64) {
		q, r := bits.Div(uint(hi), uint(lo), uint(y))
		return uint64(q), uint64(r)
	})
	check("Div32", 32, func(hi, lo, y uint64) (uint64, uint64) {
		q, r := bits.Div32(uint32(hi), uint32(lo), uint32(y))
		return uint64(q), uint64(r)
	})
	check("Div64", 64, bits.Div64)
}
