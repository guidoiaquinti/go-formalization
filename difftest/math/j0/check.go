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
	0x1p-27, 0x1p-13, 1, 2, 2.404825557695773, 2.8571,
	5 * math.Pi / 4, 4.5454, 5.520078110286311, 8, 20,
	0x1p29, 0x1p129, math.MaxFloat64 / 2, math.MaxFloat64,
}

func main() {
	profile := nativeProfile()
	edges := binary64.Boundaries()
	for _, center := range centers {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(x), binary64.Load(-x))
		}
	}
	var maximum float64
	check.CompareSamples("godiff-math-j0", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.J0(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x):
			return "NaN"
		case math.IsInf(x, 0):
			return "0"
		case x == 0:
			return observe(1)
		}
		magnitude := math.Abs(x)
		got, want := math.J0(magnitude), bessel.JZero(magnitude)
		scaled := math.Abs(got-want) * math.Sqrt(max(1, magnitude))
		if math.IsNaN(scaled) || scaled > 0x1p-45 {
			log.Fatalf("J0(%#016x) = %.17g, reference = %.17g, scaled error = %.5g", word, got, want, scaled)
		}
		maximum = max(maximum, scaled)
		return observe(got)
	})
	fmt.Printf("J0 independent reference: maximum scaled absolute error %.17g.\n", maximum)
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load(4*rng.Float64() - 2)
	case 2:
		word := binary64.Load(centers[rng.Intn(len(centers))])
		return (word + uint64(rng.Intn(129)-64)) | (rng.Uint64() & (1 << 63))
	case 3:
		return binary64.Load(40*rng.Float64() - 20)
	default:
		return binary64.Load(2048*rng.Float64() - 1024)
	}
}
func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine J0 build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled J0 architecture %s", runtime.GOARCH)
		return 0
	}
}
