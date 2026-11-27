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
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if runtime.GOARCH == "s390x" {
		fmt.Println("The modeled archLog profiles do not cover s390x assembly.")
		return
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	edges := binary64.Boundaries()
	for _, x := range []float64{0.5, math.Sqrt2 / 2, 1, math.Sqrt2, 2, 0x1p-1022, 0x1p-1021} {
		for _, word := range []float64{math.Nextafter(x, 0), x, math.Nextafter(x, math.Inf(1))} {
			edges = append(edges, binary64.Load(word), binary64.Load(-word))
		}
	}
	for k := -1074; k <= 1023; k++ {
		edges = append(edges, binary64.Load(math.Ldexp(1, k)))
	}
	ln2 := atanhLog(precise(1).Quo(precise(1), precise(3)))
	nativeAMD64 := runtime.GOARCH == "amd64"
	var maxDistance uint64
	var subnormals int
	var fixture bytes.Buffer
	compare := func(label string, cases []uint64) error {
		const executor = "godiff-math-arch-log"
		var amd64, fallback bytes.Buffer
		for _, word := range cases {
			fmt.Fprintln(&amd64, 0, word)
			fmt.Fprintln(&fallback, 1, word)
		}
		stubs := check.LeanRows(executor, len(cases), &fallback)
		for i, model := range check.LeanRows(executor, len(cases), &amd64) {
			word := cases[i]
			if stubs[i] != "panic:not implemented" {
				return fmt.Errorf("archLog stub(%#016x): Lean=%q", word, stubs[i])
			}
			x := binary64.Store(word)
			if want, exact := specialLog(x); exact {
				if model != want {
					return fmt.Errorf("archLog(%#016x): Lean=%q, exact reference=%q", word, model, want)
				}
			} else {
				got, err := strconv.ParseUint(model, 10, 64)
				if err != nil {
					return err
				}
				adjustExponent := 0
				if word < 1<<52 {
					// The assembly's leading bit gives (1+fraction) * 2^-1023.
					x = binary64.Store(word | 1<<52)
					adjustExponent = -1
					subnormals++
				}
				expected := binary64.Load(idealLog(x, adjustExponent, ln2))
				distance := max(got, expected) - min(got, expected)
				if distance > 1 {
					return fmt.Errorf("archLog(%#016x): Lean=%#016x, independent=%#016x, distance=%d", word, got, expected, distance)
				}
				maxDistance = max(maxDistance, distance)
			}
			if !nativeAMD64 {
				model = stubs[i]
			}
			fmt.Fprintf(&fixture, "%d %q\n", word, model)
		}
		fmt.Printf("archLog %s: %d modeled rows per profile.\n", label, len(cases))
		return nil
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]uint64, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	fmt.Printf("archLog amd64 model independent reference: maximum %d steps, %d effective-argument subnormal cases.\n", maxDistance, subnormals)
	dir, err := os.MkdirTemp("", "golean-math-arch-log-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_ARCH_LOG_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/arch-log/testdata/"
	return check.TestInPackage("math", "TestGoLeanArchLog", data+"bridge_test.go", data+"arch_log_test.go")
}

func specialLog(x float64) (string, bool) {
	switch {
	case math.IsNaN(x) || x < 0:
		return "NaN", true
	case x == 0:
		return observe(math.Inf(-1)), true
	case math.IsInf(x, 1):
		return observe(x), true
	case x == 1:
		return "0", true
	default:
		return "", false
	}
}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return r.Uint64()
	case 1:
		return r.Uint64() & (1<<52 - 1)
	case 2:
		return uint64(int64(binary64.Load(1)) + r.Int63n(8193) - 4096)
	case 3:
		k := r.Intn(2098) - 1074
		return binary64.Load(math.Ldexp(math.Sqrt2, k)) + uint64(r.Intn(9)-4)
	default:
		return r.Uint64() &^ (1 << 63)
	}
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

const precision = 256

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(precision).SetFloat64(value)
}

// Compute 2*atanh(s) with |s| at most 1/3.
func atanhLog(s *big.Float) *big.Float {
	square := precise(0).Mul(s, s)
	power, sum := precise(0).Set(s), precise(0)
	term, divisor := precise(0), precise(0)
	for n := int64(1); n <= 255; n += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(n)))
		power.Mul(power, square)
	}
	return sum.Mul(sum, precise(2))
}

func idealLog(x float64, adjustExponent int, ln2 *big.Float) float64 {
	u, mantissa := precise(x), precise(0)
	exponent := u.MantExp(mantissa)
	// Center the mantissa around 1 to avoid cancellation in the sum near 1.
	if mantissa.Cmp(precise(0.75)) < 0 {
		exponent--
		mantissa.SetMantExp(mantissa, 1)
	}
	ratio := precise(0).Quo(precise(0).Sub(mantissa, precise(1)), precise(0).Add(mantissa, precise(1)))
	result := atanhLog(ratio)
	exponent += adjustExponent
	result.Add(result, precise(float64(exponent)).Mul(precise(float64(exponent)), ln2))
	rounded, _ := result.Float64()
	return rounded
}
