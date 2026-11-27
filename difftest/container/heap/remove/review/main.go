package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type item struct{ key, tag int }

type probe struct {
	items     []item
	trace     strings.Builder
	decisions uint64
	compares  uint
}

func (p *probe) Len() int {
	p.trace.WriteString(" len")
	return len(p.items)
}

func (p *probe) Less(i, j int) bool {
	fmt.Fprintf(&p.trace, " L:%d:%d", i, j)
	less := p.items[i].key < p.items[j].key
	if less {
		p.decisions |= 1 << p.compares
	}
	p.compares++
	return less
}

func (p *probe) Swap(i, j int) {
	fmt.Fprintf(&p.trace, " S:%d:%d", i, j)
	p.items[i], p.items[j] = p.items[j], p.items[i]
}

func (*probe) Push(any) { panic("unexpected Push") }

func (p *probe) Pop() any {
	p.trace.WriteString(" pop")
	last := len(p.items) - 1
	removed := p.items[last]
	p.items = p.items[:last]
	return removed
}

func main() {
	cases := []struct {
		keys  []int
		index int
		want  string
		tags  []int
	}{
		{[]int{9}, 0, "40 len pop", nil},
		{[]int{0, 4, 2}, 2, "42 len pop", []int{40, 41}},
		{[]int{0, 1, 2, 3}, 1, "41 len S:1:3 L:1:0 pop", []int{40, 43, 42}},
		{[]int{0, 5, 1, 6, 7, 2, 3}, 4, "44 len S:4:6 L:4:1 S:1:4 L:1:0 pop", []int{40, 46, 42, 43, 41, 45}},
		{[]int{0, 1, 2, 3, 4, 5, 6}, 0, "40 len S:0:6 L:2:1 L:1:0 S:0:1 L:4:3 L:3:1 S:1:3 pop", []int{41, 43, 42, 46, 44, 45}},
		{[]int{0, 1, 1}, 1, "41 len S:1:2 L:1:0 pop", []int{40, 42}},
	}
	var input bytes.Buffer
	failures := [3]int{}
	for i, c := range cases {
		p := newProbe(c.keys)
		removed := heap.Remove(p, c.index).(item)
		got := fmt.Sprint(removed.tag) + p.trace.String()
		tags := remainingTags(p)
		if removed != (item{c.keys[c.index], 40 + c.index}) || got != c.want || !reflect.DeepEqual(tags, c.tags) {
			log.Fatalf("Go fixture %d returned item %+v trace %q survivors %v, want trace %q survivors %v", i, removed, got, tags, c.want, c.tags)
		}
		for j := 1; j < len(p.items); j++ {
			if p.items[j].key < p.items[(j-1)/2].key {
				log.Fatalf("Go fixture %d broke the heap at index %d", i, j)
			}
		}
		fmt.Fprintln(&input, c.index, len(c.keys), p.decisions, 40+c.index)
		for mode := range failures {
			probe := newProbe(c.keys)
			value := mutatedRemove(probe, c.index, mode)
			if got := fmt.Sprint(value.tag) + probe.trace.String(); got != c.want || !reflect.DeepEqual(remainingTags(probe), c.tags) {
				failures[mode]++
			}
		}
	}
	for i, row := range check.LeanRows("godiff-container-heap-remove", len(cases), &input) {
		if row != cases[i].want {
			log.Fatalf("Lean fixture %d returned %q, want %q", i, row, cases[i].want)
		}
	}
	if failures != [3]int{2, 1, 3} {
		log.Fatalf("mutation counts %v, want [2 1 3]", failures)
	}
	fmt.Printf("heap.Remove passed %d tagged Go/Lean fixtures; forcing last swap, using full length, and dropping upward fallback failed %v.\n", len(cases), failures)
}

func newProbe(keys []int) *probe {
	p := new(probe)
	for i, key := range keys {
		p.items = append(p.items, item{key, 40 + i})
	}
	return p
}

func remainingTags(p *probe) []int {
	var result []int
	for _, element := range p.items {
		result = append(result, element.tag)
	}
	return result
}

func mutatedRemove(p *probe, i, mode int) item {
	last := p.Len() - 1
	if i != last || mode == 0 {
		p.Swap(i, last)
		count := last
		if mode == 1 {
			count++
		}
		start := i
		for 2*i+1 < count {
			child := 2*i + 1
			if child+1 < count && p.Less(child+1, child) {
				child++
			}
			if !p.Less(child, i) {
				break
			}
			p.Swap(i, child)
			i = child
		}
		if i == start && mode != 2 {
			for i > 0 && p.Less(i, (i-1)/2) {
				parent := (i - 1) / 2
				p.Swap(parent, i)
				i = parent
			}
		}
	}
	return p.Pop().(item)
}
