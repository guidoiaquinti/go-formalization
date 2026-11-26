package main

import (
	"fmt"
	"log"
	"math/bits"
	"runtime"
)

func main() {
	fmt.Println(runtime.Version(), runtime.GOOS, runtime.GOARCH)
	for _, impl := range []struct {
		name           string
		f              func(uint8) int
		zero, failures int
	}{
		{"Go", bits.TrailingZeros8, 8, 0},
		{"seven-bit scan", func(x uint8) int {
			n := 0
			for n < 7 && x&1 == 0 {
				n++
				x >>= 1
			}
			return n
		}, 7, 0},
		{"wrong direction", bits.LeadingZeros8, 8, 170},
	} {
		z := impl.f(0)
		bad := 0
		for i := 1; i < 256; i++ {
			n := impl.f(uint8(i))
			if n < 0 || n > 8 {
				panic(n)
			}
			p := 1 << n
			r := i % (2 * p)
			if r != p {
				if bad == 0 {
					fmt.Printf("%s: first nonzero failure: x=%d n=%d remainder=%d want=%d\n", impl.name, i, n, r, p)
				}
				bad++
			}
		}
		fmt.Printf("%s: f(0)=%d zero passes=%v nonzero failures=%d/255\n", impl.name, z, z == 8, bad)
		if z != impl.zero || bad != impl.failures {
			log.Fatalf("%s: expected f(0)=%d and %d nonzero failures", impl.name, impl.zero, impl.failures)
		}
	}
}
