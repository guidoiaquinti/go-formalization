package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type triple struct{ y, cb, cr uint8 }

func main() {
	bounds := []uint8{0, 1, 127, 128, 129, 254, 255}
	var edges []triple
	for _, y := range bounds {
		for _, cb := range bounds {
			for _, cr := range bounds {
				edges = append(edges, triple{y, cb, cr})
			}
		}
	}
	for cb := range 256 {
		for cr := range 256 {
			edges = append(edges, triple{128, uint8(cb), uint8(cr)})
		}
	}
	check.CompareSamples("godiff-image-color-ycbcr-to-rgb", edges, func(rng *rand.Rand, i int) triple {
		return triple{uint8(rng.Uint32()), uint8(rng.Uint32()), uint8(rng.Uint32())}
	}, func(c triple) string { return fields(c.y, c.cb, c.cr) },
		func(c triple) string { return fields(color.YCbCrToRGB(c.y, c.cb, c.cr)) }, expected)
}

func expected(c triple) string {
	y := int64(c.y)
	u, v := int64(c.cb)-128, int64(c.cr)-128
	return fields(
		saturate(y+floorDiv(257*y+91881*v, 65536)),
		saturate(y+floorDiv(257*y-22554*u-46802*v, 65536)),
		saturate(y+floorDiv(257*y+116130*u, 65536)),
	)
}

func floorDiv(n, d int64) int64 {
	q := n / d
	if n%d < 0 {
		q--
	}
	return q
}

func saturate(x int64) uint8 { return uint8(min(255, max(0, x))) }

func fields(r, g, b uint8) string { return fmt.Sprintf("%d %d %d", r, g, b) }
