package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	bounds := []uint8{0, 1, 127, 128, 129, 254, 255}
	var edges []color.YCbCr
	for _, y := range bounds {
		for _, cb := range bounds {
			for _, cr := range bounds {
				edges = append(edges, color.YCbCr{Y: y, Cb: cb, Cr: cr})
			}
		}
	}
	for _, y := range []uint8{0, 128, 255} {
		for cb := range 256 {
			for cr := range 256 {
				edges = append(edges, color.YCbCr{Y: y, Cb: uint8(cb), Cr: uint8(cr)})
			}
		}
	}
	check.CompareSamples("godiff-image-color-ycbcr-rgba", edges, func(rng *rand.Rand, _ int) color.YCbCr {
		return color.YCbCr{Y: uint8(rng.Uint32()), Cb: uint8(rng.Uint32()), Cr: uint8(rng.Uint32())}
	}, func(c color.YCbCr) string {
		return fmt.Sprintf("%d %d %d", c.Y, c.Cb, c.Cr)
	}, func(c color.YCbCr) string {
		return fields(c.RGBA())
	}, expected)
}

func expected(c color.YCbCr) string {
	l := int64(c.Y)
	u, v := int64(c.Cb)-128, int64(c.Cr)-128
	return fields(
		saturate(257*l+floor256(l+91881*v)),
		saturate(257*l+floor256(l-22554*u-46802*v)),
		saturate(257*l+floor256(l+116130*u)),
		65535,
	)
}

func floor256(x int64) int64 {
	q := x / 256
	if x%256 < 0 {
		q--
	}
	return q
}

func saturate(x int64) uint32 { return uint32(min(65535, max(0, x))) }

func fields(r, g, b, a uint32) string { return fmt.Sprintf("%d %d %d %d", r, g, b, a) }
