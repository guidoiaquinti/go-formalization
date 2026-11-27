package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const binary = "godiff-isnan"

func actual(repr uint64) uint64 {
	if math.IsNaN(math.Float64frombits(repr)) {
		return 1
	}
	return 0
}

func expected(repr uint64) uint64 {
	if (0x7ff0000000000000 < repr && repr < 0x8000000000000000) ||
		0xfff0000000000000 < repr {
		return 1
	}
	return 0
}

func main() {
	// On each side of the sign bit, exercise the finite-to-infinity and
	// infinity-to-NaN boundaries and both sides of the quiet-bit transition.
	var edges []uint64
	for _, sign := range []uint64{0, 1 << 63} {
		for _, encoding := range []uint64{
			0,
			1,
			0x7fefffffffffffff,
			0x7ff0000000000000,
			0x7ff0000000000001,
			0x7ff7ffffffffffff,
			0x7ff8000000000000,
			0x7ff8000000000001,
			0x7fffffffffffffff,
		} {
			edges = append(edges, sign|encoding)
		}
	}

	var input bytes.Buffer
	for _, repr := range edges {
		fmt.Fprintln(&input, repr)
	}
	for i, lean := range check.LeanRows(binary, len(edges), &input) {
		repr := edges[i]
		goResult, want := actual(repr), expected(repr)
		if goResult != want || lean != strconv.FormatUint(want, 10) {
			log.Fatalf("bits %#016x: Go=%d Lean=%s expected=%d", repr, goResult, lean, want)
		}
	}
	fmt.Printf("%s binary64 boundaries: %d inputs agree.\n", binary, len(edges))
	check.SampleUnary(binary, 64, actual, expected)
}
