package main

import (
	"fmt"
	"image/color"
)

type probe struct {
	words [4]uint32
	calls int
	fault *int
}

func (p *probe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	if p.fault != nil {
		panic(p.fault)
	}
	return p.words[0], p.words[1], p.words[2], p.words[3]
}

func main() {
	original := color.YCbCr{Y: 1, Cb: 2, Cr: 3}
	if got := color.YCbCrModel.Convert(original); got != original {
		panic("YCbCr input changed")
	}
	want := color.YCbCr{Y: 76, Cb: 85, Cr: 255}
	var lowByte, clampWord int
	for _, alpha := range []uint32{0, 65535, ^uint32(0)} {
		input := &probe{words: [4]uint32{0x1234ff01, 0xffff0002, 0x10000003, alpha}}
		if got := color.YCbCrModel.Convert(input); got != want || input.calls != 1 {
			panic(fmt.Sprintf("YCbCr conversion=%v callback count=%d", got, input.calls))
		}
		if convert(uint8(input.words[0]), uint8(input.words[1]), uint8(input.words[2])) != want {
			lowByte++
		}
		if convert(255, 255, 255) != want {
			clampWord++
		}
	}
	failure := 42
	input := &probe{fault: &failure}
	if got := fault(input); got != &failure || input.calls != 1 {
		panic("callback panic or call count changed")
	}
	if fault(nil) == nil {
		panic("nil Color returned normally")
	}
	fmt.Printf("YCbCr model: Go passed 6 fixtures, low-byte=%d/3 clamp-word=%d/3\n", lowByte, clampWord)
	if lowByte != 3 || clampWord != 3 {
		panic("unexpected YCbCr control counts")
	}
}

func convert(r, g, b uint8) color.YCbCr {
	y, cb, cr := color.RGBToYCbCr(r, g, b)
	return color.YCbCr{Y: y, Cb: cb, Cr: cr}
}

func fault(c color.Color) (result any) {
	defer func() { result = recover() }()
	color.YCbCrModel.Convert(c)
	return nil
}
