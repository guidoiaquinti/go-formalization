package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type rgb struct{ r, g, b uint8 }

func main() {
	// Sweep every pair with the third channel at its lower or upper bound.
	edges := make([]rgb, 0, 6*256*256)
	for _, fixed := range []uint8{0, 255} {
		for x := range 256 {
			for y := range 256 {
				u, v := uint8(x), uint8(y)
				edges = append(edges, rgb{fixed, u, v}, rgb{u, fixed, v}, rgb{u, v, fixed})
			}
		}
	}
	check.CompareSamples("godiff-image-color-rgb-to-ycbcr", edges,
		func(rng *rand.Rand, _ int) rgb {
			return rgb{uint8(rng.Uint32()), uint8(rng.Uint32()), uint8(rng.Uint32())}
		},
		func(c rgb) string { return fmt.Sprintf("%d %d %d", c.r, c.g, c.b) },
		func(c rgb) string {
			y, cb, cr := color.RGBToYCbCr(c.r, c.g, c.b)
			return fmt.Sprintf("%d %d %d", y, cb, cr)
		}, expected)
}

func expected(c rgb) string {
	r, g, b := int64(c.r), int64(c.g), int64(c.b)
	round := func(delta int64) int64 {
		// Round the signed fraction by shifting the numerator to a
		// nonnegative range before dividing, then undoing the shift.
		return (delta+257*(1<<15))/(1<<16) - 128
	}
	y := g + round(19595*(r-g)+7471*(b-g))
	cb := min(255, max(0, 128+round(32768*(b-g)-11056*(r-g))))
	cr := min(255, max(0, 128+round(32768*(r-g)-5328*(b-g))))
	return fmt.Sprintf("%d %d %d", y, cb, cr)
}
