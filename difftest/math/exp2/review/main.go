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
	var cases []fixture
	add := func(word, want uint64) {
		cases = append(cases, fixture{word, [3]uint64{want, want, want}})
	}
	for k := -1074; k <= 1023; k++ {
		if k < -1022 {
			add(bits(float64(k)), 1<<(k+1074))
		} else {
			add(bits(float64(k)), uint64(k+1023)<<52)
		}
	}
	add(sign, bits(1))
	add(inf, inf)
	add(sign|inf, 0)
	add(inf+1, inf+1)
	add(sign|inf|0x8000000000005, sign|inf|0x8000000000005)
	add(bits(-1075), 0)
	add(bits(-1074.75), 0)
	add(bits(-1074.5), 0)
	add(bits(math.Nextafter(-1074, math.Inf(-1))), 0)
	add(bits(1024), inf)
	add(bits(0.5), 0x3ff6a09e667f3bcc)
	cases = append(cases, fixture{0xc079d9c9b9f0787a, [3]uint64{0x2614f0ddbe62e0c1, 0x2614f0ddbe62e0c1, 0x2614f0ddbe62e0c0}})

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
	var input bytes.Buffer
	var changedCutoffFailures, exchangedFailures int
	for _, c := range cases {
		x := math.Float64frombits(c.word)
		got := bits(math.Exp2(x))
		if want := c.want[profile]; got != want {
			log.Fatalf("Exp2(%#x)=%#x want %#x", c.word, got, want)
		}
		laterCutoff := got
		if -1075 < x && x < -1074 {
			laterCutoff = 1
		}
		if laterCutoff != got {
			changedCutoffFailures++
		}
		other := 2
		if profile == 2 {
			other = 0
		}
		if c.want[other] != got {
			exchangedFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&input, mode, c.word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-exp2", 3*len(cases), &input) {
		if want := strconv.FormatUint(cases[i/3].want[i%3], 10); row != want {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/3, i%3, row, want)
		}
	}
	if changedCutoffFailures != 3 || exchangedFailures != 1 {
		log.Fatalf("later cutoff failed %d, exchanged profile failed %d", changedCutoffFailures, exchangedFailures)
	}
	fmt.Printf("Exp2 passed %d fixtures; later cutoff failed %d; exchanging profiles failed %d\n",
		len(cases), changedCutoffFailures, exchangedFailures)
}
