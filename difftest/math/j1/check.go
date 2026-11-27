package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/bessel"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

var centers = []float64{
	0x1p-27, 1, 2, 2.8571, 3 * math.Pi / 4, 3.831705970207512,
	4.5454, 7.015586669815619, 8, 20, 0x1p29, 0x1p129,
	math.MaxFloat64 / 2, math.MaxFloat64,
}

var tinyWitnesses = []struct{ word, want uint64 }{
	{0x8000000000000001, 0},
	{0x8000000000000002, 1},
	{0xbe30000000000000, 0x3e20000000000000},
	{0xbe3fffffffffffff, 0x3e2fffffffffffff},
	{0xbe40000000000000, 0xbe30000000000000},
	{0xbe40000000000001, 0xbe30000000000001},
}

func main() {
	profile := nativeProfile()
	edges := binary64.Boundaries()
	for _, fixture := range tinyWitnesses {
		if got := binary64.Load(math.J1(binary64.Store(fixture.word))); got != fixture.want {
			log.Fatalf("J1 pinned sign witness %#016x: got %#016x want %#016x", fixture.word, got, fixture.want)
		}
		edges = append(edges, fixture.word)
	}
	for _, center := range centers {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(x), binary64.Load(-x))
		}
	}
	var maximum float64
	check.CompareSamples("godiff-math-j1", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.J1(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x):
			return "NaN"
		case math.IsInf(x, 0):
			return "0"
		case x == 0:
			return "0"
		}
		magnitude := math.Abs(x)
		if magnitude < 0x1p-27 {
			// Go returns early on the magnitude, including for negative input.
			return observe(math.Ldexp(magnitude, -1))
		}
		got, want := math.J1(x), math.Copysign(1, x)*bessel.JOne(magnitude)
		scaled := math.Abs(got-want) * math.Sqrt(max(1, magnitude))
		if math.IsNaN(scaled) || scaled > 0x1p-45 {
			log.Fatalf("J1(%#016x) = %.17g, reference = %.17g, scaled error = %.5g", word, got, want, scaled)
		}
		maximum = max(maximum, scaled)
		return observe(got)
	})
	fmt.Printf("J1 independent reference: maximum scaled absolute error %.17g.\n", maximum)
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 6 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load(4*rng.Float64() - 2)
	case 2:
		word := binary64.Load(centers[rng.Intn(len(centers))])
		return (word + uint64(rng.Intn(129)-64)) | (rng.Uint64() & (1 << 63))
	case 3:
		return binary64.Load(40*rng.Float64() - 20)
	case 4:
		return binary64.Load(2048*rng.Float64() - 1024)
	default:
		return (rng.Uint64() % 0x3e40000000000000) | (rng.Uint64() & (1 << 63))
	}
}
func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine J1 build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled J1 architecture %s", runtime.GOARCH)
		return 0
	}
}
