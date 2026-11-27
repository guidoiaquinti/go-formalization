package math_test

import (
	"math"
	"testing"
)

func TestGoLeanNormalizeControl(t *testing.T) {
	const sign uint64 = 1 << 63
	cases := []struct {
		input, word          uint64
		adjustment, exponent int
	}{
		{1, 1 << 52, -52, -1074},
		{sign | 3, sign | 0x0028000000000000, -52, -1073},
		{1 << 51, 52 << 52, -52, -1023},
		{(1 << 52) - 1, 0x034ffffffffffffe, -52, -1023},
		{1 << 52, 1 << 52, 0, -1022},
		{0x3ff0000000000000, 0x3ff0000000000000, 0, 0},
		{0xffefffffffffffff, 0xffefffffffffffff, 0, 1023},
	}
	var skipScaling, wrongScale, noAdjust, exponentOnly int
	for _, c := range cases {
		x := math.Float64frombits(c.input)
		normal, adjustment := math.GoLeanNormalize(x)
		exponent := math.GoLeanIlogb(x)
		bits := math.Float64bits(normal)
		if bits != c.word || adjustment != c.adjustment || exponent != c.exponent {
			t.Fatalf("%#016x: got %#016x %d %d, want %#016x %d %d",
				c.input, bits, adjustment, exponent, c.word, c.adjustment, c.exponent)
		}
		if c.input != c.word || c.adjustment != 0 {
			skipScaling++
		}
		mutant, adjust := x, 0
		if math.Abs(x) < 0x1p-1022 {
			mutant, adjust = x*0x1p51, -51
		}
		if math.Float64bits(mutant) != c.word || adjust != c.adjustment {
			wrongScale++
		}
		if int(bits>>52&2047)-1023 != c.exponent {
			noAdjust++
		}
		if int(c.input>>52&2047)-1023 != c.exponent {
			exponentOnly++
		}
	}
	if skipScaling != 4 || wrongScale != 4 || noAdjust != 4 || exponentOnly != 2 {
		t.Fatalf("mutation failures skip=%d wrong scale=%d missing adjustment=%d exponent field only=%d",
			skipScaling, wrongScale, noAdjust, exponentOnly)
	}
	t.Logf("%d inputs pass, mutant failures skip=%d wrong scale=%d missing adjustment=%d exponent field only=%d",
		len(cases), skipScaling, wrongScale, noAdjust, exponentOnly)
}
