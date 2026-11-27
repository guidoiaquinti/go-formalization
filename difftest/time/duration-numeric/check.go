package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct{ duration, divisor int64 }

func native(c fixture) string {
	d, m := time.Duration(c.duration), time.Duration(c.divisor)
	return fmt.Sprintf("%d|%d|%d|%d|%d|%d|%d|%d|%d", d.Nanoseconds(), d.Microseconds(), d.Milliseconds(),
		math.Float64bits(d.Seconds()), math.Float64bits(d.Minutes()), math.Float64bits(d.Hours()), d.Truncate(m), d.Round(m), d.Abs())
}

func oracleFloat(n int64, scale int64) uint64 {
	whole, remainder := n/scale, n%scale
	fractional, _ := new(big.Rat).SetFrac64(remainder, scale).Float64()
	sum := new(big.Rat).Add(new(big.Rat).SetInt64(whole), new(big.Rat).SetFloat64(fractional))
	result, _ := sum.Float64()
	return math.Float64bits(result)
}

func oracle(c fixture) string {
	d, m := c.duration, c.divisor
	magnitude := uint64(d)
	negative := d < 0
	if negative {
		magnitude = 0 - magnitude
	}
	withSign := func(x uint64) int64 {
		if negative {
			return int64(0 - x)
		}
		return int64(x)
	}
	truncated, rounded := d, d
	if m > 0 {
		step := uint64(m)
		q, r := magnitude/step, magnitude%step
		truncated = withSign(q * step)
		qRound := q
		if r >= (step+1)/2 {
			qRound++
		}
		cap := uint64(math.MaxInt64)
		if negative {
			cap++
		}
		if qRound > cap/step {
			rounded = withSign(cap)
		} else {
			rounded = withSign(qRound * step)
		}
	}
	parts := []string{
		strconv.FormatInt(d, 10), strconv.FormatInt(withSign(magnitude/1000), 10), strconv.FormatInt(withSign(magnitude/1_000_000), 10),
		strconv.FormatUint(oracleFloat(d, 1_000_000_000), 10), strconv.FormatUint(oracleFloat(d, 60_000_000_000), 10),
		strconv.FormatUint(oracleFloat(d, 3_600_000_000_000), 10), strconv.FormatInt(truncated, 10), strconv.FormatInt(rounded, 10),
		strconv.FormatUint(min(magnitude, math.MaxInt64), 10),
	}
	return strings.Join(parts, "|")
}

func compare(label string, cases []fixture) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d\n", c.duration, c.divisor)
	}
	for i, got := range check.LeanRows("godiff-time-duration-numeric", len(cases), &wire) {
		c := cases[i]
		n, want := native(c), oracle(c)
		if got != want || n != want {
			log.Fatalf("%s row %d %+v: Lean %s Go %s want %s", label, i, c, got, n, want)
		}
	}
	fmt.Printf("time Duration %s: %d matched across nine declarations.\n", label, len(cases))
}

func main() {
	bound := []int64{math.MinInt64, math.MinInt64 + 1, -3_600_000_000_001, -60_000_000_001, -1_000_000_001,
		-1_000_001, -1001, -1, 0, 1, 3, 4, 5, 999, 1000, 1001, 999_999, 1_000_000, 1_000_001,
		999_999_999, 1_000_000_000, 1_000_000_001, 60_000_000_000, 3_600_000_000_000, math.MaxInt64 - 1, math.MaxInt64}
	var fixed []fixture
	for _, d := range bound {
		for _, m := range bound {
			fixed = append(fixed, fixture{d, m})
		}
	}
	for _, m := range []int64{2, 3, 10, 1_000_000, 3_600_000_000_000, math.MaxInt64 / 3} {
		for _, offset := range []int64{-1, 0, 1} {
			d := m + m/2 + offset
			fixed = append(fixed, fixture{d, m}, fixture{-d, m})
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				d, m := int64(rng.Uint64()), int64(rng.Uint64())
				if i%3 == 0 {
					m = bound[rng.Intn(len(bound))]
				}
				if i%7 == 0 {
					d = bound[rng.Intn(len(bound))]
				}
				cases[i] = fixture{d, m}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases)
		}
	}
}
