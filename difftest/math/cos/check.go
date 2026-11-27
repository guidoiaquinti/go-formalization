package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
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
			return fmt.Errorf("cannot determine Cos build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 1
			}
		}
	default:
		return fmt.Errorf("unmodeled Cos architecture %s", runtime.GOARCH)
	}

	edges := binary64.Boundaries()
	for _, x := range []float64{0x1p-27, math.Pi / 4, math.Pi / 2, math.Pi, 3 * math.Pi / 2, 2 * math.Pi,
		0x1p29, 0x1p30, 0x1p64, 0x1p128, 0x1p512, 0x1p1023} {
		for _, value := range []float64{math.Nextafter(x, 0), x, math.Nextafter(x, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	halfPi := trig.HalfPi()
	var maxError float64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-cos", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Cos(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x) || math.IsInf(x, 0):
			fmt.Fprintln(&fixture, word, "NaN", 0, true)
			return "NaN"
		case x == 0:
			expected := observe(1)
			fmt.Fprintln(&fixture, word, expected, 0, true)
			return expected
		}
		got, want := math.Cos(x), trig.Cos(x, halfPi)
		if opposite := binary64.Load(math.Cos(-x)); opposite != binary64.Load(got) {
			log.Fatalf("Cos evenness failed for %#016x", word)
		}
		if err := math.Abs(got - want); math.IsNaN(err) || err > 0x1p-50 {
			log.Fatalf("Cos(%#016x) = %.17g, reference = %.17g, absolute error = %g", word, got, want, err)
		} else {
			maxError = max(maxError, err)
		}
		expected := observe(got)
		// The public wrapper selects cos on both modeled profiles, and
		// CompareSamples checks that this word also equals Lean's result.
		fmt.Fprintln(&fixture, word, expected, binary64.Load(want), false)
		return expected
	})
	fmt.Printf("Cos 1,408-bit reduction reference: maximum absolute difference %g.\n", maxError)
	dir, err := os.MkdirTemp("", "golean-math-cos-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_COS_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/cos/testdata/"
	return check.TestInPackage("math", "TestGoLeanCos", fixtures+"bridge_test.go", fixtures+"cos_test.go")
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
