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

type input struct {
	order int
	word  uint64
}

var centers = []float64{0x1p-54, 0x1p-27, 1, 2, 2.197141326031017,
	5.429681040794135, 8, 20, 33, 34, 64, 0x1p29, 0x1p302, math.MaxFloat64}

func main() {
	profile := nativeProfile()
	edges := boundaries()
	const tolerance = 0x1p-40
	var maximum float64
	var frontierDifferences int
	check.CompareSamples("godiff-math-yn", edges, sample, func(c input) string {
		return fmt.Sprintf("%d %d %d", profile, c.order, c.word)
	}, func(c input) string {
		return observe(math.Yn(c.order, binary64.Store(c.word)))
	}, func(c input) string {
		x := binary64.Store(c.word)
		got := math.Yn(c.order, x)
		n := abs(c.order)
		reverse := c.order < 0 && n%2 != 0
		switch {
		case math.IsNaN(x), x < 0:
			return "NaN"
		case math.IsInf(x, 1):
			return "0"
		case x == 0:
			if reverse {
				return observe(math.Inf(1))
			}
			return observe(math.Inf(-1))
		}
		want, limitGap := bessel.YOrder(n, x, profile != 2)
		if reverse {
			want = -want
		}
		if got == want {
			return observe(got)
		}
		if math.IsInf(got, 0) != math.IsInf(want, 0) && math.Signbit(got) == math.Signbit(want) &&
			limitGap <= tolerance && (math.IsInf(got, 0) || 1-math.Abs(got)/math.MaxFloat64 <= tolerance) {
			frontierDifferences++
			return observe(got)
		}
		scaled := (math.Abs(got-want) / max(1, math.Abs(want))) * math.Sqrt(max(1, x))
		if math.IsNaN(scaled) || scaled > tolerance {
			log.Fatalf("Yn(%d, %#016x) = %.17g, reference = %.17g, normalized error = %.5g", c.order, c.word, got, want, scaled)
		}
		maximum = max(maximum, scaled)
		return observe(got)
	})
	fmt.Printf("Yn independent reference: maximum normalized absolute error %.17g, %d binary64 frontier classification differences.\n",
		maximum, frontierDifferences)
}

func boundaries() []input {
	var edges []input
	orders := []int{-64, -34, -33, -3, -2, -1, 0, 1, 2, 3, 33, 34, 64}
	for _, n := range orders {
		for _, word := range binary64.Boundaries() {
			edges = append(edges, input{n, word})
		}
		for _, center := range centers {
			for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
				edges = append(edges, input{n, binary64.Load(x)}, input{n, binary64.Load(-x)})
			}
		}
	}
	// Use neighboring words around each observed overflow frontier. The
	// native result selects the inputs only; the oracle checks them separately.
	for n := 1; n <= 64; n++ {
		low, high := uint64(1), binary64.Load(1)
		for high-low > 1 {
			mid := (low + high) / 2
			if math.IsInf(math.Yn(n, binary64.Store(mid)), -1) {
				low = mid
			} else {
				high = mid
			}
		}
		for word := low - 1; word <= high+1; word++ {
			edges = append(edges, input{n, word}, input{-n, word})
		}
	}
	return edges
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func sample(rng *rand.Rand, i int) input {
	n := rng.Intn(129) - 64
	var word uint64
	switch i % 6 {
	case 0:
		word = rng.Uint64()
	case 1:
		word = binary64.Load(4*rng.Float64() - 2)
	case 2:
		word = binary64.Load(centers[rng.Intn(len(centers))]) + uint64(rng.Intn(129)-64)
		word |= rng.Uint64() & (1 << 63)
	case 3:
		word = binary64.Load(144*rng.Float64() - 72)
	case 4:
		word = binary64.Load(math.Ldexp(1+rng.Float64(), -2-rng.Intn(40)))
	default:
		word = 1 + uint64(rng.Int63n((1<<52)-1))
	}
	return input{n, word}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine Yn build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Yn architecture %s", runtime.GOARCH)
		return 0
	}
}
