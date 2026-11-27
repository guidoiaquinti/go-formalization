package math_test

import (
	"math"
	"testing"
)

func TestGoLeanGenericFloorCeilControl(t *testing.T) {
	z := math.Copysign(0, -1)
	cases := []struct{ input, floor, ceil float64 }{
		{0, 0, 0}, {z, z, z}, {0x1p-1074, 0, 1}, {-0x1p-1074, -1, z},
		{1.75, 1, 2}, {-1.75, -2, -1}, {0x1.fffffffffffffp+51, 0x1.ffffffffffffep+51, 0x1p52},
		{-0x1.fffffffffffffp+51, -0x1p52, -0x1.ffffffffffffep+51},
		{math.Inf(1), math.Inf(1), math.Inf(1)},
		{math.Inf(-1), math.Inf(-1), math.Inf(-1)},
		{math.Float64frombits(0x7ff0000000000001), math.NaN(), math.NaN()},
	}
	same := func(actual, expected float64) bool {
		return math.Float64bits(actual) == math.Float64bits(expected) ||
			(math.IsNaN(actual) && math.IsNaN(expected))
	}
	var truncateFloor, truncateCeil, eraseFloorZero, eraseCeilZero int
	for _, c := range cases {
		floor, ceil := math.GoLeanFloor(c.input), math.GoLeanCeil(c.input)
		if !same(floor, c.floor) || !same(ceil, c.ceil) {
			t.Fatalf("input %016x: floor=%016x ceil=%016x", math.Float64bits(c.input), math.Float64bits(floor), math.Float64bits(ceil))
		}
		if !same(math.Trunc(c.input), c.floor) {
			truncateFloor++
		}
		if !same(math.Trunc(c.input), c.ceil) {
			truncateCeil++
		}
		if floor == 0 && !same(0, c.floor) {
			eraseFloorZero++
		}
		if ceil == 0 && !same(0, c.ceil) {
			eraseCeilZero++
		}
	}
	if truncateFloor != 3 || truncateCeil != 3 || eraseFloorZero != 1 || eraseCeilZero != 2 {
		t.Fatalf("mutation counts truncation=%d/%d zero=%d/%d", truncateFloor, truncateCeil, eraseFloorZero, eraseCeilZero)
	}
	t.Logf("%d cases pass, truncate floor/ceil=%d/%d, erase signed zero=%d/%d",
		len(cases), truncateFloor, truncateCeil, eraseFloorZero, eraseCeilZero)
}
