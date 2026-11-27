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
	const sign = uint64(1) << 63
	const inf = uint64(0x7ff0000000000000)
	bits := math.Float64bits
	type fixture struct {
		word  uint64
		want  [2]uint64
		power bool
	}
	same := func(word, want uint64) fixture {
		return fixture{word, [2]uint64{want, want}, false}
	}
	cases := []fixture{
		same(0, sign|inf), same(sign, sign|inf),
		same(inf, inf), same(sign|inf, inf+1),
		same(inf+1, inf+1), same(sign|inf|2, inf+1),
		same(bits(-1), inf+1), same(3, 0xc090c1a8ff971811),
		same(6, 0xc090bda8ff971811), same(0x3fefffffffffffff, 0xbca71547652b82fe),
		{0x3ff0000000000001, [2]uint64{0x3cb8000000000000, 0x3cba82b75b40966d}, false},
	}
	for k := -1074; k <= 1023; k++ {
		want := bits(float64(k))
		cases = append(cases, fixture{bits(math.Ldexp(1, k)), [2]uint64{want, want}, true})
	}

	observe := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
	}
	profile := 0
	if runtime.GOARCH == "arm64" {
		profile = 1
	}

	var inputs bytes.Buffer
	var wrongExponentFailures, exchangedFailures int
	for _, c := range cases {
		x := math.Float64frombits(c.word)
		got := observe(bits(math.Log2(x)))
		if want := observe(c.want[profile]); got != want {
			log.Fatalf("Log2(%#x)=%s want %s", c.word, got, want)
		}
		if c.power {
			_, exponent := math.Frexp(x)
			if got != observe(bits(float64(exponent))) {
				wrongExponentFailures++
			}
		}
		if observe(c.want[1-profile]) != got {
			exchangedFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode, c.word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-log2", 2*len(cases), &inputs) {
		if want := observe(cases[i/2].want[i%2]); row != want {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/2, i%2, row, want)
		}
	}
	if wrongExponentFailures != 2098 || exchangedFailures != 1 {
		log.Fatalf("unadjusted exponent failed %d, exchanged profile failed %d", wrongExponentFailures, exchangedFailures)
	}
	fmt.Printf("Log2 passed %d fixtures; using the unadjusted Frexp exponent failed %d; exchanging profiles failed %d\n",
		len(cases), wrongExponentFailures, exchangedFailures)
}
