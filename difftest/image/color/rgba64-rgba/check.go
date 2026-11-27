package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	values := []uint16{0, 1, 32767, 32768, 65534, 65535}
	var edges []color.RGBA64
	for _, r := range values {
		for _, g := range values {
			for _, b := range values {
				for _, a := range values {
					edges = append(edges, color.RGBA64{R: r, G: g, B: b, A: a})
				}
			}
		}
	}
	check.CompareSamples("godiff-image-color-rgba64-rgba", edges, sample,
		func(c color.RGBA64) string {
			return fields(uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A))
		},
		func(c color.RGBA64) string { return fields(c.RGBA()) },
		func(c color.RGBA64) string {
			return fields(uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A))
		})
}

func sample(rng *rand.Rand, index int) color.RGBA64 {
	a := uint16(rng.Intn(65536))
	limit := 65536
	if index%2 == 0 {
		limit = int(a) + 1
	}
	return color.RGBA64{
		R: uint16(rng.Intn(limit)),
		G: uint16(rng.Intn(limit)),
		B: uint16(rng.Intn(limit)),
		A: a,
	}
}

func fields(r, g, b, a uint32) string {
	return fmt.Sprintf("%d %d %d %d", r, g, b, a)
}
