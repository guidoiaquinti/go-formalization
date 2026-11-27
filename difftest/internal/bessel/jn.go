package bessel

// JOrder evaluates the integer-order J function for nonnegative orders
// through 64 and positive finite arguments. A direct Taylor series avoids
// the unstable recurrence when the argument is below the order.
func JOrder(n int, x float64) float64 {
	if x > 64 {
		a, b := precise(JZero(x)), precise(JOne(x))
		if n == 0 {
			value, _ := a.Float64()
			return value
		}
		argument := precise(x)
		for i := 1; i < n; i++ {
			factor := precise(float64(2 * i))
			factor.Quo(factor, argument).Mul(factor, b).Sub(factor, a)
			a, b = b, factor
		}
		value, _ := b.Float64()
		return value
	}
	half := precise(x)
	half.SetMantExp(half, -1)
	term, divisor := precise(1), precise(0)
	for i := 1; i <= n; i++ {
		term.Mul(term, half).Quo(term, divisor.SetInt64(int64(i)))
	}
	sum := precise(0).Set(term)
	square := precise(0).Mul(half, half)
	square.Neg(square)
	cutoff := term.MantExp(nil) - 250
	for k := 1; k <= 300; k++ {
		term.Mul(term, square).Quo(term, divisor.SetInt64(int64(k*(n+k))))
		sum.Add(sum, term)
		if k > 64 && term.MantExp(nil) < cutoff {
			break
		}
	}
	rounded, _ := sum.Float64()
	return rounded
}
