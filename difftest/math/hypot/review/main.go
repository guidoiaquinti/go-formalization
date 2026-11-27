package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const (
		inf  = uint64(0x7ff0000000000000)
		nan  = inf + 1
		sign = uint64(1) << 63
	)
	bits := math.Float64bits
	cases := []struct {
		x, y, fused, separate uint64
		rejectUnscaled        bool
	}{
		{bits(3), bits(-4), bits(5), bits(5), false},
		{sign, sign, 0, 0, false},
		{inf | sign, nan, inf, inf, true},
		{nan, inf, inf, inf, true},
		{nan, bits(1), nan, nan, false},
		{bits(math.MaxFloat64), 0, bits(math.MaxFloat64), bits(math.MaxFloat64), true},
		{1, 0, 1, 1, true},
		{0x3ff305267911c858, 0x3feaee7002cc1d39, 0x3ff74de5f2bff45e, 0x3ff74de5f2bff45f, false},
		{0x7fe6a09e667f3bcc, 0x7fe6a09e667f3bcc, inf, inf, false},
	}
	encode := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
	}
	var input bytes.Buffer
	var expected []string
	for _, c := range cases {
		fmt.Fprintln(&input, 0, c.x, c.y)
		fmt.Fprintln(&input, 1, c.x, c.y)
		expected = append(expected, encode(c.separate), encode(c.fused))
	}
	for i, row := range check.LeanRows("godiff-math-hypot", len(expected), &input) {
		if row != expected[i] {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/2, i%2, row, expected[i])
		}
	}
	var unscaledFailures, exchangedProfileFailures int
	for _, c := range cases {
		x, y := math.Float64frombits(c.x), math.Float64frombits(c.y)
		want, exchanged := c.separate, c.fused
		if runtime.GOARCH == "arm64" {
			want, exchanged = exchanged, want
		}
		got := encode(bits(math.Hypot(x, y)))
		if got != encode(want) || got != encode(bits(math.Hypot(y, x))) {
			log.Fatalf("Hypot(%#x, %#x) = %s, want %s in both orders", c.x, c.y, got, encode(want))
		}
		if got != encode(exchanged) {
			exchangedProfileFailures++
		}
		if c.rejectUnscaled && encode(bits(math.Sqrt(x*x+y*y))) != got {
			unscaledFailures++
		}
	}
	if unscaledFailures != 4 || exchangedProfileFailures != 1 {
		log.Fatalf("unscaled failures=%d, exchanged profile failures=%d", unscaledFailures, exchangedProfileFailures)
	}
	fmt.Printf("Hypot passed %d fixtures; unscaled formula failed %d selected probes; exchanging profiles failed %d\n",
		len(cases), unscaledFailures, exchangedProfileFailures)
}
