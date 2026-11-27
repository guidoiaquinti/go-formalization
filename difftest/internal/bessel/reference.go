// Package bessel contains high-precision references for the Bessel comparisons.
package bessel

import "math/big"

var (
	pi    = decimal("3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117068")
	euler = decimal("0.5772156649015328606065120900824024310421593359399235988057672348848677267776646709369470632917467495")
	ln2   = atanhLog(precise(0).Quo(precise(1), precise(3)))
)

func decimal(value string) *big.Float {
	number, _, err := big.ParseFloat(value, 10, 256, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	return number
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

// PZero evaluates the order-zero Bessel amplitude for finite x >= 2.
func PZero(x float64) float64 {
	argument := precise(x)
	var amplitude *big.Float
	if x <= 20 {
		// Invert the J0/Y0 phase rotation, using independent power series.
		j, y := seriesZero(argument)
		theta := precise(0).SetMantExp(pi, -2)
		theta.Sub(argument, theta)
		sine, cosine := sineCosine(theta)
		dilation := precise(0).Mul(pi, argument)
		dilation.Sqrt(dilation.SetMantExp(dilation, -1))
		amplitude = precise(0).Mul(j, cosine)
		amplitude.Add(amplitude, precise(0).Mul(y, sine))
		amplitude.Mul(amplitude, dilation)
	} else {
		amplitude = asymptoticPZero(argument)
	}
	result, _ := amplitude.Float64()
	return result
}

func seriesZero(x *big.Float) (*big.Float, *big.Float) {
	z := precise(0).Mul(x, x)
	z.Neg(z.SetMantExp(z, -2))
	term, j, harmonic, weighted := precise(1), precise(1), precise(0), precise(0)
	numerator, denominator := precise(0), precise(0)
	for k := int64(1); k <= 145; k++ {
		term.Mul(term, z)
		term.Quo(term, denominator.SetInt64(k*k))
		j.Add(j, term)
		harmonic.Add(harmonic, numerator.SetInt64(1).Quo(numerator, denominator.SetInt64(k)))
		weighted.Add(weighted, numerator.Mul(harmonic, term))
		if k > 20 && term.MantExp(nil) < -225 {
			break
		}
	}
	y := logarithm(precise(0).SetMantExp(x, -1), ln2)
	y.Add(y, euler).Mul(y, j).Sub(y, weighted)
	y.Mul(y, precise(2).Quo(precise(2), pi))
	return j, y
}

func sineCosine(x *big.Float) (*big.Float, *big.Float) {
	factor := precise(0).Mul(x, x)
	factor.Neg(factor)
	sinTerm, cosTerm := precise(0).Set(x), precise(1)
	sin, cos, divisor := precise(0).Set(x), precise(1), precise(0)
	for n := int64(1); n <= 130; n++ {
		sinTerm.Mul(sinTerm, factor).Quo(sinTerm, divisor.SetInt64(2*n*(2*n+1)))
		cosTerm.Mul(cosTerm, factor).Quo(cosTerm, divisor.SetInt64(2*n*(2*n-1)))
		sin.Add(sin, sinTerm)
		cos.Add(cos, cosTerm)
		if n > 20 && sinTerm.MantExp(nil) < -225 && cosTerm.MantExp(nil) < -225 {
			break
		}
	}
	return sin, cos
}

func asymptoticPZero(x *big.Float) *big.Float {
	sum, term := precise(1), precise(1)
	denominator, numerator, previous, magnitude := precise(0), precise(0), precise(1), precise(0)
	for n := int64(1); n <= 120; n++ {
		// The Hankel expansion contributes alternating even-index terms to P.
		term.Mul(term, numerator.SetInt64(-(2*n-1)*(2*n-1)))
		denominator.Mul(x, denominator.SetInt64(8*n))
		term.Quo(term, denominator)
		magnitude.Abs(term)
		if magnitude.Cmp(previous) > 0 {
			break
		}
		previous.Set(magnitude)
		switch n % 4 {
		case 0:
			sum.Add(sum, term)
		case 2:
			sum.Sub(sum, term)
		}
		if n > 5 && term.MantExp(nil) < -230 {
			break
		}
	}
	return sum
}
