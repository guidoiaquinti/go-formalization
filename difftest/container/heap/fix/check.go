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
func (*probe) Pop() any { panic("unexpected Pop") }

func main() {
	var edges []input
	for length := 1; length <= 33; length++ {
		for index := range length {
			for mask := range 16 {
				edges = append(edges, input{index, length, big.NewInt(int64(mask))})
			}
		}
	}
	full := new(big.Int).Lsh(big.NewInt(1), 128)
	full.Sub(full, big.NewInt(1))
	for _, index := range []int{0, math.MaxInt/2 - 1, math.MaxInt / 2, math.MaxInt - 1} {
		for _, mask := range []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(2), full} {
			edges = append(edges, input{index, math.MaxInt, mask})
		}
	}
	check.CompareSamples("godiff-container-heap-fix", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) input {
	var raw [16]byte
	r.Read(raw[:])
	length := r.Intn(math.MaxInt) + 1
	index := r.Intn(length)
	if i%3 == 0 {
		index = 0
	}
	return input{index, length, new(big.Int).SetBytes(raw[:])}
}

func encode(c input) string {
	return fmt.Sprintf("%d %d %s", c.index, c.length, c.decisions)
}

func native(c input) string {
	p := probe{input: c}
	heap.Fix(&p, c.index)
	return strings.TrimPrefix(p.trace.String(), " ")
}

func expected(c input) string {
	var trace strings.Builder
	trace.WriteString("len")
	step := 0
	less := func(a, b int) bool {
		fmt.Fprintf(&trace, " L:%d:%d", a, b)
		proceed := c.decisions.Bit(step) != 0
		step++
		return proceed
	}
	start := c.index + 1
	node := start
	for node <= c.length/2 {
		child := node * 2
		if child < c.length && less(child, child-1) {
			child++
		}
		if !less(child-1, node-1) {
			break
		}
		fmt.Fprintf(&trace, " S:%d:%d", node-1, child-1)
		node = child
	}
	if node == start {
		for node > 1 && less(node-1, node/2-1) {
			fmt.Fprintf(&trace, " S:%d:%d", node/2-1, node-1)
			node /= 2
		}
	}
	return trace.String()
}
