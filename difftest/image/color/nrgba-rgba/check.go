package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	// For each alpha, each channel sees all 256 values. The three
	// permutations also exercise the channel order on unequal RGB inputs.
	cases := make([]color.NRGBA, 0, 1<<16)
	for alpha := range 256 {
		for channel := range 256 {
			x := uint8(channel)
			cases = append(cases, color.NRGBA{R: x, G: x ^ 0x80, B: x ^ 0xff, A: uint8(alpha)})
		}
	}
	check.CompareSamples("godiff-image-color-nrgba-rgba", cases, sample,
		func(c color.NRGBA) string {
			return fields(uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A))
		},
		func(c color.NRGBA) string { return fields(c.RGBA()) },
		func(c color.NRGBA) string {
			a := uint32(c.A)
			return fields(scale(c.R, a), scale(c.G, a), scale(c.B, a), 257*a)
		})
}

func sample(rng *rand.Rand, _ int) color.NRGBA {
	return color.NRGBA{R: uint8(rng.Intn(256)), G: uint8(rng.Intn(256)),
		B: uint8(rng.Intn(256)), A: uint8(rng.Intn(256))}
}

func scale(channel uint8, alpha uint32) uint32 {
	product := uint32(channel) * alpha
	return product + 2*product/255
}

func fields(r, g, b, a uint32) string {
	return fmt.Sprintf("%d %d %d %d", r, g, b, a)
}
