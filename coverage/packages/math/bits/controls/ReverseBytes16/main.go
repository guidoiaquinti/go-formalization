package main

import (
	"fmt"
	"math/bits"
)

func main() {
	for _, t := range []struct {
		name                      string
		impl                      func(uint16) uint16
		specFailures, invFailures int
	}{
		{"Go", bits.ReverseBytes16, 0, 0},
		{"identity", func(x uint16) uint16 { return x }, 65280, 0},
		{"drop upper output byte", func(x uint16) uint16 { return x >> 8 }, 65280, 65535},
	} {
		specFailures, invFailures := 0, 0
		for i := range 1 << 16 {
			x := uint16(i)
			want := uint16((i%256)*256 + i/256)
			got := t.impl(x)
			if got != want {
				specFailures++
			}
			if t.impl(got) != x {
				invFailures++
			}
		}
		fmt.Printf("%s: %d spec failures, %d involution failures\n", t.name, specFailures, invFailures)
		if specFailures != t.specFailures || invFailures != t.invFailures {
			panic(t.name)
		}
	}
}
