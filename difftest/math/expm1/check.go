package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
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
			return fmt.Errorf("cannot determine Expm1 build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 1
			}
		}
	default:
		return fmt.Errorf("unmodeled Expm1 architecture %s", runtime.GOARCH)
	}

	edges := binary64.Boundaries()
	for _, value := range []float64{0x1p-54, 0x1p-53, 0x1p-28, 0.25, math.Ln2 / 2,
		math.Ln2, 1.5 * math.Ln2, 2 * math.Ln2, 19.5 * math.Ln2, 20 * math.Ln2,
		56 * math.Ln2, 56.5 * math.Ln2, 709.782712893384} {
		for _, x := range []float64{math.Nextafter(value, 0), value, math.Nextafter(value, math.Inf(1))} {
			edges = append(edges, binary64.Load(x), binary64.Load(-x))
		}
	}
	ln2 := highPrecisionLn2()
	var maxDistance uint64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-expm1", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return strconv.FormatUint(binary64.Load(math.Expm1(binary64.Store(word))), 10)
	}, func(word uint64) string {
		x := binary64.Store(word)
		want, exact := reference(word, ln2)
		if exact {
			fmt.Fprintln(&fixture, word, want, want, true)
			return strconv.FormatUint(want, 10)
		}
		// These public profiles select the generic helper. CompareSamples
		// checks that this word also equals the model for every row.
		got := binary64.Load(math.Expm1(x))
		fmt.Fprintln(&fixture, word, got, want, false)
		distance := max(got, want) - min(got, want)
		if distance > 1 {
			log.Fatalf("Expm1(%#016x) bits %#016x differ from reference %#016x by %d steps", word, got, want, distance)
		}
		maxDistance = max(maxDistance, distance)
		return strconv.FormatUint(got, 10)
	})
	fmt.Printf("Expm1 high-precision reference: maximum encoded-step distance %d.\n", maxDistance)
	dir, err := os.MkdirTemp("", "golean-math-expm1-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_EXPM1_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/expm1/testdata/"
	return check.TestInPackage("math", "TestGoLeanExpm1", fixtures+"bridge_test.go", fixtures+"expm1_test.go")
}

func reference(word uint64, ln2 *big.Float) (uint64, bool) {
	x := binary64.Store(word)
	switch {
	case math.IsNaN(x) || x == 0:
		return word, true
	case x >= 709.782712893384:
		return binary64.Load(math.Inf(1)), true
	case x <= -56*math.Ln2:
		return binary64.Load(-1), true
	default:
		return binary64.Load(idealExpm1(x, ln2)), false
	}
}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return r.Uint64()
	case 1:
		return binary64.Load(750*r.Float64() - 39)
	case 2:
		k := float64(r.Intn(1080) - 56)
		return binary64.Load((k+0.5)*math.Ln2) + uint64(r.Intn(9)-4)
	case 3:
		return binary64.Load(r.Float64() - 0.5)
	default:
		return binary64.Load(math.Ldexp(r.Float64()-0.5, -r.Intn(1074)))
	}
}

const precision = 256

func highPrecisionLn2() *big.Float {
	value := func(x float64) *big.Float { return new(big.Float).SetPrec(precision).SetFloat64(x) }
	third := value(1).Quo(value(1), value(3))
	ninth := value(1).Mul(third, third)
	power, sum := value(1).Set(third), value(0)
	term, divisor := value(0), value(0)
	// ln(2) = 2 * atanh(1/3).
	for n := int64(1); n <= 255; n += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(n)))
		power.Mul(power, ninth)
	}
	return sum.Mul(sum, value(2))
}

func idealExpm1(x float64, ln2 *big.Float) float64 {
	value := func(x float64) *big.Float { return new(big.Float).SetPrec(precision).SetFloat64(x) }
	k := math.Round(x * math.Log2E)
	r := value(x).Sub(value(x), value(k).Mul(value(k), ln2))
	term, sum, divisor := value(1), value(1), value(0)
	if k == 0 {
		// Sum from x without first rounding 1+x and losing the small term.
		sum.SetInt64(0)
	}
	for n := int64(1); n <= 64; n++ {
		term.Mul(term, r).Quo(term, divisor.SetInt64(n))
		sum.Add(sum, term)
	}
	if k != 0 {
		sum.SetMantExp(sum, int(k)).Sub(sum, value(1))
	}
	result, _ := sum.Float64()
	return result
}
