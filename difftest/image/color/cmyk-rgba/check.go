package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	// Each ink position sees every ink–black pair, with permuted channels
	// so the comparison also exercises the channel order.
	cases := make([]color.CMYK, 0, 1<<16)
	for black := range 256 {
		for ink := range 256 {
			x := uint8(ink)
			cases = append(cases, color.CMYK{C: x, M: x ^ 0x80, Y: x ^ 0xff, K: uint8(black)})
		}
	}
	check.CompareSamples("godiff-image-color-cmyk-rgba", cases, sample,
		func(c color.CMYK) string {
			return fields(uint32(c.C), uint32(c.M), uint32(c.Y), uint32(c.K))
		},
		func(c color.CMYK) string { return fields(c.RGBA()) },
		func(c color.CMYK) string {
			return fields(scale(c.C, c.K), scale(c.M, c.K), scale(c.Y, c.K), 65535)
		})
}

func sample(rng *rand.Rand, _ int) color.CMYK {
	return color.CMYK{C: uint8(rng.Intn(256)), M: uint8(rng.Intn(256)),
		Y: uint8(rng.Intn(256)), K: uint8(rng.Intn(256))}
}

func scale(ink, black uint8) uint32 {
	return 257 * (255 - uint32(ink)) * (255 - uint32(black)) / 255
}

func fields(r, g, b, a uint32) string {
	return fmt.Sprintf("%d %d %d %d", r, g, b, a)
}
