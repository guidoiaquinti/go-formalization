package main

import (
	"fmt"
	"image/color"
)

type probe struct {
	red, alpha uint32
	calls      int
}

func (p *probe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	return p.red, 0, 0, p.alpha
}

func main() {
	ycbcr := color.YCbCr{Y: 21, Cb: 43, Cr: 65}
	original := color.NYCbCrA{YCbCr: ycbcr, A: 7}
	if color.NYCbCrAModel.Convert(original) != original ||
		color.NYCbCrAModel.Convert(ycbcr) != (color.NYCbCrA{YCbCr: ycbcr, A: 255}) {
		panic("incorrect concrete value path")
	}
	wrapped := &probe{red: 65538, alpha: 300}
	transparent := &probe{red: 0xff00}
	black := color.NYCbCrA{YCbCr: color.YCbCr{Cb: 128, Cr: 128}, A: 1}
	red := color.NYCbCrA{YCbCr: color.YCbCr{Y: 76, Cb: 85, Cr: 255}}
	if color.NYCbCrAModel.Convert(wrapped) != black ||
		color.NYCbCrAModel.Convert(transparent) != red ||
		wrapped.calls != 1 || transparent.calls != 1 {
		panic("incorrect callback conversion")
	}
	wideRed := uint8((uint64(wrapped.red) * 65535 / uint64(wrapped.alpha)) >> 8)
	y, cb, cr := color.RGBToYCbCr(wideRed, 0, 0)
	wide := color.NYCbCrA{YCbCr: color.YCbCr{Y: y, Cb: cb, Cr: cr}, A: 1}
	zeroed := color.NYCbCrA{YCbCr: black.YCbCr}
	if wide == black || zeroed == red {
		panic("expected both mutants to disagree")
	}
	fmt.Println("NYCbCrA model: Go passed 4 fixtures; widened and zero-alpha RGB mutants rejected")
}
