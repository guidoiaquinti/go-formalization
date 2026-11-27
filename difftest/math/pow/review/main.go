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
	const nan = inf | 1
	bits := math.Float64bits
	type fixture struct {
		x, y uint64
		want [3]uint64
	}
	same := func(x, y, want uint64) fixture {
		return fixture{x, y, [3]uint64{want, want, want}}
	}
	one := bits(1)
	cases := []fixture{
		same(nan, 0, one),
		same(nan, sign, one),
		same(one, nan, one),
		same(nan, one, nan),
		same(bits(-3), one, bits(-3)),
		same(sign, bits(3), sign),
		same(sign, bits(-3), sign|inf),
		same(sign, bits(2), 0),
		same(sign, bits(-2), inf),
		same(sign|inf, bits(3), sign|inf),
		same(sign|inf, bits(-3), sign),
		same(bits(-2), bits(0.5), nan),
		same(bits(4), bits(0.5), bits(2)),
		same(bits(4), bits(-0.5), bits(0.5)),
		same(bits(-2), bits(3), bits(-8)),
		same(bits(0.5), bits(1074), 1),
		same(bits(2), bits(-1074), 1),
		same(bits(-1), bits(0x1p63), one),
		same(bits(0.5), bits(0x1p63), 0),
		same(bits(2), bits(0x1p63), inf),
		{bits(3), bits(0.75), [3]uint64{0x40023c6e3224f9d0, 0x40023c6e3224f9d0, 0x40023c6e3224f9d1}},
	}
	observe := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
	}
	profile := 0
	switch runtime.GOARCH {
	case "arm64":
		profile = 2
	case "amd64":
		// The fixtures have the same expected words for both amd64 Exp paths.
	default:
		log.Fatalf("unreviewed architecture %s", runtime.GOARCH)
	}
	var inputs bytes.Buffer
	var nanFirstFailures, unsignedZeroFailures, logShortcutFailures int
	for _, c := range cases {
		x, y := math.Float64frombits(c.x), math.Float64frombits(c.y)
		got := observe(bits(math.Pow(x, y)))
		if want := observe(c.want[profile]); want != got {
			log.Fatalf("Pow(%#x, %#x)=%s want %s", c.x, c.y, got, want)
		}
		if (math.IsNaN(x) || math.IsNaN(y)) && got != "NaN" {
			nanFirstFailures++
		}
		if c.x == sign && got != observe(bits(math.Pow(0, y))) {
			unsignedZeroFailures++
		}
		if x == -2 && y == 3 && got != observe(bits(math.Exp(y*math.Log(x)))) {
			logShortcutFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode, c.x, c.y)
		}
	}
	for i, row := range check.LeanRows("godiff-math-pow", 3*len(cases), &inputs) {
		if want := observe(cases[i/3].want[i%3]); want != row {
			log.Fatalf("fixture %d profile %d: Lean=%s want=%s", i/3, i%3, row, want)
		}
	}
	if nanFirstFailures != 3 || unsignedZeroFailures != 2 || logShortcutFailures != 1 {
		log.Fatalf("NaN-first failed %d; unsigned-zero failed %d; log/exp shortcut failed %d",
			nanFirstFailures, unsignedZeroFailures, logShortcutFailures)
	}
	fmt.Printf("Pow passed %d fixtures; NaN-first failed %d; unsigned-zero failed %d; log/exp shortcut failed %d\n",
		len(cases), nanFirstFailures, unsignedZeroFailures, logShortcutFailures)
}
