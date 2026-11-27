package main

import (
	"fmt"
	"image/color"
	"log"
)

type channels [4]uint32

func compare(c color.Color, want channels) {
	r, g, b, a := c.RGBA()
	if got := (channels{r, g, b, a}); got != want {
		log.Fatalf("%T %+v: got %v, want %v", c, c, got, want)
	}
}

func main() {
	var grayTruncate, cmykIgnoreBlack, ycbcrNarrow, alphaBeforeClamp int
	for i := range 1 << 16 {
		u := uint32(i)
		compare(color.Gray16{Y: uint16(i)}, channels{u, u, u, 65535})
		if uint32(uint8(i)) != u {
			grayTruncate++
		}
	}
	cmyk := []struct {
		in   color.CMYK
		want channels
	}{
		{color.CMYK{C: 0, M: 127, Y: 255, K: 0}, channels{65535, 32896, 0, 65535}},
		{color.CMYK{C: 0, M: 127, Y: 255, K: 255}, channels{0, 0, 0, 65535}},
		{color.CMYK{C: 13, M: 103, Y: 170, K: 194}, channels{14877, 9344, 5225, 65535}},
		{color.CMYK{C: 127, M: 128, Y: 129, K: 127}, channels{16512, 16383, 16254, 65535}},
	}
	for _, c := range cmyk {
		compare(c.in, c.want)
		ignoreBlack := func(v uint8) uint32 { return (255 - uint32(v)) * 257 }
		if (channels{ignoreBlack(c.in.C), ignoreBlack(c.in.M), ignoreBlack(c.in.Y), 65535}) != c.want {
			cmykIgnoreBlack++
		}
	}
	ycbcr := []struct {
		in   color.YCbCr
		want channels
	}{
		{color.YCbCr{Y: 0, Cb: 128, Cr: 128}, channels{0, 0, 0, 65535}},
		{color.YCbCr{Y: 255, Cb: 128, Cr: 128}, channels{65535, 65535, 65535, 65535}},
		{color.YCbCr{Y: 127, Cb: 128, Cr: 129}, channels{32998, 32456, 32639, 65535}},
		{color.YCbCr{Y: 0, Cb: 127, Cr: 128}, channels{0, 88, 0, 65535}},
		{color.YCbCr{Y: 128, Cb: 0, Cr: 255}, channels{65535, 20955, 0, 65535}},
		{color.YCbCr{Y: 255, Cb: 255, Cr: 128}, channels{65535, 54347, 65535, 65535}},
	}
	for _, c := range ycbcr {
		compare(c.in, c.want)
		r, g, b := color.YCbCrToRGB(c.in.Y, c.in.Cb, c.in.Cr)
		if (channels{uint32(r) * 257, uint32(g) * 257, uint32(b) * 257, 65535}) != c.want {
			ycbcrNarrow++
		}
	}
	nycbcra := []struct {
		in   color.NYCbCrA
		want channels
	}{
		{color.NYCbCrA{YCbCr: ycbcr[4].in, A: 0}, channels{0, 0, 0, 0}},
		{color.NYCbCrA{YCbCr: ycbcr[4].in, A: 128}, channels{32896, 10518, 0, 32896}},
		{color.NYCbCrA{YCbCr: ycbcr[4].in, A: 255}, ycbcr[4].want},
		{color.NYCbCrA{YCbCr: ycbcr[5].in, A: 128}, channels{32896, 27280, 32896, 32896}},
		{color.NYCbCrA{YCbCr: ycbcr[2].in, A: 128}, channels{16563, 16291, 16383, 32896}},
	}
	for _, c := range nycbcra {
		compare(c.in, c.want)
		l := int64(c.in.Y) * 65793
		u, v := int64(c.in.Cb)-128, int64(c.in.Cr)-128
		a := int64(c.in.A)
		scaleFirst := func(raw int64) uint32 { return uint32(min(65535, max(0, (raw>>8)*a/255))) }
		if (channels{scaleFirst(l + 91881*v), scaleFirst(l - 22554*u - 46802*v),
			scaleFirst(l + 116130*u), uint32(a) * 257}) != c.want {
			alphaBeforeClamp++
		}
	}
	if grayTruncate != 65280 || cmykIgnoreBlack != 3 || ycbcrNarrow != 4 || alphaBeforeClamp != 2 {
		log.Fatalf("mutation failures: Gray16=%d CMYK=%d YCbCr=%d NYCbCrA=%d",
			grayTruncate, cmykIgnoreBlack, ycbcrNarrow, alphaBeforeClamp)
	}
	fmt.Printf("wide RGBA: all fixtures passed, mutant failures Gray16=%d/65536 CMYK=%d/%d YCbCr=%d/%d NYCbCrA=%d/%d\n",
		grayTruncate, cmykIgnoreBlack, len(cmyk), ycbcrNarrow, len(ycbcr), alphaBeforeClamp, len(nycbcra))
}
