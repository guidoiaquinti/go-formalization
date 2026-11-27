package bessel

import (
	"math"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/trig"
)

var halfPi = trig.HalfPi()

// JZero evaluates J0 for finite positive x. The central series is used up
// to 20, and larger inputs use independent Hankel amplitudes and phase.
func JZero(x float64) float64 {
	if x <= 20 {
		value, _ := seriesZero(precise(x))
		rounded, _ := value.Float64()
		return rounded
	}
	argument := precise(x)
	sine := precise(trig.Sin(x, halfPi))
	cosine := precise(trig.Cos(x, halfPi))
	sum := precise(0).Add(sine, cosine)
	difference := precise(0).Sub(sine, cosine)
	leading := asymptoticPZero(argument)
	companion := asymptoticQZero(argument)
	leading.Mul(leading, sum)
	companion.Mul(companion, difference)
	leading.Sub(leading, companion)
	dilation := precise(0).Mul(pi, argument)
	dilation.Sqrt(dilation)
	rounded, _ := leading.Quo(leading, dilation).Float64()
	// The independent expression preserves the sign through a rounded zero.
	if rounded == 0 && leading.Sign() < 0 {
		return math.Copysign(0, -1)
	}
	return rounded
}
