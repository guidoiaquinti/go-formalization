package main

import (
	"fmt"
	"log"
	"math/bits"
)

// Evaluate the specification's count of bit positions via unsigned division.
func spec(x int) (n int) {
	for position := 0; position < 8; position++ {
		n += x % 2
		x /= 2
	}
	return n
}

// Mutate the clear-lowest-bit loop by giving it fuel for seven iterations.
func capped(x uint8) int {
	n := 0
	for x != 0 && n < 7 {
		x &= x - 1
		n++
	}
	return n
}

func check(label string, f func(uint8) int, wantFailures int) {
	failures := 0
	for x := 0; x < 256; x++ {
		got, want := f(uint8(x)), spec(x)
		if got != want {
			if failures == 0 {
				fmt.Printf("%s first counterexample x=%#02x, got=%d want=%d\n", label, x, got, want)
			}
			failures++
		}
	}
	fmt.Printf("%s failures %d/256\n", label, failures)
	if failures != wantFailures {
		log.Fatalf("%s: got %d failures, want %d", label, failures, wantFailures)
	}
}

func main() {
	check("go1.27.0", bits.OnesCount8, 0)
	check("seven iterations", capped, 1)
}
