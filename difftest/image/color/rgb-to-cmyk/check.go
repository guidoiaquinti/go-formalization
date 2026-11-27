package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type rgb struct{ r, g, b uint8 }

func main() {
	// For each possible channel and white value, place the channel in each
	// position and let an adjacent channel supply white.
	edges := make([]rgb, 0, 3*256*257/2)
	for white := range 256 {
		for channel := range white + 1 {
			v, w := uint8(channel), uint8(white)
			edges = append(edges, rgb{v, w, 0}, rgb{0, v, w}, rgb{w, 0, v})
		}
	}
	check.CompareSamples("godiff-image-color-rgb-to-cmyk", edges,
		func(rng *rand.Rand, _ int) rgb {
			return rgb{uint8(rng.Uint32()), uint8(rng.Uint32()), uint8(rng.Uint32())}
		}, func(c rgb) string {
			return fmt.Sprintf("%d %d %d", c.r, c.g, c.b)
		}, func(c rgb) string {
			cyan, magenta, yellow, black := color.RGBToCMYK(c.r, c.g, c.b)
			return fmt.Sprintf("%d %d %d %d", cyan, magenta, yellow, black)
		}, expected)
}

func expected(c rgb) string {
	white := max(int(c.r), int(c.g), int(c.b))
	if white == 0 {
		return "0 0 0 255"
	}
	// Select the largest ink whose weighted sum does not exceed full scale.
	ink := func(channel uint8) int {
		return 255 - (255*int(channel)+white-1)/white
	}
	return fmt.Sprintf("%d %d %d %d", ink(c.r), ink(c.g), ink(c.b), 255-white)
}
