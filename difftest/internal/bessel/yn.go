package bessel

import (
	"math"
	"math/big"
)

// YOrder evaluates the integer-order Y function for nonnegative orders
// through 64 and positive finite arguments. Seeds near the pole stay
// unrounded until the final recurrence has been evaluated. The second
// result measures the unrounded distance from the binary64 finite limit.
func YOrder(n int, x float64, amd64 bool) (value, relativeLimitGap float64) {
	if n == 0 {
		return YZero(x, amd64), 1
	}
	argument := precise(x)
	var a, b *big.Float
	if x <= 20 {
		_, a = seriesZero(argument)
		_, b = seriesOne(argument)
	} else {
		a, b = precise(YZero(x, amd64)), precise(YOne(x))
	}
	for i := 1; i < n; i++ {
		factor := precise(float64(2 * i))
		factor.Quo(factor, argument).Mul(factor, b).Sub(factor, a)
		a, b = b, factor
	}
	value, _ = b.Float64()
	magnitude := precise(0).Abs(b)
	if magnitude.Cmp(precise(math.MaxFloat64/2)) < 0 {
		return value, 1
	}
	limit := precise(math.MaxFloat64)
	distance := magnitude.Sub(magnitude, limit)
	distance.Abs(distance).Quo(distance, limit)
	relativeLimitGap, _ = distance.Float64()
	return value, relativeLimitGap
}
