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
	cases := []struct {
		word uint64
		want int
	}{
		{0, math.MinInt32}, {1 << 63, math.MinInt32},
		{1, -1074}, {1<<63 | (1 << 51), -1023},
		{1 << 52, -1022}, {0x3ff0000000000000, 0},
		{0x7fefffffffffffff, 1023}, {0x7ff0000000000000, math.MaxInt32},
		{0xfff0000000000000, math.MaxInt32}, {0x7ff0000000000001, math.MaxInt32},
	}
	var real, exponentOnly int
	for _, c := range cases {
		got := math.Ilogb(math.Float64frombits(c.word))
		if got != c.want {
			real++
		}
		mutant := got
		magnitude := c.word & (1<<63 - 1)
		if magnitude > 0 && magnitude < 1<<52 {
			mutant = -1023
		}
		if mutant != c.want {
			exponentOnly++
		}
	}
	if real != 0 || exponentOnly != 1 {
		log.Fatalf("Ilogb mismatches: Go=%d subnormal exponent mutant=%d", real, exponentOnly)
	}
	fmt.Printf("%s %s/%s: Ilogb Go=%d/10 subnormal exponent mutant=%d/10\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, exponentOnly)
}
