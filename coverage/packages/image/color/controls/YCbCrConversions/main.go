package main

import (
	"fmt"
	"image/color"
	"log"
)

type triple [3]uint8

func main() {
	rgb := []struct{ in, want triple }{
		{triple{0, 0, 0}, triple{0, 128, 128}},
		{triple{255, 255, 255}, triple{255, 128, 128}},
		{triple{255, 0, 0}, triple{76, 85, 255}},
		{triple{0, 255, 0}, triple{150, 44, 21}},
		{triple{0, 0, 255}, triple{29, 255, 107}},
		{triple{255, 255, 0}, triple{226, 1, 149}},
		{triple{0, 3, 0}, triple{2, 127, 127}},
		{triple{1, 2, 3}, triple{2, 129, 127}},
	}
	saturate := func(input int64) uint8 { return uint8(min(255, max(0, input>>16))) }
	var rgbFloor, rgbWrap int
	for _, c := range rgb {
		y, cb, cr := color.RGBToYCbCr(c.in[0], c.in[1], c.in[2])
		if got := (triple{y, cb, cr}); got != c.want {
			log.Fatalf("RGBToYCbCr %v got %v, want %v", c.in, got, c.want)
		}
		r, g, b := int64(c.in[0]), int64(c.in[1]), int64(c.in[2])
		l := 19595*r + 38470*g + 7471*b
		u := -11056*r - 21712*g + 32768*b + 128<<16
		v := 32768*r - 27440*g - 5328*b + 128<<16
		if (triple{saturate(l), saturate(u), saturate(v)}) != c.want {
			rgbFloor++
		}
		if (triple{uint8((l + 32768) >> 16), uint8((u + 32768) >> 16), uint8((v + 32768) >> 16)}) != c.want {
			rgbWrap++
		}
	}
	ycbcr := []struct{ in, want triple }{
		{triple{0, 128, 128}, triple{0, 0, 0}},
		{triple{255, 128, 128}, triple{255, 255, 255}},
		{triple{128, 0, 255}, triple{255, 81, 0}},
		{triple{10, 129, 127}, triple{8, 10, 11}},
		{triple{128, 129, 127}, triple{127, 128, 130}},
		{triple{255, 0, 255}, triple{255, 209, 29}},
		{triple{1, 255, 128}, triple{1, 0, 226}},
		{triple{250, 127, 128}, triple{250, 251, 249}},
		{triple{0, 128, 129}, triple{1, 0, 0}},
	}
	var rgbConstant, ycbcrWrap int
	for _, c := range ycbcr {
		r, g, b := color.YCbCrToRGB(c.in[0], c.in[1], c.in[2])
		if got := (triple{r, g, b}); got != c.want {
			log.Fatalf("YCbCrToRGB %v got %v, want %v", c.in, got, c.want)
		}
		l, u, v := int64(c.in[0]), int64(c.in[1])-128, int64(c.in[2])-128
		dr, dg, db := 91881*v, -22554*u-46802*v, 116130*u
		constant := l*65536 + 32768
		if (triple{saturate(constant + dr), saturate(constant + dg), saturate(constant + db)}) != c.want {
			rgbConstant++
		}
		adjusted := l * 65793
		if (triple{uint8((adjusted + dr) >> 16), uint8((adjusted + dg) >> 16), uint8((adjusted + db) >> 16)}) != c.want {
			ycbcrWrap++
		}
	}
	if rgbFloor != 5 || rgbWrap != 2 || rgbConstant != 3 || ycbcrWrap != 4 {
		log.Fatalf("mutation failures: RGB floor=%d, RGB wrap=%d, YCbCr constant=%d, YCbCr wrap=%d",
			rgbFloor, rgbWrap, rgbConstant, ycbcrWrap)
	}
	fmt.Printf("YCbCr conversions: %d+%d fixtures pass, mutant failures RGB floor=%d, RGB wrap=%d, YCbCr constant=%d, YCbCr wrap=%d\n",
		len(rgb), len(ycbcr), rgbFloor, rgbWrap, rgbConstant, ycbcrWrap)
}
