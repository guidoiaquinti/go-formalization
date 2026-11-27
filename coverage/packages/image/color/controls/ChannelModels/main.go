package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
)

type words [4]uint32

type probe struct {
	words
	calls   int
	failure *int
}

func (p *probe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	if p.failure != nil {
		panic(p.failure)
	}
	return p.words[0], p.words[1], p.words[2], p.words[3]
}

func try(model color.Model, input color.Color) (value color.Color, failure any) {
	defer func() { failure = recover() }()
	return model.Convert(input), nil
}

func main() {
	cases := []struct {
		words
		want [4]uint32
	}{
		{words{0, 0, 0, 0x12345678}, [4]uint32{86, 22136, 0, 0}},
		{words{65535, 65535, 65535, 0x89abcdef}, [4]uint32{205, 52719, 255, 65535}},
		{words{0x12345678, 0xabcdef01, 0xfedcba98, 0xf0000001}, [4]uint32{0, 1, 220, 56346}},
		{words{0, 436, 0, 0x01020304}, [4]uint32{3, 772, 1, 256}},
	}
	models := []struct {
		name         string
		model        color.Model
		fast         color.Color
		pack         func(uint32) color.Color
		mutant       func(words) uint32
		wantFailures int
	}{
		{"alphaModel", color.AlphaModel, color.Alpha{A: 0xd3},
			func(v uint32) color.Color { return color.Alpha{A: uint8(v)} },
			func(w words) uint32 { return w[3] & 255 }, 4},
		{"alpha16Model", color.Alpha16Model, color.Alpha16{A: 0xcdef},
			func(v uint32) color.Color { return color.Alpha16{A: uint16(v)} },
			func(w words) uint32 { return w[3] >> 16 }, 4},
		{"grayModel", color.GrayModel, color.Gray{Y: 0xa7},
			func(v uint32) color.Color { return color.Gray{Y: uint8(v)} },
			func(w words) uint32 { return (19595*w[0] + 38470*w[1] + 7471*w[2]) >> 24 }, 1},
		{"gray16Model", color.Gray16Model, color.Gray16{Y: 0x82ad},
			func(v uint32) color.Color { return color.Gray16{Y: uint16(v)} },
			func(w words) uint32 { return (19595*w[0] + 38470*w[1] + 7471*w[2]) >> 16 }, 2},
	}
	for i, m := range models {
		if got := m.model.Convert(m.fast); got != m.fast {
			log.Fatalf("%s fast path: got %#v", m.name, got)
		}
		failures := 0
		for _, c := range cases {
			p := &probe{words: c.words}
			if got := m.model.Convert(p); got != m.pack(c.want[i]) || p.calls != 1 {
				log.Fatalf("%s words=%#v: got %#v in %d calls", m.name, c.words, got, p.calls)
			}
			if m.mutant(c.words) != c.want[i] {
				failures++
			}
		}
		payload := new(int)
		p := &probe{failure: payload}
		if _, failure := try(m.model, p); failure != payload || p.calls != 1 {
			log.Fatalf("%s callback panic: recovered %#v in %d calls", m.name, failure, p.calls)
		}
		if _, failure := try(m.model, nil); failure == nil {
			log.Fatalf("%s returned from nil Color", m.name)
		} else if _, isRuntimeError := failure.(runtime.Error); !isRuntimeError {
			log.Fatalf("%s nil Color: recovered %T", m.name, failure)
		}
		if failures != m.wantFailures {
			log.Fatalf("%s mutation failures=%d, want %d", m.name, failures, m.wantFailures)
		}
		fmt.Printf("%s: 7 fixtures pass, arithmetic mutant fails %d/%d\n", m.name, failures, len(cases))
	}
}
