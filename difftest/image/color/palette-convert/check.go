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

type channels [4]uint32

const (
	returns = iota
	panics
	nilColor
)

type callback struct {
	kind int
	rgba channels
}

type entry struct {
	callback
	replacement int
}

const (
	unchanged = iota
	newColor
	replaceWithNil
)

type input struct {
	query      callback
	entries    []entry
	nilPalette bool
}

type probe struct {
	callback
	trace   *[]int
	id      int
	replace func()
}

func (p *probe) RGBA() (uint32, uint32, uint32, uint32) {
	*p.trace = append(*p.trace, p.id)
	if p.replace != nil {
		p.replace()
	}
	if p.kind == panics {
		panic(p)
	}
	return p.rgba[0], p.rgba[1], p.rgba[2], p.rgba[3]
}

func resultToken(e entry, i int) int {
	switch e.replacement {
	case newColor:
		return i + 1 + 100
	case replaceWithNil:
		return 0
	default:
		return i + 1
	}
}

func actual(c input) (row string) {
	var calls []int
	var p color.Palette
	if !c.nilPalette {
		p = make(color.Palette, len(c.entries))
	}
	makeColor := func(event callback, id int) color.Color {
		if event.kind == nilColor {
			return nil
		}
		return &probe{callback: event, trace: &calls, id: id}
	}
	identity := make(map[color.Color]int)
	for i, e := range c.entries {
		p[i] = makeColor(e.callback, i+1)
		if original, ok := p[i].(*probe); ok {
			identity[original] = i + 1
			if e.replacement == replaceWithNil {
				original.replace = func() { p[i] = nil }
			}
			if e.replacement == newColor {
				replacement := makeColor(callback{kind: panics, rgba: channels{37}}, i+1+100)
				identity[replacement] = i + 1 + 100
				original.replace = func() { p[i] = replacement }
			}
		}
	}
	defer func() {
		switch failure := recover().(type) {
		case *probe:
			row = format("panic", failure.rgba[0], calls)
		case runtime.Error:
			row = format("panic", 0, calls)
		case nil:
		default:
			panic(failure)
		}
	}()
	value := p.Convert(makeColor(c.query, 0))
	if value == nil {
		return format("nil", 0, calls)
	}
	token, found := identity[value]
	if !found {
		return "unknown identity"
	}
	return format("color", uint32(token), calls)
}

func format(kind string, value uint32, calls []int) string {
	fields := []string{kind}
	if kind != "nil" {
		fields = append(fields, strconv.FormatUint(uint64(value), 10))
	}
	for _, call := range calls {
		fields = append(fields, strconv.Itoa(call))
	}
	return strings.Join(fields, " ")
}

func expected(c input) string {
	if len(c.entries) == 0 {
		return format("nil", 0, nil)
	}
	if c.query.kind == nilColor {
		return format("panic", 0, nil)
	}
	calls := []int{0}
	if c.query.kind == panics {
		return format("panic", c.query.rgba[0], calls)
	}
	bestIndex := 0
	bestScore := ^uint64(0)
	for i, e := range c.entries {
		if e.kind == nilColor {
			return format("panic", 0, calls)
		}
		calls = append(calls, i+1)
		if e.kind == panics {
			return format("panic", e.rgba[0], calls)
		}
		var score uint64
		for j, lhs := range c.query.rgba {
			diff := int64(lhs) - int64(e.rgba[j])
			score += uint64(diff * diff / 4)
		}
		if score < bestScore {
			bestIndex, bestScore = i, score
		}
		if score == 0 {
			break
		}
	}
	token := resultToken(c.entries[bestIndex], bestIndex)
	if token == 0 {
		return format("nil", 0, calls)
	}
	return format("color", uint32(token), calls)
}

func main() {
	clear := callback{}
	one := callback{rgba: channels{1, 1, 1, 1}}
	dark := callback{rgba: channels{0, 0, 0, 65535}}
	bright := callback{rgba: channels{65535, 65535, 65535, 65535}}
	crash := callback{kind: panics, rgba: channels{17}}
	missing := callback{kind: nilColor}
	edges := []input{
		{query: clear},
		{query: crash, nilPalette: true},
		{query: crash},
		{query: missing},
		{query: missing, entries: []entry{{callback: dark}}},
		{query: crash, entries: []entry{{callback: dark}}},
		{query: clear, entries: []entry{{callback: one}, {callback: crash}}},
		{query: clear, entries: []entry{{callback: bright}, {callback: dark}, {callback: dark}}},
		{query: clear, entries: []entry{{callback: dark}, {callback: crash}, {callback: clear}}},
		{query: clear, entries: []entry{{callback: dark}, {callback: missing}}},
		{query: clear, entries: []entry{{callback: clear, replacement: newColor}, {callback: crash}}},
		{query: clear, entries: []entry{{callback: clear, replacement: replaceWithNil}}},
		{query: clear, entries: []entry{{callback: bright}, {callback: dark, replacement: newColor}}},
		{query: clear, entries: []entry{{callback: dark}, {callback: bright, replacement: replaceWithNil}}},
	}
	check.CompareSamples("godiff-image-color-palette-convert", edges, func(r *rand.Rand, i int) input {
		color := func() callback {
			alpha := r.Uint32() >> 16
			return callback{rgba: channels{r.Uint32() % (alpha + 1), r.Uint32() % (alpha + 1), r.Uint32() % (alpha + 1), alpha}}
		}
		query := color()
		entries := make([]entry, r.Intn(13))
		for j := range entries {
			entries[j] = entry{color(), r.Intn(3)}
		}
		if len(entries) != 0 {
			j := r.Intn(len(entries))
			switch i % 6 {
			case 0:
				entries[j].callback = query
			case 1:
				entries[j] = entry{callback: crash}
			case 2:
				entries[j] = entry{callback: missing}
			case 3:
				entries[len(entries)-1].callback = entries[0].callback
			case 4:
				query = crash
			}
		}
		return input{query, entries, len(entries) == 0 && i%2 == 0}
	}, func(c input) string {
		var row strings.Builder
		writeCall := func(event callback) {
			fmt.Fprintf(&row, "%d %d %d %d %d", event.kind, event.rgba[0], event.rgba[1], event.rgba[2], event.rgba[3])
		}
		writeCall(c.query)
		for i, e := range c.entries {
			fmt.Fprintf(&row, " %d ", resultToken(e, i))
			writeCall(e.callback)
		}
		return row.String()
	}, actual, expected)
}
