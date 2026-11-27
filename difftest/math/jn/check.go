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

var centers = []float64{0x1p-29, 0x1p-27, 1, 2, 2.404825557695773,
	5.520078110286311, 8, 20, 33, 34, 64, 0x1p29, 0x1p302, math.MaxFloat64}

func main() {
	profile := nativeProfile()
	edges := boundaries()
	checkRootWitness(profile)
	var maximum float64
	var nearRoot int
	check.CompareSamples("godiff-math-jn", edges, sample, func(c input) string {
		return fmt.Sprintf("%d %d %d", profile, c.order, c.word)
	}, func(c input) string {
		return observe(math.Jn(c.order, binary64.Store(c.word)))
	}, func(c input) string {
		x := binary64.Store(c.word)
		got := math.Jn(c.order, x)
		switch {
		case math.IsNaN(x):
			return "NaN"
		case math.IsInf(x, 0):
			return "0"
		case x == 0:
			if c.order == 0 {
				return observe(1)
			}
			return "0"
		}
		n := c.order
		if n < 0 {
			n, x = -n, -x
		}
		magnitude := math.Abs(x)
		if n == 1 && magnitude < 0x1p-27 {
			return observe(math.Ldexp(magnitude, -1))
		}
		want := bessel.JOrder(n, magnitude)
		if math.Signbit(x) && n%2 != 0 {
			want = -want
		}
		scaled := math.Abs(got-want) * math.Sqrt(max(1, magnitude))
		if math.IsNaN(scaled) || scaled > 0x1p-35 {
			// Backward recurrence normalizes with J0(x) and its recurrence
			// estimate. Near a J0 zero the division can be ill-conditioned.
			if n > 1 && magnitude < float64(n) && math.Abs(bessel.JZero(magnitude)) <= 0x1p-32 {
				nearRoot++
				return observe(got)
			}
			log.Fatalf("Jn(%d, %#016x) = %.17g, reference = %.17g, scaled error = %.5g", c.order, c.word, got, want, scaled)
		}
		maximum = max(maximum, scaled)
		return observe(got)
	})
	fmt.Printf("Jn independent reference: maximum scaled absolute error %.17g, %d exact-model rows exceeded the numeric bound near J0 zeros.\n",
		maximum, nearRoot)
}

func checkRootWitness(profile int) {
	const center = uint64(0x40033d152e971b40)
	x := binary64.Store(center)
	want := bessel.JOrder(3, x)
	if math.Abs(want-0.19899990535769083) > 0x1p-50 || math.Abs(bessel.JZero(x)) > 0x1p-50 {
		log.Fatalf("Jn root reference at %#016x changed to %.17g", center, want)
	}
	if profile == 2 {
		words := []uint64{0x3fc69a0af29eb3ad, 0xfff0000000000000, 0x3fc66acb2c505f28}
		for i, expected := range words {
			arg := center - 1 + uint64(i)
			if got := binary64.Load(math.Jn(3, binary64.Store(arg))); got != expected {
				log.Fatalf("Jn pinned arm64 root witness %#016x: got %#016x want %#016x", arg, got, expected)
			}
		}
	}
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
	// Enter backward recurrence at the lowest step size near each n, where
	// the continued-fraction threshold generally takes longest to reach.
	for n := 2; n <= 64; n++ {
		x := math.Nextafter(float64(n), 0)
		edges = append(edges, input{n, binary64.Load(x)})
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
		word |= rng.Uint64() & (1 << 63)
	default:
		word = binary64.Load(float64(abs(n)) + 0.25*(rng.Float64()-0.5))
		word |= rng.Uint64() & (1 << 63)
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
			log.Fatal("cannot determine Jn build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Jn architecture %s", runtime.GOARCH)
		return 0
	}
}
