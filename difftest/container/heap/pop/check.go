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
	length    int
	decisions *big.Int
	value     uint64
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
	for length := range 129 {
		for mask := range 16 {
			edges = append(edges, input{length, big.NewInt(int64(mask)), uint64(length + mask)})
		}
	}
	full := new(big.Int).Lsh(big.NewInt(1), 128)
	full.Sub(full, big.NewInt(1))
	for _, length := range []int{129, 1024, math.MaxInt - 1, math.MaxInt} {
		edges = append(edges, input{length, big.NewInt(0), 0}, input{length, full, math.MaxUint64})
	}
	check.CompareSamples("godiff-container-heap-pop", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) input {
	var raw [16]byte
	r.Read(raw[:])
	if i%4 == 0 {
		for j := range raw {
			raw[j] = 0xff
		}
	}
	return input{r.Int(), new(big.Int).SetBytes(raw[:]), r.Uint64()}
}

func encode(c input) string {
	return fmt.Sprintf("%d %s %d", c.length, c.decisions, c.value)
}

func native(c input) string {
	p := probe{input: c}
	got := heap.Pop(&p).(uint64)
	return fmt.Sprint(got) + p.trace.String()
}

func expected(c input) string {
	var trace strings.Builder
	last := c.length - 1
	fmt.Fprintf(&trace, "%d len S:0:%d", c.value, last)
	step := 0
	less := func(a, b int) bool {
		fmt.Fprintf(&trace, " L:%d:%d", a, b)
		proceed := c.decisions.Bit(step) != 0
		step++
		return proceed
	}
	for node := 1; node <= last/2; {
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
	trace.WriteString(" pop")
	return trace.String()
}
