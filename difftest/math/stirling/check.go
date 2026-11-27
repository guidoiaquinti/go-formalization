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
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const splitThreshold = 143.01608

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	profile, err := nativeProfile()
	if err != nil {
		return err
	}
	pi, _, err := big.ParseFloat(
		"3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117068",
		10, 256, big.ToNearestEven)
	if err != nil {
		return err
	}
	sqrtTwoPi := precise(0).Sqrt(pi.SetMantExp(pi, 1))
	ln2 := atanhLog(precise(0).Quo(precise(1), precise(3)))
	edges := []float64{math.Inf(1), math.MaxFloat64, 33}
	for _, center := range []float64{34, 100, splitThreshold, 171.5, 172, 180, 200} {
		edges = append(edges, math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1)))
	}
	dir, err := os.MkdirTemp("", "golean-math-stirling-")
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
		var modelInput bytes.Buffer
		for _, x := range cases {
			fmt.Fprintln(&modelInput, profile, binary64.Load(x))
		}
		for i, model := range check.LeanRows("godiff-math-stirling", len(cases), &modelInput) {
			x := cases[i]
			first, second := idealStirling(x, ln2, sqrtTwoPi)
			if _, err := fmt.Fprintln(fixture, binary64.Load(x), model, binary64.Load(first), binary64.Load(second)); err != nil {
				return err
			}
		}
		fmt.Printf("Stirling %s: %d Lean pairs.\n", label, len(cases))
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
	if err := os.Setenv("GOLEAN_MATH_STIRLING_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/stirling/testdata/"
	return check.TestInPackage("math", "TestGoLeanStirling", data+"bridge_test.go", data+"stirling_test.go")
}

func sample(rng *rand.Rand, i int) float64 {
	switch i % 4 {
	case 0:
		return 33 + 167*rng.Float64()
	case 1:
		center := []float64{splitThreshold, 172, 180, 200}[rng.Intn(4)]
		word := binary64.Load(center) + uint64(rng.Intn(129)-64)
		return binary64.Store(word)
	case 2:
		return 100 + 100*rng.Float64()
	default:
		word := binary64.Load(33) + uint64(rng.Int63n(int64(binary64.Load(math.Inf(1))-binary64.Load(33))))
		return binary64.Store(word)
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func atanhLog(x *big.Float) *big.Float {
	square := precise(0).Mul(x, x)
	power, sum := precise(0).Set(x), precise(0)
	term, divisor := precise(0), precise(0)
	for k := int64(1); k <= 159; k += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(k)))
		power.Mul(power, square)
	}
	return sum.SetMantExp(sum, 1)
}

func logarithm(x float64, ln2 *big.Float) *big.Float {
	fraction := precise(0)
	k := precise(x).MantExp(fraction) - 1
	fraction.SetMantExp(fraction, 1)
	ratio := precise(0).Quo(precise(0).Sub(fraction, precise(1)), precise(0).Add(fraction, precise(1)))
	fractionLog := atanhLog(ratio)
	return fractionLog.Add(fractionLog, precise(0).Mul(precise(float64(k)), ln2))
}

func exponential(x, ln2 *big.Float) *big.Float {
	approximation, _ := x.Float64()
	k := math.Round(approximation * math.Log2E)
	remainder := precise(0).Sub(x, precise(0).Mul(precise(k), ln2))
	term, sum, divisor := precise(1), precise(1), precise(0)
	for n := int64(1); n <= 64; n++ {
		term.Mul(term, remainder)
		term.Quo(term, divisor.SetInt64(n))
		sum.Add(sum, term)
	}
	return sum.SetMantExp(sum, int(k))
}

func idealStirling(x float64, ln2, sqrtTwoPi *big.Float) (float64, float64) {
	if x > 200 {
		return math.Inf(1), 1
	}
	// The Bernoulli expansion of the logarithmic Stirling correction.
	reciprocal := precise(0).Quo(precise(1), precise(x))
	square := precise(0).Mul(reciprocal, reciprocal)
	coefficient := func(numerator, denominator int64) *big.Float {
		n := precise(0).SetInt64(numerator)
		return n.Quo(n, precise(0).SetInt64(denominator))
	}
	correction := coefficient(-691, 360360)
	for _, c := range [][2]int64{{1, 1188}, {-1, 1680}, {1, 1260}, {-1, 360}, {1, 12}} {
		correction.Mul(correction, square)
		correction.Add(correction, coefficient(c[0], c[1]))
	}
	correction.Mul(correction, reciprocal)
	second := exponential(correction, ln2)
	second.Mul(second, sqrtTwoPi)
	logX := logarithm(x, ln2)
	power := precise(0).Sub(precise(x), precise(0.5))
	var first *big.Float
	if x > splitThreshold {
		power.SetMantExp(power, -1)
		power.Mul(power, logX)
		first = exponential(power, ln2)
		remaining := precise(0).Sub(power, precise(x))
		second.Mul(second, exponential(remaining, ln2))
	} else {
		power.Mul(power, logX)
		power.Sub(power, precise(x))
		first = exponential(power, ln2)
	}
	y1, _ := first.Float64()
	y2, _ := second.Float64()
	return y1, y2
}

func nativeProfile() (int, error) {
	if runtime.GOARCH == "arm64" {
		return 4, nil
	}
	if runtime.GOARCH != "amd64" {
		return 0, fmt.Errorf("unmodeled Stirling architecture %s", runtime.GOARCH)
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return 0, fmt.Errorf("cannot determine Stirling build profile")
	}
	profile := 0
	for _, setting := range info.Settings {
		if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
			profile = 2
		}
	}
	dir, err := os.MkdirTemp("", "golean-stirling-exp-profile-")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "profile")
	if err := os.Setenv("GOLEAN_MATH_EXP_PROFILE", path); err != nil {
		return 0, err
	}
	const fixtures = "difftest/math/exp/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanExpProfile",
		fixtures+"profile_amd64_test.go", fixtures+"observe_profile_amd64_test.go"); err != nil {
		return 0, err
	}
	value, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	fused, err := strconv.ParseBool(string(value))
	if fused {
		profile++
	}
	return profile, err
}
