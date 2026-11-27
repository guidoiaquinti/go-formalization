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
	type fixture struct {
		input uint64
		want  [3]uint64
	}
	same := func(input, want uint64) fixture {
		return fixture{input, [3]uint64{want, want, want}}
	}
	cases := []fixture{
		same(0, 0x3ff0000000000000), same(sign, 0x3ff0000000000000),
		same(inf, inf+1), same(inf|sign, inf+1), same(inf+7, inf+1),
		same(1, 0x3ff0000000000000), same(sign|1, 0x3ff0000000000000),
		same(0x3fe921fb54442d18, 0x3fe6a09e667f3bcd),
		same(0x3ff921fb54442d18, 0x3c91a62633145c00),
		same(0x400921fb54442d18, 0xbff0000000000000),
		same(sign|0x400921fb54442d18, 0xbff0000000000000),
		same(0x41bfffffffffffff, 0xbfee3edd238bb17a),
		same(0x41c0000000000000, 0xbfee3edd2dfeef90),
		same(0x41c0000000000001, 0xbfee3edd42e56b61),
		same(0x47f0000000000000, 0x3fcc841684311c82),
		same(0x5ff0000000000000, 0x3fc187699284d0b4),
		same(0x7fe0000000000000, 0xbfea719f26c232bf),
		same(0x7fefffffffffffff, 0xbfefffe62ecfab75),
	}
	observe := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
	}
	var inputs bytes.Buffer
	profile := nativeProfile()
	var reflectionFailures, reductionFailures int
	for _, c := range cases {
		x := math.Float64frombits(c.input)
		got := observe(math.Float64bits(math.Cos(x)))
		if want := observe(c.want[profile]); got != want {
			log.Fatalf("Cos(%#x)=%s want %s", c.input, got, want)
		}
		if got != observe(math.Float64bits(math.Abs(math.Cos(x)))) {
			reflectionFailures++
		}
		if x >= 0x1p29 && !math.IsInf(x, 0) && math.Abs(math.Cos(x)-math.Cos(math.Mod(x, 2*math.Pi))) > 0x1p-30 {
			reductionFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode, c.input)
		}
	}
	for i, row := range check.LeanRows("godiff-math-cos", 3*len(cases), &inputs) {
		if want := observe(cases[i/3].want[i%3]); row != want {
			log.Fatalf("fixture %d profile %d: Lean=%s want=%s", i/3, i%3, row, want)
		}
	}
	if reflectionFailures != 7 || reductionFailures != 6 {
		log.Fatalf("lost-octant-sign mutant failed %d; rounded-period mutant failed %d", reflectionFailures, reductionFailures)
	}
	fmt.Printf("Cos passed %d fixtures; lost octant sign fails %d; rounded-period reduction fails %d\n",
		len(cases), reflectionFailures, reductionFailures)
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine Cos build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unreviewed architecture %s", runtime.GOARCH)
		return 0
	}
}
