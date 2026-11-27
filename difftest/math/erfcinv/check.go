package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	profile := nativeProfile()
	edges := binary64.Boundaries()
	for _, center := range []float64{math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64,
		0x1p-54, -0x1p-53, 0.15, 1, 1.85, 2, -1} {
		edges = append(edges, binary64.Load(math.Nextafter(center, math.Inf(-1))), binary64.Load(center),
			binary64.Load(math.Nextafter(center, math.Inf(1))))
	}
	check.CompareSamples("godiff-math-erfcinv", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Erfcinv(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x) || math.IsInf(x, 0):
			return "NaN"
		case x == 0:
			return observe(math.Inf(1))
		case x == 1:
			return "0"
		case x == 2:
			return observe(math.Inf(-1))
		}
		// Binary64 operands can span 1075 significant places in 1-x.
		// Subtract at higher precision, then round once to binary64.
		one := new(big.Float).SetPrec(2048).SetInt64(1)
		staged, _ := one.Sub(one, new(big.Float).SetFloat64(x)).Float64()
		if reference, actual := binary64.Load(staged), binary64.Load(1-x); reference != actual {
			log.Fatalf("1 - %#016x rounds to %#016x, reference %#016x", word, actual, reference)
		}
		return observe(math.Erfinv(staged))
	})
	checkLostInputs()
}

func checkLostInputs() {
	cases := []struct {
		input, stage uint64
		result       string
	}{
		{0x8000000000000001, 0x3ff0000000000000, observe(math.Inf(1))},
		{0xbca0000000000000, 0x3ff0000000000000, observe(math.Inf(1))},
		{0xbca0000000000001, 0x3ff0000000000001, "NaN"},
		{0x0000000000000001, 0x3ff0000000000000, observe(math.Inf(1))},
		{0x3c90000000000000, 0x3ff0000000000000, observe(math.Inf(1))},
	}
	for _, c := range cases {
		x := binary64.Store(c.input)
		got := math.Erfcinv(x)
		if staged := binary64.Load(1 - x); staged != c.stage || observe(got) != c.result {
			log.Fatalf("Erfcinv(%#016x): staged=%#016x result=%#016x, want %#016x and %s",
				c.input, staged, binary64.Load(got), c.stage, c.result)
		}
		fmt.Printf("Erfcinv(%#016x): staged=%#016x result=%#016x.\n", c.input, c.stage, binary64.Load(got))
	}
	x := binary64.Store(0x3c90000000000001)
	if got := math.Erfcinv(x); !(got > 0 && !math.IsInf(got, 0)) ||
		binary64.Load(1-x) != 0x3fefffffffffffff {
		log.Fatalf("Erfcinv(first positive input past cancellation) = %g", got)
	}
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 6 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load(2 * rng.Float64())
	case 2:
		centers := []float64{0x1p-54, -0x1p-53, 1, 2, 0.15, 1.85}
		center := binary64.Load(centers[rng.Intn(len(centers))])
		return center + uint64(rng.Intn(129)-64)
	case 3:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))
	case 4:
		return binary64.Load(math.Ldexp(rng.Float64(), -rng.Intn(60)))
	default:
		return binary64.Load(4*rng.Float64() - 1)
	}
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine Erfcinv build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Erfcinv architecture %s", runtime.GOARCH)
		return 0
	}
}
