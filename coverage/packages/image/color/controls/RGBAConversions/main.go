package main

import (
	"fmt"
	"image/color"
	"log"
)

type channels [4]uint32

func rgba(c color.Color) channels {
	r, g, b, a := c.RGBA()
	return channels{r, g, b, a}
}

func main() {
	var rescale, div256, div65536, grayAlpha int
	rgba8 := []struct {
		in   color.RGBA
		want channels
	}{
		{color.RGBA{R: 1, G: 128, B: 255}, channels{0x0101, 0x8080, 0xffff, 0}},
		{color.RGBA{R: 2, G: 65, B: 128, A: 128}, channels{0x0202, 0x4141, 0x8080, 0x8080}},
		{color.RGBA{R: 128, G: 1, B: 254, A: 255}, channels{0x8080, 0x0101, 0xfefe, 0xffff}},
	}
	for _, c := range rgba8 {
		if got := rgba(c.in); got != c.want {
			log.Fatalf("RGBA %+v: got %v, want %v", c.in, got, c.want)
		}
		if rgba(color.NRGBA(c.in)) != c.want {
			rescale++
		}
	}
	nrgba := []struct {
		in   color.NRGBA
		want channels
	}{
		{color.NRGBA{R: 1, G: 128, B: 255}, channels{0, 0, 0, 0}},
		{color.NRGBA{R: 1, G: 128, B: 255, A: 255}, channels{257, 32896, 65535, 65535}},
		{color.NRGBA{R: 127, G: 128, B: 129, A: 128}, channels{16383, 16512, 16641, 32896}},
		{color.NRGBA{R: 1, G: 64, B: 254, A: 254}, channels{255, 16383, 65022, 65278}},
		{color.NRGBA{R: 255, G: 170, B: 85, A: 1}, channels{257, 171, 85, 257}},
	}
	for _, c := range nrgba {
		if got := rgba(c.in); got != c.want {
			log.Fatalf("NRGBA %+v: got %v, want %v", c.in, got, c.want)
		}
		a := uint32(c.in.A)
		if (channels{257 * uint32(c.in.R) * a / 256, 257 * uint32(c.in.G) * a / 256, 257 * uint32(c.in.B) * a / 256, 257 * a}) != c.want {
			div256++
		}
	}
	nrgba64 := []struct {
		in   color.NRGBA64
		want channels
	}{
		{color.NRGBA64{R: 1, G: 32768, B: 65535}, channels{0, 0, 0, 0}},
		{color.NRGBA64{R: 1, G: 32768, B: 65535, A: 65535}, channels{1, 32768, 65535, 65535}},
		{color.NRGBA64{R: 32767, G: 32768, B: 32769, A: 32768}, channels{16383, 16384, 16384, 32768}},
		{color.NRGBA64{R: 65535, G: 65534, B: 1, A: 65534}, channels{65534, 65533, 0, 65534}},
		{color.NRGBA64{R: 65535, G: 32768, B: 1, A: 1}, channels{1, 0, 0, 1}},
	}
	for _, c := range nrgba64 {
		if got := rgba(c.in); got != c.want {
			log.Fatalf("NRGBA64 %+v: got %v, want %v", c.in, got, c.want)
		}
		a := uint32(c.in.A)
		if (channels{uint32(c.in.R) * a / 65536, uint32(c.in.G) * a / 65536, uint32(c.in.B) * a / 65536, a}) != c.want {
			div65536++
		}
	}
	for i := range 256 {
		y := uint32(i) * 257
		want := channels{y, y, y, 65535}
		if got := rgba(color.Gray{Y: uint8(i)}); got != want {
			log.Fatalf("Gray %d: got %v, want %v", i, got, want)
		}
		if (channels{y, y, y, y}) != want {
			grayAlpha++
		}
	}
	if rescale != 2 || div256 != 4 || div65536 != 3 || grayAlpha != 255 {
		log.Fatalf("mutation counts rescale=%d divide256=%d divide65536=%d grayAlpha=%d", rescale, div256, div65536, grayAlpha)
	}
	fmt.Printf("RGBA conversions: rescale=%d/%d divide256=%d/%d divide65536=%d/%d grayAlpha=%d/256\n",
		rescale, len(rgba8), div256, len(nrgba), div65536, len(nrgba64), grayAlpha)
}
