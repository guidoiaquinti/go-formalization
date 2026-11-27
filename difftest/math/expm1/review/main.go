package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const sign = uint64(1) << 63
	const inf = uint64(0x7ff0000000000000)
	bits := math.Float64bits
	type fixture struct {
		word uint64
		want [3]uint64
	}
	same := func(word, want uint64) fixture {
		return fixture{word, [3]uint64{want, want, want}}
	}
	cases := []fixture{
		same(0, 0), same(sign, sign),
		same(1, 1), same(sign|1, sign|1),
		same(bits(0x1p-55), bits(0x1p-55)), same(bits(-0x1p-55), bits(-0x1p-55)),
		same(inf, inf), same(sign|inf, bits(-1)),
		same(inf+1, inf+1), same(sign|inf|0x8000000000005, sign|inf|0x8000000000005),
		same(0xc043687a9f1af2b1, bits(-1)),
		same(0x40862e42fefa39ef, inf),
		same(bits(-0.5), 0xbfd92e9a0720d3ec),
		same(bits(1.5), 0x400bda7fcc9f8a83),
		{0x4039e2ce850175b3, [3]uint64{0x4244548712ed6298, 0x4244548712ed6298, 0x4244548712ed6297}},
	}
	profile := 0
	if runtime.GOARCH == "arm64" {
		profile = 2
	} else {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("no build metadata")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 1
			}
		}
	}

	var inputs bytes.Buffer
	var directSubtractFailures, exchangedFailures int
	for _, c := range cases {
		x := math.Float64frombits(c.word)
		got := bits(math.Expm1(x))
		if want := c.want[profile]; got != want {
			log.Fatalf("Expm1(%#x)=%#x want %#x", c.word, got, want)
		}
		if math.Abs(x) < 0x1p-54 && bits(math.Exp(x)-1) != got {
			directSubtractFailures++
		}
		other := 2
		if profile == 2 {
			other = 0
		}
		if c.want[other] != got {
			exchangedFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode, c.word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-expm1", 3*len(cases), &inputs) {
		if want := strconv.FormatUint(cases[i/3].want[i%3], 10); row != want {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/3, i%3, row, want)
		}
	}
	if directSubtractFailures != 5 || exchangedFailures != 1 {
		log.Fatalf("direct subtraction failed %d, exchanged profile failed %d", directSubtractFailures, exchangedFailures)
	}
	fmt.Printf("Expm1 passed %d fixtures; Exp(x)-1 failed %d tiny probes; exchanging profiles failed %d\n",
		len(cases), directSubtractFailures, exchangedFailures)
}
