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
	dir, err := os.MkdirTemp("", "golean-math-asinh-")
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
	if err := os.Setenv("GOLEAN_MATH_ASINH_CASES", path); err != nil {
		return err
	}
	exact := func(word uint64, expected string) string {
		fmt.Fprintln(fixture, word, expected, expected, "exact")
		return expected
	}

	profile := nativeProfile()
	edges := binary64.Boundaries()
	for _, center := range []float64{0x1p-28, 0.25, 1, 2, 0x1p28, 0x1p512, math.MaxFloat64} {
		for _, value := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	ln2Half := atanhSeries(precise(0).Quo(precise(1), precise(3)))
	var maximum uint64
	check.CompareSamples("godiff-math-asinh", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Asinh(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		magnitude := math.Abs(x)
		switch {
		case math.IsNaN(x):
			return exact(word, "NaN")
		case math.IsInf(x, 0):
			return exact(word, observe(x))
		case magnitude < 0x1p-28:
			return exact(word, observe(x))
		}
		got, want := math.Asinh(x), idealAsinh(x, ln2Half)
		a, b := binary64.Load(got), binary64.Load(want)
		distance := max(a, b) - min(a, b)
		if distance > 4 {
			log.Fatalf("Asinh(%#016x) = %#016x, reference = %#016x, distance = %d", word, a, b, distance)
		}
		maximum = max(maximum, distance)
		expected := observe(got)
		fmt.Fprintln(fixture, word, expected, observe(want), "numeric")
		return expected
	})
	fmt.Printf("Asinh 256-bit reference: maximum distance of %d binary64 steps.\n", maximum)
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	const data = "difftest/math/asinh/testdata/"
	return check.TestInPackage("math", "TestGoLeanAsinh", data+"bridge_test.go", data+"asinh_test.go")
}

func sample(rng *rand.Rand, i int) uint64 {
	const sign = uint64(1) << 63
	switch i % 6 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load(4*rng.Float64() - 2)
	case 2:
		center := []float64{0x1p-28, 2, 0x1p28, math.MaxFloat64}[rng.Intn(4)]
		return (binary64.Load(center) + uint64(rng.Intn(129)-64)) | (rng.Uint64() & sign)
	case 3:
		return (binary64.Load(math.MaxFloat64) - uint64(rng.Intn(1<<24))) | (rng.Uint64() & sign)
	case 4:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))
	default:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, rng.Intn(1024)))
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

func idealAsinh(x float64, ln2Half *big.Float) float64 {
	magnitude := math.Abs(x)
	argument := precise(magnitude)
	root := precise(0).Mul(argument, argument)
	root.Sqrt(root.Add(root, precise(1)))
	exponent := 0
	if magnitude <= 0.25 {
		// asinh(x) = atanh(x/sqrt(1+x*x)).
		argument.Quo(argument, root)
	} else {
		// asinh(x) = log(x+sqrt(1+x*x)); normalize before applying the log series.
		positive := precise(0).Add(argument, root)
		mantissa := precise(0)
		exponent = positive.MantExp(mantissa) - 1
		mantissa.SetMantExp(mantissa, 1)
		argument.Quo(precise(0).Sub(mantissa, precise(1)), precise(0).Add(mantissa, precise(1)))
	}
	sum := atanhSeries(argument)
	if magnitude > 0.25 {
		sum.Add(sum, precise(0).Mul(precise(float64(exponent)), ln2Half))
		sum.SetMantExp(sum, 1)
	}
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
			log.Fatal("cannot determine Asinh build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Asinh architecture %s", runtime.GOARCH)
		return 0
	}
}
