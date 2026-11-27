package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	values := []uint16{0, 1, 254, 255, 256, 257, 32767, 32768, 65534, 65535}
	var edges []color.NRGBA64
	for _, r := range values {
		for _, g := range values {
			for _, b := range values {
				for _, a := range values {
					edges = append(edges, color.NRGBA64{R: r, G: g, B: b, A: a})
				}
			}
		}
	}
	check.CompareSamples("godiff-image-color-nrgba64-rgba", edges, sample,
		func(c color.NRGBA64) string {
			return fields(uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A))
		},
		func(c color.NRGBA64) string { return fields(c.RGBA()) },
		func(c color.NRGBA64) string {
			return fields(scale(c.R, c.A), scale(c.G, c.A), scale(c.B, c.A), uint32(c.A))
		})
}

func sample(rng *rand.Rand, _ int) color.NRGBA64 {
	return color.NRGBA64{
		R: uint16(rng.Uint32()),
		G: uint16(rng.Uint32()),
		B: uint16(rng.Uint32()),
		A: uint16(rng.Uint32()),
	}
}

func scale(channel, alpha uint16) uint32 {
	return uint32(uint64(channel) * uint64(alpha) / 65535)
}

func fields(r, g, b, a uint32) string {
	return fmt.Sprintf("%d %d %d %d", r, g, b, a)
}
