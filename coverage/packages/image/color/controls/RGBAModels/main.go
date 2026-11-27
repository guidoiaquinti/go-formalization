package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
)

type callback struct {
	rgba    [4]uint32
	calls   int
	payload *int
}

func (c *callback) RGBA() (uint32, uint32, uint32, uint32) {
	c.calls++
	if c.payload != nil {
		panic(c.payload)
	}
	return c.rgba[0], c.rgba[1], c.rgba[2], c.rgba[3]
}

func capture(m color.Model, c color.Color) (recovered any) {
	defer func() { recovered = recover() }()
	m.Convert(c)
	return nil
}

func main() {
	cases := []struct {
		words  [4]uint32
		want8  color.RGBA
		want16 color.RGBA64
	}{
		{[4]uint32{}, color.RGBA{}, color.RGBA64{}},
		{[4]uint32{0x12345678, 0xabcdef01, 0xfedcba98, 0xf0000001},
			color.RGBA{R: 0x56, G: 0xef, B: 0xba, A: 0}, color.RGBA64{R: 0x5678, G: 0xef01, B: 0xba98, A: 1}},
		{[4]uint32{0x00000100, 0xff000080, 0xffff0000, 0x11112222},
			color.RGBA{R: 1, G: 0, B: 0, A: 0x22}, color.RGBA64{R: 0x100, G: 0x80, B: 0, A: 0x2222}},
		{[4]uint32{0x000000ff, 0x0000ffff, 0x00010000, 0xffffffff},
			color.RGBA{R: 0, G: 255, B: 0, A: 255}, color.RGBA64{R: 255, G: 65535, B: 0, A: 65535}},
	}
	var lowByte, highHalf int
	for _, c := range cases {
		p := &callback{rgba: c.words}
		if got := color.RGBAModel.Convert(p); got != c.want8 || p.calls != 1 {
			log.Fatalf("rgbaModel %x got %v, calls %d", c.words, got, p.calls)
		}
		p.calls = 0
		if got := color.RGBA64Model.Convert(p); got != c.want16 || p.calls != 1 {
			log.Fatalf("rgba64Model %x got %v, calls %d", c.words, got, p.calls)
		}
		r, g, b, a := c.words[0], c.words[1], c.words[2], c.words[3]
		if (color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}) != c.want8 {
			lowByte++
		}
		if (color.RGBA64{R: uint16(r >> 16), G: uint16(g >> 16), B: uint16(b >> 16), A: uint16(a >> 16)}) != c.want16 {
			highHalf++
		}
	}
	models := []struct {
		model color.Model
		same  color.Color
	}{
		{color.RGBAModel, color.RGBA{R: 199, G: 23, B: 88, A: 7}},
		{color.RGBA64Model, color.RGBA64{R: 0xabcd, G: 0x0102, B: 0x9876, A: 0x0089}},
	}
	for _, c := range models {
		if got := c.model.Convert(c.same); got != c.same {
			log.Fatalf("same-type %T: got %v, want %v", c.same, got, c.same)
		}
		token := new(int)
		p := &callback{payload: token}
		if got := capture(c.model, p); got != token || p.calls != 1 {
			log.Fatalf("%T callback panic: got %v, calls %d", c.same, got, p.calls)
		}
		if p := capture(c.model, nil); p == nil {
			log.Fatalf("%T returned from nil Color", c.same)
		} else if _, ok := p.(runtime.Error); !ok {
			log.Fatalf("%T expected runtime.Error, got %T", c.same, p)
		}
	}
	if lowByte != 3 || highHalf != 3 {
		log.Fatalf("mutation counts: low byte=%d high half=%d", lowByte, highHalf)
	}
	fmt.Printf("RGBA models: 7 fixtures per model pass, mutant failures low byte=%d high half=%d\n", lowByte, highHalf)
}
