package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"slices"
)

type marker struct {
	id     int
	value  uint32
	seen   *[]int
	action func()
}

func (m *marker) RGBA() (uint32, uint32, uint32, uint32) {
	*m.seen = append(*m.seen, m.id)
	if m.action != nil {
		m.action()
	}
	return m.value, m.value, m.value, m.value
}

type test struct {
	palette   color.Palette
	query     color.Color
	want      color.Color
	wantTrace []int
	caught    any
}
type convert func(color.Palette, color.Color) color.Color

func main() {
	payload := new(int)
	candidates := []struct {
		name string
		run  convert
		want int
	}{
		{"Palette.Convert", color.Palette.Convert, 0},
		{"Read palette slot before callbacks", cachedSlot, 3},
		{"Always visit the query", func(p color.Palette, c color.Color) color.Color { c.RGBA(); return p.Convert(c) }, 7},
		{"Swallow the panic", swallowPanic, 2},
	}
	for _, candidate := range candidates {
		mismatches := 0
		for mode := range 8 {
			if !matches(candidate.run, mode, payload) {
				mismatches++
				if candidate.want == 0 {
					log.Fatalf("Go disagreed on fixture %d", mode)
				}
			}
		}
		if mismatches != candidate.want {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, mismatches, candidate.want)
		}
		fmt.Printf("%s disagreed on %d of 8 fixtures\n", candidate.name, mismatches)
	}
}

func matches(candidate convert, mode int, payload *int) bool {
	var trace []int
	query := &marker{id: 0, seen: &trace}
	entry := &marker{id: 1, seen: &trace}
	later := &marker{id: 2, value: 3, seen: &trace}
	replacement := &marker{id: 3, seen: &trace, action: func() { panic("replacement was called") }}
	p := color.Palette{entry, later}
	var c test
	switch mode {
	case 0:
		query.action = func() { panic(payload) }
		c = test{nil, query, nil, nil, nil}
	case 1:
		c = test{color.Palette{}, nil, nil, nil, nil}
	case 2:
		entry.action = func() { p[0] = replacement }
		c = test{p, query, replacement, []int{0, 1}, nil}
	case 3:
		entry.action = func() { p[0] = nil }
		c = test{p, query, nil, []int{0, 1}, nil}
	case 4:
		entry.value = 2
		later.action = func() { p[0] = replacement }
		c = test{p, query, replacement, []int{0, 1, 2}, nil}
	case 5:
		entry.action = func() { panic(payload) }
		c = test{p, query, nil, []int{0, 1}, payload}
	case 6:
		c = test{p, nil, nil, nil, "runtime"}
	case 7:
		c = test{p, query, entry, []int{0, 1}, nil}
	}
	got, caught := observe(candidate, c.palette, c.query)
	if c.caught == "runtime" {
		if _, ok := caught.(runtime.Error); !ok {
			return false
		}
	} else if caught != c.caught {
		return false
	}
	return got == c.want && slices.Equal(trace, c.wantTrace)
}

func observe(run convert, p color.Palette, c color.Color) (result color.Color, caught any) {
	defer func() { caught = recover() }()
	return run(p, c), nil
}

func cachedSlot(p color.Palette, c color.Color) color.Color {
	if len(p) == 0 {
		return nil
	}
	old := slices.Clone(p)
	return old[p.Index(c)]
}

func swallowPanic(p color.Palette, c color.Color) color.Color {
	defer func() { recover() }()
	return p.Convert(c)
}
