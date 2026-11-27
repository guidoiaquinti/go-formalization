package bessel

import (
	"math"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/trig"
)

// YZero evaluates the order-zero function for positive finite x. For amd64
// subnormals it evaluates the tiny branch with the effective logarithm
// argument implied by the pinned public Log assembly.
func YZero(x float64, amd64 bool) float64 {
	word := math.Float64bits(x)
	if amd64 && word < 1<<52 {
		effective := precise(math.Float64frombits(word | (1 << 52)))
		effective.SetMantExp(effective, -1)
		_, value := seriesZero(effective)
		rounded, _ := value.Float64()
		return rounded
	}
	if x <= 20 {
		_, value := seriesZero(precise(x))
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
	leading.Mul(leading, difference)
	companion.Mul(companion, sum)
	leading.Add(leading, companion)
	dilation := precise(0).Mul(pi, argument)
	dilation.Sqrt(dilation)
	rounded, _ := leading.Quo(leading, dilation).Float64()
	return rounded
}
