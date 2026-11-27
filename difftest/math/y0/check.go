package main

import (
	"bytes"
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
	0x1p-27, 0.8935769662791675, 1, 2, 2.8571, 3 * math.Pi / 4,
	3.9576784193148578, 4.5454, 8, 20, 0x1p29, 0x1p129,
	math.MaxFloat64 / 2, math.MaxFloat64,
}

func main() {
	checkSubnormalProfiles()
	profile := nativeProfile()
	edges := binary64.Boundaries()
	for _, center := range centers {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(x), binary64.Load(-x))
		}
	}
	var maximum float64
	check.CompareSamples("godiff-math-y0", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Y0(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x), x < 0:
			return "NaN"
		case x == 0:
			return observe(math.Inf(-1))
		case math.IsInf(x, 1):
			return "0"
		}
		got := math.Y0(x)
		maximum = max(maximum, checkReference(word, profile, got))
		return observe(got)
	})
	fmt.Printf("Y0 independent reference: maximum normalized absolute error %.17g.\n", maximum)
}

func checkReference(word uint64, profile int, got float64) float64 {
	x := binary64.Store(word)
	want := bessel.YZero(x, profile != 2)
	scaled := math.Abs(got-want) * math.Sqrt(max(1, x)) / max(1, math.Abs(want))
	if math.IsNaN(scaled) || scaled > 0x1p-45 {
		log.Fatalf("Y0 profile %d (%#016x) = %.17g, reference = %.17g, normalized error = %.5g",
			profile, word, got, want, scaled)
	}
	return scaled
}

func checkSubnormalProfiles() {
	words := []uint64{1, 2, 1 << 51, (1 << 52) - 1}
	const profiles = 3
	var input bytes.Buffer
	for profile := range profiles {
		for _, word := range words {
			fmt.Fprintln(&input, profile, word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-y0", profiles*len(words), &input) {
		got, err := strconv.ParseUint(row, 10, 64)
		if err != nil {
			log.Fatalf("invalid Y0 subnormal result %q: %v", row, err)
		}
		checkReference(words[i%len(words)], i/len(words), binary64.Store(got))
	}
	fmt.Printf("Y0 checked %d synthetic subnormal rows across three Lean profiles.\n", profiles*len(words))
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
		return binary64.Load(2 * rng.Float64())
	case 2:
		word := binary64.Load(centers[rng.Intn(len(centers))])
		return (word + uint64(rng.Intn(129)-64)) | (rng.Uint64() & (1 << 63))
	case 3:
		return binary64.Load(20 * rng.Float64())
	case 4:
		return binary64.Load(2048 * rng.Float64())
	default:
		return 1 + uint64(rng.Int63n((1<<52)-1))
	}
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine Y0 build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Y0 architecture %s", runtime.GOARCH)
		return 0
	}
}
