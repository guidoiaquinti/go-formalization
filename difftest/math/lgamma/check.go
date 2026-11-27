package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const ymin = 1.461632144968362245

func main() {
	profile := nativeProfile()
	pi, _, err := big.ParseFloat(
		"3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117068",
		10, 256, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)
	}
	ln2 := atanhLog(precise(0).Quo(precise(1), precise(3)))
	logSqrtTwoPi := logarithm(precise(0).SetMantExp(pi, 1), ln2)
	logSqrtTwoPi.SetMantExp(logSqrtTwoPi, -1)
	edges := binary64.Boundaries()
	for _, center := range []float64{math.SmallestNonzeroFloat64, 0x1p-70, ymin - 1 - 0.27, ymin - 1 + 0.27, 0.9, 1,
		ymin - 0.27, ymin, ymin + 0.27, 2, 3, 4, 5, 6, 7, 8, 0x1p52, 0x1p58, 1e305, math.MaxFloat64} {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(x), binary64.Load(-x))
		}
	}
	checkSubnormalProfiles(ln2)
	var maximum float64
	var stagedSubnormals int
	check.CompareSamples("godiff-math-lgamma", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Lgamma(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x):
			return "NaN 1"
		case math.IsInf(x, 0):
			return observe(x, 1)
		case x == 0 || x < 0 && math.Trunc(x) == x:
			return observe(math.Inf(1), 1)
		case x == 1 || x == 2:
			return "0 1"
		}
		value, sign := math.Lgamma(x)
		if magnitude := math.Abs(x); profile != 2 && magnitude < 0x1p-1022 {
			// Lgamma returns the rounded public Log of the magnitude on this branch.
			if binary64.Load(value) != binary64.Load(-math.Log(magnitude)) {
				log.Fatalf("Lgamma(%#016x)=%g disagrees with staged -Log", word, value)
			}
			want, wantSign := subnormalReference(word, profile, ln2)
			checkSubnormalResult(word, profile, binary64.Load(value), sign, want, wantSign)
			stagedSubnormals++
			return observe(value, sign)
		}
		want, wantSign := idealLgamma(x, ln2, pi, logSqrtTwoPi)
		if math.IsNaN(value) || math.IsInf(value, 0) != math.IsInf(want, 0) || sign != wantSign {
			log.Fatalf("Lgamma(%#016x)=(%g,%d) reference=(%g,%d)", word, value, sign, want, wantSign)
		}
		if !math.IsInf(value, 0) {
			delta := math.Abs(value-want) / max(1, math.Abs(want))
			if delta > 0x1p-42 {
				log.Fatalf("Lgamma(%#016x)=%g reference=%g normalized difference %g", word, value, want, delta)
			}
			maximum = max(maximum, delta)
		}
		return observe(value, sign)
	})
	fmt.Printf("Lgamma 256-bit reference maximum normalized difference %g; %d native amd64 subnormal staged-log cases.\n",
		maximum, stagedSubnormals)
}

func subnormalReference(word uint64, profile int, ln2 *big.Float) (uint64, int) {
	magnitude := word &^ (uint64(1) << 63)
	effectiveInput := precise(binary64.Store(magnitude))
	if profile != 2 {
		// The amd64 assembly supplies the leading mantissa bit at exponent
		// -1023. Dividing the corresponding normal value by two at high
		// precision preserves the extra bit at the subnormal boundary.
		effectiveInput.SetFloat64(binary64.Store(magnitude | (1 << 52)))
		effectiveInput.SetMantExp(effectiveInput, -1)
	}
	value, _ := logarithm(effectiveInput, ln2).Float64()
	sign := 1
	if word>>63 != 0 {
		sign = -1
	}
	return binary64.Load(-value), sign
}

func checkSubnormalResult(word uint64, profile int, got uint64, sign int, want uint64, wantSign int) {
	if distance := max(got, want) - min(got, want); distance > 1 || sign != wantSign {
		log.Fatalf("Lgamma profile %d subnormal %#016x: got %#016x %d, effective reference %#016x %d",
			profile, word, got, sign, want, wantSign)
	}
}

func checkSubnormalProfiles(ln2 *big.Float) {
	words := []uint64{1, (1 << 63) | 1, (1 << 52) - 1, (1 << 63) | ((1 << 52) - 1)}
	const profiles = 3
	var input bytes.Buffer
	for profile := range profiles {
		for _, word := range words {
			fmt.Fprintln(&input, profile, word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-lgamma", profiles*len(words), &input) {
		profile, word := i/len(words), words[i%len(words)]
		value, signText, ok := strings.Cut(row, " ")
		got, valueErr := strconv.ParseUint(value, 10, 64)
		sign, signErr := strconv.Atoi(signText)
		if !ok || valueErr != nil || signErr != nil {
			log.Fatalf("invalid Lgamma witness row %q", row)
		}
		want, wantSign := subnormalReference(word, profile, ln2)
		checkSubnormalResult(word, profile, got, sign, want, wantSign)
	}
	fmt.Printf("Lgamma checked %d synthetic subnormal rows across three Lean profiles.\n", profiles*len(words))
}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 6 {
	case 0:
		return r.Uint64()
	case 1:
		return binary64.Load(4*r.Float64() - 2)
	case 2:
		return binary64.Load(260*r.Float64() - 200)
	case 3:
		cutoffs := []float64{0x1p-70, ymin - 1 - 0.27, ymin - 1 + 0.27, 0.9, ymin - 0.27, ymin + 0.27, 2, 8, 0x1p52, 0x1p58}
		word := binary64.Load(cutoffs[r.Intn(len(cutoffs))]) + uint64(r.Intn(129)-64)
		return word | (r.Uint64() & (1 << 63))
	case 4:
		return binary64.Load(math.Ldexp(r.Float64()-0.5, -r.Intn(1074)))
	default:
		n := 1 + r.Intn(200)
		return binary64.Load(-float64(n)) + uint64(r.Intn(129)-64)
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

func logarithm(x, ln2 *big.Float) *big.Float {
	fraction := precise(0)
	k := x.MantExp(fraction) - 1
	fraction.SetMantExp(fraction, 1)
	ratio := precise(0).Quo(precise(0).Sub(fraction, precise(1)), precise(0).Add(fraction, precise(1)))
	fractionLog := atanhLog(ratio)
	return fractionLog.Add(fractionLog, precise(0).Mul(precise(float64(k)), ln2))
}

func correction(x *big.Float) *big.Float {
	reciprocal := precise(0).Quo(precise(1), x)
	square := precise(0).Mul(reciprocal, reciprocal)
	coef := func(n, d int64) *big.Float {
		return precise(0).Quo(precise(0).SetInt64(n), precise(0).SetInt64(d))
	}
	sum := coef(-691, 360360)
	for _, c := range [][2]int64{{1, 1188}, {-1, 1680}, {1, 1260}, {-1, 360}, {1, 12}} {
		sum.Mul(sum, square)
		sum.Add(sum, coef(c[0], c[1]))
	}
	return sum.Mul(sum, reciprocal)
}

func logGammaPositive(x float64, ln2, logSqrtTwoPi *big.Float) *big.Float {
	argument, product, one, minimum := precise(x), precise(1), precise(1), precise(36)
	for argument.Cmp(minimum) < 0 {
		product.Mul(product, argument)
		argument.Add(argument, one)
	}
	result := precise(0).Sub(argument, precise(0.5))
	result.Mul(result, logarithm(argument, ln2))
	result.Sub(result, argument)
	result.Add(result, correction(argument))
	result.Add(result, logSqrtTwoPi)
	return result.Sub(result, logarithm(product, ln2))
}

func sinPiMagnitude(x, floor float64, pi *big.Float) *big.Float {
	fraction := precise(0).Sub(precise(x), precise(floor))
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
	return sum
}

func idealLgamma(x float64, ln2, pi, logSqrtTwoPi *big.Float) (float64, int) {
	magnitude := math.Abs(x)
	result := logGammaPositive(magnitude, ln2, logSqrtTwoPi)
	sign := 1
	if x < 0 {
		floor := math.Floor(magnitude)
		if int64(floor)&1 == 0 {
			sign = -1
		}
		logSin := logarithm(sinPiMagnitude(magnitude, floor, pi), ln2)
		reflection := logarithm(pi, ln2)
		reflection.Sub(reflection, logSin)
		reflection.Sub(reflection, logarithm(precise(magnitude), ln2))
		result.Sub(reflection, result)
	}
	out, _ := result.Float64()
	return out, sign
}

func observe(value float64, sign int) string {
	if math.IsNaN(value) {
		return fmt.Sprintf("NaN %d", sign)
	}
	return fmt.Sprintf("%d %d", binary64.Load(value), sign)
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine Lgamma build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Lgamma architecture %s", runtime.GOARCH)
		return 0
	}
}
