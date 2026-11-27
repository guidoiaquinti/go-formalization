package main

import (
	"bufio"
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
	dir, err := os.MkdirTemp("", "golean-math-atanh-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fixture := bufio.NewWriter(file)
	if err := os.Setenv("GOLEAN_MATH_ATANH_CASES", path); err != nil {
		return err
	}
	exact := func(word uint64, expected string) string {
		fmt.Fprintln(fixture, word, expected, expected, "exact")
		return expected
	}

	profile := nativeProfile()
	edges := binary64.Boundaries()
	for _, center := range []float64{0x1p-28, 0.125, 0.25, 0.5, 0.75, 0.875, 1} {
		for _, value := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	ln2Half := atanhSeries(precise(0).Quo(precise(1), precise(3)))
	var maximum uint64
	check.CompareSamples("godiff-math-atanh", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Atanh(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		magnitude := math.Abs(x)
		switch {
		case math.IsNaN(x) || magnitude > 1:
			return exact(word, "NaN")
		case magnitude == 1:
			return exact(word, observe(math.Copysign(math.Inf(1), x)))
		case magnitude < 0x1p-28:
			return exact(word, observe(x))
		}
		got, want := math.Atanh(x), idealAtanh(x, ln2Half)
		a, b := binary64.Load(got), binary64.Load(want)
		distance := max(a, b) - min(a, b)
		if distance > 4 {
			log.Fatalf("Atanh(%#016x) = %#016x, reference = %#016x, distance = %d", word, a, b, distance)
		}
		maximum = max(maximum, distance)
		expected := observe(got)
		fmt.Fprintln(fixture, word, expected, observe(want), "numeric")
		return expected
	})
	fmt.Printf("Atanh 256-bit reference: maximum distance of %d binary64 steps.\n", maximum)
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	const data = "difftest/math/atanh/testdata/"
	return check.TestInPackage("math", "TestGoLeanAtanh", data+"bridge_test.go", data+"atanh_test.go")
}

func sample(rng *rand.Rand, i int) uint64 {
	const sign = uint64(1) << 63
	switch i % 5 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load(2*rng.Float64() - 1)
	case 2:
		center := []float64{0x1p-28, 0.5, 1}[rng.Intn(3)]
		return (binary64.Load(center) + uint64(rng.Intn(129)-64)) | (rng.Uint64() & sign)
	case 3:
		return (binary64.Load(1) - uint64(rng.Intn(1<<24))) | (rng.Uint64() & sign)
	default:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func atanhSeries(x *big.Float) *big.Float {
	square := precise(0).Mul(x, x)
	power, sum := precise(0).Set(x), precise(0)
	term, divisor := precise(0), precise(0)
	for n := int64(1); n <= 255; n += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(n)))
		power.Mul(power, square)
	}
	return sum
}

func idealAtanh(x float64, ln2Half *big.Float) float64 {
	magnitude := math.Abs(x)
	argument := precise(magnitude)
	exponent := 0
	if magnitude > 0.25 {
		// (1+|x|)/(1-|x|) = 2^exponent*mantissa, with 1 <= mantissa < 2.
		ratio := precise(0).Quo(precise(0).Add(precise(1), argument), precise(0).Sub(precise(1), argument))
		mantissa := precise(0)
		exponent = ratio.MantExp(mantissa) - 1
		mantissa.SetMantExp(mantissa, 1)
		argument.Quo(precise(0).Sub(mantissa, precise(1)), precise(0).Add(mantissa, precise(1)))
	}
	sum := atanhSeries(argument)
	sum.Add(sum, precise(0).Mul(precise(float64(exponent)), ln2Half))
	rounded, _ := sum.Float64()
	return math.Copysign(rounded, x)
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
			log.Fatal("cannot determine Atanh build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Atanh architecture %s", runtime.GOARCH)
		return 0
	}
}
