package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/trig"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	profile := 0
	switch runtime.GOARCH {
	case "arm64":
		profile = 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return fmt.Errorf("cannot determine Sincos build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 1
			}
		}
	default:
		return fmt.Errorf("unmodeled Sincos architecture %s", runtime.GOARCH)
	}

	edges := binary64.Boundaries()
	for _, x := range []float64{0x1p-27, math.Pi / 4, math.Pi / 2, math.Pi, 3 * math.Pi / 2, 2 * math.Pi, 5 * math.Pi / 4, 7 * math.Pi / 4,
		0x1p29, 0x1p30, 0x1p64, 0x1p128, 0x1p512, 0x1p1023} {
		for _, value := range []float64{math.Nextafter(x, 0), x, math.Nextafter(x, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	halfPi := trig.HalfPi()
	var maxSinError, maxCosError float64
	check.CompareSamples("godiff-math-sincos", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observePair(math.Sincos(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x) || math.IsInf(x, 0):
			return "NaN NaN"
		case x == 0:
			return observePair(x, 1)
		}
		sine, cosine := math.Sincos(x)
		wantSine, wantCosine := trig.Sin(x, halfPi), trig.Cos(x, halfPi)
		oppositeSine, oppositeCosine := math.Sincos(-x)
		if binary64.Load(oppositeSine) != binary64.Load(sine)^(1<<63) ||
			binary64.Load(oppositeCosine) != binary64.Load(cosine) {
			log.Fatalf("Sincos sign symmetry failed for %#016x", word)
		}
		sinError := math.Abs(sine - wantSine)
		cosError := math.Abs(cosine - wantCosine)
		if math.IsNaN(sinError) || math.IsNaN(cosError) || max(sinError, cosError) > 0x1p-50 {
			log.Fatalf("Sincos(%#016x): sine absolute error %g, cosine absolute error %g", word, sinError, cosError)
		}
		maxSinError = max(maxSinError, sinError)
		maxCosError = max(maxCosError, cosError)
		return observePair(sine, cosine)
	})
	fmt.Printf("Sincos 1,408-bit reduction reference: maximum absolute errors sine=%g cosine=%g.\n", maxSinError, maxCosError)
	return nil
}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return r.Uint64()
	case 1:
		return binary64.Load((2*r.Float64() - 1) * 0x1p29)
	case 2:
		return binary64.Load(math.Ldexp(r.Float64()+0.5, 29+r.Intn(995))) | r.Uint64()&(1<<63)
	case 3:
		multiple := (r.Float64() - 0.5) * 0x1p29
		return binary64.Load(math.Round(multiple)*math.Pi/4) + uint64(r.Intn(9)-4)
	default:
		return binary64.Load(math.Ldexp(r.Float64()-0.5, -r.Intn(1074)))
	}
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func observePair(sine, cosine float64) string {
	return observe(sine) + " " + observe(cosine)
}
