package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

type grow func(*strings.Builder, int)

func panicValue(fn grow, b *strings.Builder, n int) (value any) {
	defer func() { value = recover() }()
	fn(b, n)
	return nil
}

func checks(fn grow) [4]bool {
	var b strings.Builder
	b.WriteString("abc")
	c := b.Cap()
	fn(&b, c-b.Len())
	fits := b.Cap() == c && b.String() == "abc"
	request := c - b.Len() + 1
	fn(&b, request)
	allocated := b.Cap() == cap(slices.Grow([]byte(nil), 2*c+request)) && b.String() == "abc"
	var zero strings.Builder
	neg := panicValue(fn, &zero, -1) == "strings.Builder.Grow: negative count"
	copyOfZero := zero
	const copyPanic = "strings: illegal use of non-zero Builder copied by value"
	identity := neg && panicValue(fn, &copyOfZero, 0) == copyPanic
	copyOfB := b
	priority := panicValue(fn, &copyOfB, -1) == copyPanic && copyOfB.String() == "abc"
	return [4]bool{fits, allocated, identity, priority}
}

func main() {
	for _, c := range []struct {
		name string
		fn   grow
		want int
	}{
		{"Go", (*strings.Builder).Grow, 0},
		{"Reject negative before copy check", func(b *strings.Builder, n int) {
			if n < 0 {
				panic("strings.Builder.Grow: negative count")
			}
			b.Grow(n)
		}, 2},
		{"Grow when request exactly fits", func(b *strings.Builder, n int) {
			if b != nil && n >= 0 && n == b.Cap()-b.Len() {
				n++
			}
			b.Grow(n)
		}, 1},
	} {
		bad := 0
		for _, pass := range checks(c.fn) {
			if !pass {
				bad++
			}
		}
		fmt.Printf("%s disagreed on %d of 4 checks\n", c.name, bad)
		if bad != c.want {
			log.Fatalf("%s wanted %d failures", c.name, c.want)
		}
	}
}
