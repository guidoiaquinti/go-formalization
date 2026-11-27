package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	value     uint64
	before    int
	decisions uint64
}

type probe struct {
	length    int
	decisions uint64
	trace     strings.Builder
}

func (p *probe) Push(x any) {
	fmt.Fprintf(&p.trace, "P:%d", x.(uint64))
	p.length++
}

func (p *probe) Len() int {
	p.trace.WriteString(" len")
	return p.length
}

func (p *probe) Less(child, parent int) bool {
	fmt.Fprintf(&p.trace, " L:%d:%d", child, parent)
	proceed := p.decisions&1 != 0
	p.decisions >>= 1
	return proceed
}

func (p *probe) Swap(parent, child int) {
	fmt.Fprintf(&p.trace, " S:%d:%d", parent, child)
}

func (*probe) Pop() any { panic("unexpected Pop") }

func main() {
	var edges []input
	for n := range 128 {
		for decisions := range 16 {
			edges = append(edges, input{uint64(n), n, uint64(decisions)})
		}
	}
	for _, n := range []int{128, 1023, 1024, math.MaxInt - 1} {
		edges = append(edges, input{0, n, 0}, input{math.MaxUint64, n, math.MaxUint64})
	}
	check.CompareSamples("godiff-container-heap-push", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) input {
	decisions := r.Uint64()
	if i%4 == 0 {
		decisions = math.MaxUint64
	}
	return input{r.Uint64(), r.Intn(math.MaxInt), decisions}
}

func encode(c input) string {
	return fmt.Sprintf("%d %d %d", c.value, c.before+1, c.decisions)
}

func native(c input) string {
	p := probe{length: c.before, decisions: c.decisions}
	heap.Push(&p, c.value)
	return p.trace.String()
}

func expected(c input) string {
	var trace strings.Builder
	fmt.Fprintf(&trace, "P:%d len", c.value)
	// Dropping low bits of the one-based new index enumerates its ancestors.
	for depth, node := 0, uint64(c.before)+1; node > 1; depth, node = depth+1, node>>1 {
		child, parent := node-1, (node>>1)-1
		fmt.Fprintf(&trace, " L:%d:%d", child, parent)
		if c.decisions&(1<<depth) == 0 {
			break
		}
		fmt.Fprintf(&trace, " S:%d:%d", parent, child)
	}
	return trace.String()
}
