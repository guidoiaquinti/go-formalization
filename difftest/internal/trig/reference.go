// Package trig provides a high-precision trigonometric reference.
package trig

import (
	"math"
	"math/big"
)

const reductionPrecision = 1408
const seriesPrecision = 256

func high(prec uint, value float64) *big.Float {
	return new(big.Float).SetPrec(prec).SetFloat64(value)
}

func HalfPi() *big.Float {
	value := func(x float64) *big.Float { return high(reductionPrecision, x) }
	arctan := func(n float64, terms int64) *big.Float {
		r := value(1).Quo(value(1), value(n))
		square := value(0).Mul(r, r)
		square.Neg(square)
		power, sum := value(0).Set(r), value(0)
		term, divisor := value(0), value(0)
		for i := int64(1); i <= 2*terms-1; i += 2 {
			sum.Add(sum, term.Quo(power, divisor.SetInt64(i)))
			power.Mul(power, square)
		}
		return sum
	}
	// Machin's formula for pi/2.
	left := value(0).Mul(value(8), arctan(5, 330))
	right := value(0).Mul(value(2), arctan(239, 100))
	return left.Sub(left, right)
}

// Sin evaluates a quadrant-reduced Taylor series for sin(x).
func Sin(x float64, halfPi *big.Float) float64 {
	return evaluate(x, halfPi, false)
}

// Cos evaluates a quadrant-reduced Taylor series for cos(x).
func Cos(x float64, halfPi *big.Float) float64 {
	return evaluate(x, halfPi, true)
}

func evaluate(x float64, halfPi *big.Float, cosine bool) float64 {
	value := func(x float64) *big.Float { return high(reductionPrecision, x) }
	absolute := value(math.Abs(x))
	quotient := value(0).Quo(absolute, halfPi)
	quadrant, _ := quotient.Add(quotient, value(0.5)).Int(nil)
	angle := value(0).Sub(absolute, value(0).Mul(value(0).SetInt(quadrant), halfPi))
	q := new(big.Int).And(quadrant, big.NewInt(3)).Uint64()

	arg := high(seriesPrecision, 0).Set(angle)
	negativeSquare := high(seriesPrecision, 0).Mul(arg, arg)
	negativeSquare.Neg(negativeSquare)
	term := high(seriesPrecision, 0).Set(arg)
	power := int64(1)
	if (q&1 != 0) != cosine {
		term.SetInt64(1)
		power = 0
	}
	sum, divisor := high(seriesPrecision, 0).Set(term), high(seriesPrecision, 0)
	for n := power + 2; n <= power+50; n += 2 {
		term.Mul(term, negativeSquare).Quo(term, divisor.SetInt64(n*(n-1)))
		sum.Add(sum, term)
	}
	negative := math.Signbit(x) != (q >= 2)
	if cosine {
		negative = q == 1 || q == 2
	}
	if negative {
		sum.Neg(sum)
	}
	result, _ := sum.Float64()
	return result
}
