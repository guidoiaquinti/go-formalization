package bessel

// QOne evaluates the complementary order-one Bessel amplitude for finite x >= 2.
func QOne(x float64) float64 {
	argument := precise(x)
	if x > 20 {
		_, q := orderOneAsymptotic(argument)
		rounded, _ := q.Float64()
		return rounded
	}
	j, y := seriesOne(argument)
	theta := precise(0).Mul(pi, precise(0.75))
	theta.Sub(argument, theta)
	sine, cosine := sineCosine(theta)
	dilation := precise(0).Mul(pi, argument)
	dilation.Sqrt(dilation.SetMantExp(dilation, -1))
	q := precise(0).Mul(y, cosine)
	q.Sub(q, precise(0).Mul(j, sine)).Mul(q, dilation)
	rounded, _ := q.Float64()
	return rounded
}
