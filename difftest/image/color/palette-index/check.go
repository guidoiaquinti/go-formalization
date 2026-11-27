package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const (
	returns = iota
	panics
	nilColor
)

type call struct {
	kind int
	rgba [4]uint32
}

type input struct {
	query   call
	palette []call
}

type probe struct {
	call
	id    int
	trace *[]int
}

func (p *probe) RGBA() (uint32, uint32, uint32, uint32) {
	*p.trace = append(*p.trace, p.id)
	if p.kind == panics {
		panic(p)
	}
	rgba := p.rgba
	return rgba[0], rgba[1], rgba[2], rgba[3]
}

func actual(c input) (row string) {
	var calls []int
	makeColor := func(event call, id int) color.Color {
		if event.kind == nilColor {
			return nil
		}
		return &probe{event, id, &calls}
	}
	palette := make(color.Palette, len(c.palette))
	for i, event := range c.palette {
		palette[i] = makeColor(event, i+1)
	}
	defer func() {
		switch p := recover().(type) {
		case *probe:
			row = format("panic", int(p.rgba[0]), calls)
		case runtime.Error:
			row = format("panic", 0, calls)
		case nil:
		default:
			panic(p)
		}
	}()
	index := palette.Index(makeColor(c.query, 0))
	return format("index", index, calls)
}

func format(status string, value int, calls []int) string {
	fields := []string{status, strconv.Itoa(value)}
	for _, id := range calls {
		fields = append(fields, strconv.Itoa(id))
	}
	return strings.Join(fields, " ")
}

func expected(c input) string {
	if c.query.kind == nilColor {
		return format("panic", 0, nil)
	}
	calls := []int{0}
	if c.query.kind == panics {
		return format("panic", int(c.query.rgba[0]), calls)
	}
	var scores []uint64
	for i, event := range c.palette {
		if event.kind == nilColor {
			return format("panic", 0, calls)
		}
		calls = append(calls, i+1)
		if event.kind == panics {
			return format("panic", int(event.rgba[0]), calls)
		}
		var score uint64
		for j, lhs := range c.query.rgba {
			delta := int64(lhs) - int64(event.rgba[j])
			if delta < 0 {
				delta = -delta
			}
			// The square of any uint32 difference fits in uint64.
			square := uint64(delta) * uint64(delta)
			score += square % (1 << 32) / 4
		}
		if score == 0 {
			return format("index", i, calls)
		}
		scores = append(scores, score)
	}
	best := 0
	for i := 1; i < len(scores); i++ {
		if scores[i] < scores[best] {
			best = i
		}
	}
	return format("index", best, calls)
}

func main() {
	zero := call{}
	maxWord := call{rgba: [4]uint32{0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff}}
	maxColor := call{rgba: [4]uint32{65535, 65535, 65535, 65535}}
	near := call{rgba: [4]uint32{1, 1, 1, 1}}
	tieR := call{rgba: [4]uint32{2, 0, 0, 0}}
	tieA := call{rgba: [4]uint32{0, 0, 0, 2}}
	panicEvent := call{kind: panics, rgba: [4]uint32{17}}
	nilEvent := call{kind: nilColor}
	edges := []input{
		{zero, nil},
		{panicEvent, nil},
		{nilEvent, nil},
		{zero, []call{near, panicEvent}},
		{zero, []call{maxColor, tieR, tieA}},
		{zero, []call{maxColor, tieA, tieR}},
		{zero, []call{maxColor, tieA, nilEvent}},
		{zero, []call{nilEvent, zero}},
		{zero, []call{tieR, panicEvent, zero}},
		{zero, []call{maxColor, maxWord}},
		{maxWord, []call{maxColor, zero, panicEvent}},
	}
	check.CompareSamples("godiff-image-color-palette-index", edges, func(r *rand.Rand, i int) input {
		randomColor := func() call {
			var c call
			for j := range c.rgba {
				c.rgba[j] = r.Uint32()
			}
			if i%2 == 0 {
				a := c.rgba[3] >> 16
				c.rgba[3] = a
				for j := range 3 {
					c.rgba[j] %= a + 1
				}
			}
			return c
		}
		query := randomColor()
		palette := make([]call, r.Intn(17))
		for j := range palette {
			palette[j] = randomColor()
		}
		if len(palette) != 0 {
			switch i % 5 {
			case 0:
				palette[r.Intn(len(palette))] = query
			case 1:
				palette[r.Intn(len(palette))] = panicEvent
			case 2:
				palette[r.Intn(len(palette))] = nilEvent
			case 3:
				palette[len(palette)-1] = palette[0]
			}
		}
		return input{query, palette}
	}, func(c input) string {
		var row strings.Builder
		fmt.Fprint(&row, len(c.palette))
		writeCall := func(event call) {
			fmt.Fprintf(&row, " %d %d %d %d %d", event.kind, event.rgba[0], event.rgba[1], event.rgba[2], event.rgba[3])
		}
		writeCall(c.query)
		for _, event := range c.palette {
			writeCall(event)
		}
		return row.String()
	}, actual, expected)
}
