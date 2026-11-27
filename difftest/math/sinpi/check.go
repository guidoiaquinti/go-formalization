package main

import (
	"bufio"
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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	profile := nativeProfile()
	pi, _, err := big.ParseFloat(
		"3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117068",
		10, 256, big.ToNearestEven)
	if err != nil {
		return err
	}
	edges := []float64{0, math.Copysign(0, -1), math.MaxFloat64}
	for _, center := range []float64{math.SmallestNonzeroFloat64, 0x1p-1022, 0.25, 0.5, 0.75, 1,
		1.25, 1.5, 1.75, 2, 3, 0x1p50, 0x1p51, 0x1p52, 0x1p53} {
		edges = append(edges, math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1)))
	}
	dir, err := os.MkdirTemp("", "golean-math-sinpi-")
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
	compare := func(label string, cases []float64) error {
		var input bytes.Buffer
		for _, x := range cases {
			fmt.Fprintln(&input, profile, binary64.Load(x))
		}
		for i, row := range check.LeanRows("godiff-math-sinpi", len(cases), &input) {
			x := cases[i]
			if _, err := fmt.Fprintln(fixture, binary64.Load(x), row, binary64.Load(idealSinPi(x, pi))); err != nil {
				return err
			}
		}
		fmt.Printf("sinPi %s: %d Lean rows.\n", label, len(cases))
		return nil
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]float64, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_SINPI_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/sinpi/testdata/"
	return check.TestInPackage("math", "TestGoLeanSinPi", data+"bridge_test.go", data+"sinpi_test.go")
}

func sample(rng *rand.Rand, i int) float64 {
	switch i % 5 {
	case 0:
		return binary64.Store(rng.Uint64() % 0x7ff0000000000000)
	case 1:
		return 16 * rng.Float64()
	case 2:
		return math.Ldexp(rng.Float64(), rng.Intn(60))
	case 3:
		center := []float64{0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2, 0x1p52, 0x1p53}[rng.Intn(10)]
		return binary64.Store(binary64.Load(center) + uint64(rng.Intn(129)-64))
	default:
		return 0.25 * float64(rng.Intn(1<<20))
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func idealSinPi(x float64, pi *big.Float) float64 {
	if x == 0 {
		return -x
	}
	exact := precise(x)
	whole, _ := exact.Int(nil)
	fraction := precise(0).Sub(exact, precise(0).SetInt(whole))
	if fraction.Sign() == 0 {
		return math.Copysign(0, -1)
	}
	if fraction.Cmp(precise(0.5)) > 0 {
		fraction.Sub(precise(1), fraction)
	}
	angle := fraction.Mul(fraction, pi)
	square := precise(0).Mul(angle, angle)
	square.Neg(square)
	sum, term, divisor := precise(0).Set(angle), precise(0).Set(angle), precise(0)
	for n := int64(3); n <= 79; n += 2 {
		term.Mul(term, square)
		term.Quo(term, divisor.SetInt64(n*(n-1)))
		sum.Add(sum, term)
	}
	result, _ := sum.Float64()
	if whole.Bit(0) == 0 {
		return -result
	}
	return result
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine sinPi build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled sinPi architecture %s", runtime.GOARCH)
		return 0
	}
}
