package main

import (
	"fmt"
	"image/color"
)

type probe struct {
	words [4]uint32
	calls int
}

func (p *probe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	return p.words[0], p.words[1], p.words[2], p.words[3]
}

func main() {
	original := color.CMYK{C: 10, M: 30, Y: 90, K: 200}
	if got := color.CMYKModel.Convert(original); got != original {
		panic("CMYK input changed")
	}
	cases := []struct {
		words [4]uint32
		want  color.CMYK
	}{
		{[4]uint32{0xffff0001, 0x10000002, 0x12340003, 65535}, color.CMYK{K: 255}},
		{[4]uint32{0xffff0001, 0x1000ff02, 0x12340003, 0}, color.CMYK{C: 255, Y: 255}},
		{[4]uint32{0xffff0001, 0x1000ff02, 0x12340003, 65535}, color.CMYK{C: 255, Y: 255}},
	}
	var lowByte, alphaScaled int
	for _, c := range cases {
		p := &probe{words: c.words}
		if got := color.CMYKModel.Convert(p); got != c.want || p.calls != 1 {
			panic(fmt.Sprintf("CMYK conversion=%v want=%v RGBA calls=%d", got, c.want, p.calls))
		}
		r, g, b, a := c.words[0], c.words[1], c.words[2], c.words[3]
		if convert(uint8(r), uint8(g), uint8(b)) != c.want {
			lowByte++
		}
		if convert(scale(r, a), scale(g, a), scale(b, a)) != c.want {
			alphaScaled++
		}
	}
	fmt.Printf("CMYK model: Go passed 4 fixtures, low-byte=%d/3 alpha-scaled=%d/3\n", lowByte, alphaScaled)
	if lowByte != 3 || alphaScaled != 1 {
		panic("unexpected CMYK control counts")
	}
}

func scale(word, alpha uint32) uint8 {
	return uint8(uint32(uint8(word>>8)) * alpha / 65535)
}

func convert(r, g, b uint8) color.CMYK {
	c, m, y, k := color.RGBToCMYK(r, g, b)
	return color.CMYK{C: c, M: m, Y: y, K: k}
}
