package bessel

import "github.com/guidoiaquinti/go-formalization/difftest/internal/trig"

// YOne evaluates the order-one Y function for positive finite x.
func YOne(x float64) float64 {
	argument := precise(x)
	if x <= 20 {
		_, value := seriesOne(argument)
		rounded, _ := value.Float64()
		return rounded
	}
	sine := precise(trig.Sin(x, halfPi))
	cosine := precise(trig.Cos(x, halfPi))
	rotatedSine := precise(0).Add(sine, cosine)
	rotatedSine.Neg(rotatedSine)
	rotatedCosine := precise(0).Sub(sine, cosine)
	p, q := orderOneAsymptotic(argument)
	p.Mul(p, rotatedSine)
	q.Mul(q, rotatedCosine)
	p.Add(p, q)
	dilation := precise(0).Mul(pi, argument)
	dilation.Sqrt(dilation)
	rounded, _ := p.Quo(p, dilation).Float64()
	return rounded
}
