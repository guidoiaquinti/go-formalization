package bessel

import "math/big"

// POne evaluates the order-one Bessel amplitude for finite x >= 2.
func POne(x float64) float64 {
	argument := precise(x)
	if x > 20 {
		p, _ := orderOneAsymptotic(argument)
		rounded, _ := p.Float64()
		return rounded
	}
	j, y := seriesOne(argument)
	theta := precise(0).Mul(pi, precise(0.75))
	theta.Sub(argument, theta)
	sine, cosine := sineCosine(theta)
	dilation := precise(0).Mul(pi, argument)
	dilation.Sqrt(dilation.SetMantExp(dilation, -1))
	p := precise(0).Mul(j, cosine)
	p.Add(p, precise(0).Mul(y, sine)).Mul(p, dilation)
	rounded, _ := p.Float64()
	return rounded
}

func seriesOne(x *big.Float) (*big.Float, *big.Float) {
	square := precise(0).Mul(x, x)
	square.Neg(square.SetMantExp(square, -2))
	termZero, j0 := precise(1), precise(1)
	termOne := precise(0).SetMantExp(x, -1)
	j1 := precise(0).Set(termOne)
	harmonic, derivative := precise(0), precise(0)
	work, divisor, index := precise(0), precise(0), precise(0)
	for k := int64(1); k <= 145; k++ {
		termZero.Mul(termZero, square).Quo(termZero, divisor.SetInt64(k*k))
		j0.Add(j0, termZero)
		termOne.Mul(termOne, square).Quo(termOne, divisor.SetInt64(k*(k+1)))
		j1.Add(j1, termOne)
		harmonic.Add(harmonic, work.SetInt64(1).Quo(work, index.SetInt64(k)))
		work.Mul(harmonic, termZero).Mul(work, index.SetInt64(2*k))
		derivative.Add(derivative, work)
		if k > 20 && termZero.MantExp(nil) < -225 {
			break
		}
	}
	// Y1 = -Y0'. Differentiate the order-zero logarithm and harmonic series.
	y1 := logarithm(precise(0).SetMantExp(x, -1), ln2)
	y1.Add(y1, euler).Mul(y1, j1)
	j0.Sub(j0, derivative)
	y1.Sub(y1, j0.Quo(j0, x))
	return j1, y1.Mul(y1, precise(2).Quo(precise(2), pi))
}

func orderOneAsymptotic(x *big.Float) (*big.Float, *big.Float) {
	p, q, term := precise(1), precise(0), precise(1)
	denominator, numerator, previous, magnitude := precise(0), precise(0), precise(1), precise(0)
	for n := int64(1); n <= 120; n++ {
		factor := 4 - (2*n-1)*(2*n-1)
		term.Mul(term, numerator.SetInt64(factor))
		denominator.Mul(x, denominator.SetInt64(8*n))
		term.Quo(term, denominator)
		magnitude.Abs(term)
		if magnitude.Cmp(previous) > 0 {
			break
		}
		previous.Set(magnitude)
		switch n % 4 {
		case 0:
			p.Add(p, term)
		case 1:
			q.Add(q, term)
		case 2:
			p.Sub(p, term)
		case 3:
			q.Sub(q, term)
		}
		if n > 5 && term.MantExp(nil) < -230 {
			break
		}
	}
	return p, q
}
