package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	values := []uint8{0, 1, 127, 128, 254, 255}
	var edges []color.RGBA
	for _, r := range values {
		for _, g := range values {
			for _, b := range values {
				for _, a := range values {
					edges = append(edges, color.RGBA{R: r, G: g, B: b, A: a})
				}
			}
		}
	}
	check.CompareSamples("godiff-image-color-rgba-rgba", edges, sample,
		func(c color.RGBA) string {
			return fields(uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A))
		},
		func(c color.RGBA) string { return fields(c.RGBA()) },
		func(c color.RGBA) string {
			return fields(257*uint32(c.R), 257*uint32(c.G), 257*uint32(c.B), 257*uint32(c.A))
		})
}

func sample(rng *rand.Rand, index int) color.RGBA {
	a := uint8(rng.Intn(256))
	limit := 256
	if index%2 == 0 {
		limit = int(a) + 1
	}
	return color.RGBA{
		R: uint8(rng.Intn(limit)),
		G: uint8(rng.Intn(limit)),
		B: uint8(rng.Intn(limit)),
		A: a,
	}
}

func fields(r, g, b, a uint32) string {
	return fmt.Sprintf("%d %d %d %d", r, g, b, a)
}
