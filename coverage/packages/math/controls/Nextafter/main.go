package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected Go version %s", runtime.Version())
	}
	check("Nextafter", uint64(1<<63), uint64(0x3ff0000000000000),
		uint64(0x7ff0000000000000), uint64(0x7ff8000000000001), func(from, toward uint64) uint64 {
			return math.Float64bits(math.Nextafter(math.Float64frombits(from), math.Float64frombits(toward)))
		})
	check("Nextafter32", uint32(1<<31), uint32(0x3f800000),
		uint32(0x7f800000), uint32(0x7fc00000), func(from, toward uint32) uint32 {
			return math.Float32bits(math.Nextafter32(math.Float32frombits(from), math.Float32frombits(toward)))
		})
}

func check[T uint32 | uint64](name string, negativeZero, one, inf, nan T, next func(T, T) T) {
	cases := []struct{ from, toward, want T }{
		{negativeZero, 0, negativeZero},
		{negativeZero + 1, one, negativeZero},
		{1, negativeZero | one, 0},
		{0, negativeZero | one, negativeZero + 1},
		{one, inf, one + 1},
		{inf, 0, inf - 1},
		{negativeZero | inf, one, negativeZero | (inf - 1)},
		{one, inf + 2, nan},
	}
	var goFailures, positiveZeroFailures int
	for _, c := range cases {
		got := next(c.from, c.toward)
		if got != c.want {
			goFailures++
		}
		if got == negativeZero {
			got = 0
		}
		if got != c.want {
			positiveZeroFailures++
		}
	}
	if goFailures != 0 || positiveZeroFailures != 2 {
		log.Fatalf("%s mismatches: Go=%d positive-zero mutant=%d", name, goFailures, positiveZeroFailures)
	}
	fmt.Printf("%s %s/%s: %s Go=%d/8 positive-zero mutant=%d/8\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, name, goFailures, positiveZeroFailures)
}
