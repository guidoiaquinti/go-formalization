package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	channels := []uint8{0, 1, 127, 128, 129, 254, 255}
	var edges []color.NYCbCrA
	for _, y := range channels {
		for _, cb := range channels {
			for _, cr := range channels {
				for a := range 256 {
					edges = append(edges, input(y, cb, cr, uint8(a)))
				}
			}
		}
	}
	check.CompareSamples("godiff-image-color-nycbcra-rgba", edges,
		func(rng *rand.Rand, _ int) color.NYCbCrA {
			w := rng.Uint32()
			return input(uint8(w), uint8(w>>8), uint8(w>>16), uint8(w>>24))
		}, func(c color.NYCbCrA) string {
			return fmt.Sprintf("%d %d %d %d", c.Y, c.Cb, c.Cr, c.A)
		}, func(c color.NYCbCrA) string {
			return fields(c.RGBA())
		}, expected)
}

func input(y, cb, cr, a uint8) color.NYCbCrA {
	return color.NYCbCrA{YCbCr: color.YCbCr{Y: y, Cb: cb, Cr: cr}, A: a}
}

func expected(c color.NYCbCrA) string {
	l := int64(c.Y)
	u, v := int64(c.Cb)-128, int64(c.Cr)-128
	return fields(
		scale(l, 91881*v, c.A),
		scale(l, -22554*u-46802*v, c.A),
		scale(l, 116130*u, c.A),
		uint32(c.A)*257,
	)
}

func scale(l, chroma int64, alpha uint8) uint32 {
	n := l + chroma
	q := n / 256
	if n%256 < 0 {
		q--
	}
	channel := min(65535, max(0, 257*l+q))
	return uint32(channel * int64(alpha) / 255)
}

func fields(r, g, b, a uint32) string {
	return fmt.Sprintf("%d %d %d %d", r, g, b, a)
}
