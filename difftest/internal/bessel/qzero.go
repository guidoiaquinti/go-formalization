package bessel

import "math/big"

// QZero evaluates the order-zero complementary Bessel amplitude for finite x >= 2.
func QZero(x float64) float64 {
	argument := precise(x)
	var amplitude *big.Float
	if x <= 20 {
		j, y := seriesZero(argument)
		theta := precise(0).SetMantExp(pi, -2)
		theta.Sub(argument, theta)
		sine, cosine := sineCosine(theta)
		dilation := precise(0).Mul(pi, argument)
		dilation.Sqrt(dilation.SetMantExp(dilation, -1))
		amplitude = precise(0).Mul(y, cosine)
		amplitude.Sub(amplitude, precise(0).Mul(j, sine))
		amplitude.Mul(amplitude, dilation)
	} else {
		amplitude = asymptoticQZero(argument)
	}
	result, _ := amplitude.Float64()
	return result
}

func asymptoticQZero(x *big.Float) *big.Float {
	sum, term := precise(0), precise(1)
	denominator, numerator, previous, magnitude := precise(0), precise(0), precise(1), precise(0)
	for n := int64(1); n <= 120; n++ {
		// The Hankel expansion contributes alternating odd-index terms to Q.
		term.Mul(term, numerator.SetInt64(-(2*n-1)*(2*n-1)))
		denominator.Mul(x, denominator.SetInt64(8*n))
		term.Quo(term, denominator)
		magnitude.Abs(term)
		if magnitude.Cmp(previous) > 0 {
			break
		}
		previous.Set(magnitude)
		switch n % 4 {
		case 1:
			sum.Add(sum, term)
		case 3:
			sum.Sub(sum, term)
		}
		if n > 5 && term.MantExp(nil) < -230 {
			break
		}
	}
	return sum
}
