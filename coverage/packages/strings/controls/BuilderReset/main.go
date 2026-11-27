package main

import (
	"fmt"
	"reflect"
	"strings"
)

type state struct {
	identity bool
	buffer   bool
	length   int
	capacity int
}

func snapshot(b *strings.Builder) state {
	fields := reflect.ValueOf(b).Elem()
	return state{!fields.FieldByName("addr").IsNil(), !fields.FieldByName("buf").IsNil(), b.Len(), b.Cap()}
}

func main() {
	builders := make([]strings.Builder, 3)
	builders[1].Grow(32)
	builders[2].WriteString("original")
	var real, retainIdentity, retainCapacity int
	for i := range builders {
		b := &builders[i]
		saved := b.String()
		before := snapshot(b)
		b.Reset()
		if snapshot(b) != (state{}) {
			real++
		}
		if (state{identity: before.identity}) != (state{}) {
			retainIdentity++
		}
		if (state{buffer: before.buffer, capacity: before.capacity}) != (state{}) {
			retainCapacity++
		}
		b.WriteString("new")
		if b.String() != "new" || saved != []string{"", "", "original"}[i] {
			panic("Reset did not preserve the old string and allow reuse")
		}
	}
	fmt.Printf("Builder.Reset: Go=%d/3 retained-identity=%d/3 retained-capacity=%d/3\n",
		real, retainIdentity, retainCapacity)
	if real != 0 || retainIdentity != 2 || retainCapacity != 2 {
		panic("unexpected Builder.Reset control counts")
	}
}
