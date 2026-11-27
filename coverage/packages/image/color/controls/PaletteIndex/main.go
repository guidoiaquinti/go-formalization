package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"slices"
)

type channels [4]uint32
type probe struct {
	id       int
	channels channels
	trace    *[]int
	payload  *panicToken
}
type panicToken struct{ id int }

func (p *probe) RGBA() (uint32, uint32, uint32, uint32) {
	*p.trace = append(*p.trace, p.id)
	if p.payload != nil {
		panic(p.payload)
	}
	return p.channels[0], p.channels[1], p.channels[2], p.channels[3]
}

type event struct {
	channels channels
	nilColor bool
	panics   bool
}
type fixture struct {
	query     event
	palette   []event
	wantIndex int
	wantTrace []int
	wantPanic string
}

func main() {
	zero := event{}
	panicEvent := event{panics: true}
	nilEvent := event{nilColor: true}
	oneEach := event{channels: channels{1, 1, 1, 1}}
	tieRed := event{channels: channels{2, 0, 0, 2}}
	tieGreen := event{channels: channels{0, 2, 0, 2}}
	wrapped := event{channels: channels{65536, 0, 0, 0}}
	far := event{channels: channels{3, 3, 3, 3}}
	cases := []fixture{
		{zero, nil, 0, []int{0}, ""},
		{nilEvent, []event{oneEach}, 0, nil, "runtime"},
		{panicEvent, []event{zero}, 0, []int{0}, "callback"},
		{zero, []event{oneEach, panicEvent}, 0, []int{0, 1}, ""},
		{zero, []event{far, tieRed, tieGreen}, 1, []int{0, 1, 2, 3}, ""},
		{zero, []event{far, nilEvent, zero}, 0, []int{0, 1}, "runtime"},
		{zero, []event{far, panicEvent, zero}, 0, []int{0, 1, 2}, "callback"},
		{zero, []event{far, wrapped, panicEvent}, 1, []int{0, 1, 2}, ""},
	}
	candidates := []struct {
		name string
		run  func(color.Palette, color.Color) int
		want int
	}{
		{"Palette.Index", color.Palette.Index, 0},
		{"prefer the last tie", func(p color.Palette, c color.Color) int { return mutant(p, c, true, false, false) }, 1},
		{"truncate after summing", func(p color.Palette, c color.Color) int { return mutant(p, c, false, true, false) }, 1},
		{"continue scanning after zero", func(p color.Palette, c color.Color) int { return mutant(p, c, false, false, true) }, 2},
	}
	for _, contestant := range candidates {
		mismatches := 0
		for i, test := range cases {
			if !matches(contestant.run, test) {
				mismatches++
				if contestant.want == 0 {
					log.Fatalf("Go disagreed on fixture %d", i)
				}
			}
		}
		if mismatches != contestant.want {
			log.Fatalf("%s disagreed on %d fixtures, want %d", contestant.name, mismatches, contestant.want)
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", contestant.name, mismatches, len(cases))
	}
}

func matches(index func(color.Palette, color.Color) int, test fixture) bool {
	var seen []int
	token := &panicToken{7}
	wrap := func(event event, id int) color.Color {
		if event.nilColor {
			return nil
		}
		var payload *panicToken
		if event.panics {
			payload = token
		}
		return &probe{id, event.channels, &seen, payload}
	}
	p := make(color.Palette, len(test.palette))
	for i, e := range test.palette {
		p[i] = wrap(e, i+1)
	}
	result := 0
	var caught any
	func() {
		defer func() { caught = recover() }()
		result = index(p, wrap(test.query, 0))
	}()
	if !slices.Equal(seen, test.wantTrace) {
		return false
	}
	switch test.wantPanic {
	case "runtime":
		_, ok := caught.(runtime.Error)
		return ok
	case "callback":
		return caught == token
	default:
		return caught == nil && result == test.wantIndex
	}
}

func mutant(p color.Palette, c color.Color, lastTie, lateTruncate, scanAfterZero bool) int {
	query := rgba(c)
	selected, best := 0, ^uint64(0)
	for i, entry := range p {
		candidate := rgba(entry)
		var score uint64
		for j, lhs := range query {
			d := lhs - candidate[j]
			square := uint64(d * d)
			if !lateTruncate {
				square /= 4
			}
			score += square
		}
		if lateTruncate {
			score /= 4
		}
		if score < best || lastTie && score == best {
			selected, best = i, score
			if score == 0 && !scanAfterZero {
				break
			}
		}
	}
	return selected
}

func rgba(c color.Color) channels {
	r, g, b, a := c.RGBA()
	return channels{r, g, b, a}
}
