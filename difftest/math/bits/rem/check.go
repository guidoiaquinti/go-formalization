// Exercise Rem at the host's uint width and its delegated implementation at the other width.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/bits"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func outcome(width int, hi, lo, y uint64) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(runtime.Error); ok && e.Error() == "runtime error: integer divide by zero" {
				result = "divide-by-zero"
			} else {
				result = fmt.Sprintf("unexpected panic (%T): %v", p, p)
			}
		}
	}()
	var r uint64
	switch {
	case width == bits.UintSize:
		r = uint64(bits.Rem(uint(hi), uint(lo), uint(y)))
	case width == 32:
		r = uint64(bits.Rem32(uint32(hi), uint32(lo), uint32(y)))
	default:
		r = bits.Rem64(hi, lo, y)
	}
	return strconv.FormatUint(r, 10)
}

type triple struct{ hi, lo, y uint64 }

func (c triple) expected(width int) string {
	if c.y == 0 {
		return "divide-by-zero"
	}
	n := new(big.Int).SetUint64(c.hi)
	n.Lsh(n, uint(width)).Add(n, new(big.Int).SetUint64(c.lo))
	return n.Mod(n, new(big.Int).SetUint64(c.y)).String()
}

func compare(width int, label string, cases []triple) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.hi, c.lo, c.y)
	}
	rows := check.LeanRows("godiff-rem", len(cases), &input, strconv.Itoa(width))
	for i, c := range cases {
		actual := outcome(width, c.hi, c.lo, c.y)
		expected := c.expected(width)
		model, source, ok := strings.Cut(rows[i], " ")
		if !ok || actual != expected || model != expected || source != expected {
			log.Fatalf("Rem/%d %s hi=%#x lo=%#x y=%#x: Go=%s, Lean abstraction=%q, Lean source=%q, expected=%s",
				width, label, c.hi, c.lo, c.y, actual, model, source, expected)
		}
	}
}

func run(width int) {
	mask := uint64(1)<<width - 1
	half, top := uint64(1)<<(width/2), uint64(1)<<(width-1)
	edges := []uint64{0, 1, 2, 3, half - 1, half, half + 1, top - 1, top, top + 1, mask - 1, mask}
	var boundaries []triple
	for _, hi := range edges {
		for _, lo := range edges {
			for _, y := range edges {
				boundaries = append(boundaries, triple{hi, lo, y})
			}
		}
	}
	compare(width, "boundaries", boundaries)
	count := len(boundaries)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]triple, 20_000)
		for i := range cases {
			c := triple{rng.Uint64() & mask, rng.Uint64() & mask, rng.Uint64() & mask}
			// Exercise hi just below, at, and just above the divisor.
			if offset := i % 8; offset < 3 {
				c.hi = (c.y + uint64(offset) - 1) & mask
			}
			cases[i] = c
		}
		compare(width, fmt.Sprintf("seed=%d", seed), cases)
		count += len(cases)
	}
	entry := "bits.Rem"
	if width != bits.UintSize {
		entry += strconv.Itoa(width)
	}
	fmt.Printf("%s (%d-bit inputs): %s (%s/%s) and Lean agree on %d cases, seeds %v.\n",
		entry, width, runtime.Version(), runtime.GOOS, runtime.GOARCH, count, seeds)
}

func main() {
	native := flag.String("native-width", "", "exercise only bits.Rem and require this Go uint width")
	flag.Parse()
	if *native != "" {
		if *native != strconv.Itoa(bits.UintSize) {
			log.Fatalf("expected %s-bit Go binary, got %d", *native, bits.UintSize)
		}
		run(bits.UintSize)
		return
	}
	run(32)
	run(64)
}
