package main

import (
	"container/heap"
	"fmt"
	"math/big"
	"math/rand"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	length    int
	decisions *big.Int
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
	full := new(big.Int).Lsh(big.NewInt(1), 256)
	full.Sub(full, big.NewInt(1))
	alternating := new(big.Int).Div(full, big.NewInt(3))
	masks := []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(2), big.NewInt(3),
		big.NewInt(5), big.NewInt(10), full, alternating}
	var edges []input
	for length := range 129 {
		for _, mask := range masks {
			edges = append(edges, input{length, mask})
		}
	}
	check.CompareSamples("godiff-container-heap-init", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, _ int) input {
	var raw [32]byte
	r.Read(raw[:])
	return input{r.Intn(65), new(big.Int).SetBytes(raw[:])}
}

func encode(c input) string {
	return fmt.Sprintf("%d %s", c.length, c.decisions)
}

func native(c input) string {
	p := probe{input: c}
	heap.Init(&p)
	return strings.TrimPrefix(p.trace.String(), " ")
}

func expected(c input) string {
	var trace strings.Builder
	trace.WriteString("len")
	probeIndex := 0
	less := func(a, b int) bool {
		fmt.Fprintf(&trace, " L:%d:%d", a, b)
		proceed := c.decisions.Bit(probeIndex) != 0
		probeIndex++
		return proceed
	}

	// Visit parents in reverse one-based order. Parent eligibility is
	// checked before doubling so the child arithmetic stays in range.
	for start := c.length / 2; start > 0; start-- {
		for node := start; node <= c.length/2; {
			child := 2 * node
			if child < c.length && less(child, child-1) {
				child++
			}
			if !less(child-1, node-1) {
				break
			}
			fmt.Fprintf(&trace, " S:%d:%d", node-1, child-1)
			node = child
		}
	}
	return trace.String()
}
