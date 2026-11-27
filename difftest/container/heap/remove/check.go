package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	index, length int
	decisions     *big.Int
	value         uint64
}

type probe struct {
	input
	next  int
	trace strings.Builder
}

func (p *probe) Len() int {
	p.trace.WriteString(" len")
	return p.length
}

func (p *probe) Less(a, b int) bool {
	fmt.Fprintf(&p.trace, " L:%d:%d", a, b)
	proceed := p.decisions.Bit(p.next) != 0
	p.next++
	return proceed
}

func (p *probe) Swap(a, b int) {
	fmt.Fprintf(&p.trace, " S:%d:%d", a, b)
}

func (*probe) Push(any) { panic("unexpected Push") }

func (p *probe) Pop() any {
	p.trace.WriteString(" pop")
	return p.value
}

func main() {
	var edges []input
	for length := 1; length <= 33; length++ {
		for index := range length {
			for mask := range 16 {
				edges = append(edges, input{index, length, big.NewInt(int64(mask)), uint64(index + mask)})
			}
		}
	}
	full := new(big.Int).Lsh(big.NewInt(1), 128)
	full.Sub(full, big.NewInt(1))
	for _, index := range []int{0, math.MaxInt/2 - 1, math.MaxInt / 2, math.MaxInt - 1} {
		edges = append(edges, input{index, math.MaxInt, big.NewInt(0), 0},
			input{index, math.MaxInt, full, math.MaxUint64})
	}
	check.CompareSamples("godiff-container-heap-remove", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) input {
	var raw [16]byte
	r.Read(raw[:])
	if i%4 == 0 {
		for j := range raw {
			raw[j] = 0xff
		}
	}
	length := r.Intn(math.MaxInt) + 1
	index := r.Intn(length)
	if i%3 == 0 {
		index = length - 1
	} else if i%3 == 1 {
		index = 0
	}
	return input{index, length, new(big.Int).SetBytes(raw[:]), r.Uint64()}
}

func encode(c input) string {
	return fmt.Sprintf("%d %d %s %d", c.index, c.length, c.decisions, c.value)
}

func native(c input) string {
	p := probe{input: c}
	got := heap.Remove(&p, c.index).(uint64)
	return fmt.Sprint(got) + p.trace.String()
}

func expected(c input) string {
	var trace strings.Builder
	last := c.length - 1
	fmt.Fprintf(&trace, "%d len", c.value)
	if c.index == last {
		return trace.String() + " pop"
	}
	fmt.Fprintf(&trace, " S:%d:%d", c.index, last)
	step := 0
	less := func(a, b int) bool {
		fmt.Fprintf(&trace, " L:%d:%d", a, b)
		proceed := c.decisions.Bit(step) != 0
		step++
		return proceed
	}
	// Check parent eligibility before doubling the one-based index.
	initial := c.index + 1
	node := initial
	for node <= last/2 {
		child := node * 2
		if child < last && less(child, child-1) {
			child++
		}
		if !less(child-1, node-1) {
			break
		}
		fmt.Fprintf(&trace, " S:%d:%d", node-1, child-1)
		node = child
	}
	if node == initial {
		for node > 1 && less(node-1, node/2-1) {
			fmt.Fprintf(&trace, " S:%d:%d", node/2-1, node-1)
			node /= 2
		}
	}
	trace.WriteString(" pop")
	return trace.String()
}
