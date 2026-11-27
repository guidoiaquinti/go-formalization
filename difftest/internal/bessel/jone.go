package bessel

import "github.com/guidoiaquinti/go-formalization/difftest/internal/trig"

// JOne evaluates the order-one J function for positive finite x.
func JOne(x float64) float64 {
	if x <= 20 {
		value, _ := seriesOne(precise(x))
		rounded, _ := value.Float64()
		return rounded
	}
	argument := precise(x)
	sine := precise(trig.Sin(x, halfPi))
	cosine := precise(trig.Cos(x, halfPi))
	rotatedSine := precise(0).Add(sine, cosine)
	rotatedSine.Neg(rotatedSine)
	rotatedCosine := precise(0).Sub(sine, cosine)
	p, q := orderOneAsymptotic(argument)
	p.Mul(p, rotatedCosine)
	q.Mul(q, rotatedSine)
	p.Sub(p, q)
	dilation := precise(0).Mul(pi, argument)
	dilation.Sqrt(dilation)
	rounded, _ := p.Quo(p, dilation).Float64()
	return rounded
}
