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
	const sign = uint64(1 << 63)
	const inf = uint64(0x7ff0000000000000)
	cases := []struct{ word, want uint64 }{
		{0, sign | inf}, {sign, sign | inf},
		{1, math.Float64bits(-1074)}, {sign | (1 << 32), math.Float64bits(-1042)},
		{1 << 52, math.Float64bits(-1022)}, {0x3ff0000000000000, 0},
		{sign | inf, inf}, {0xfff800000000002a, 0xfff800000000002a},
	}
	var real, exponentOnly int
	for _, c := range cases {
		got := math.Float64bits(math.Logb(math.Float64frombits(c.word)))
		if got != c.want {
			real++
		}
		mutant := got
		magnitude := c.word & (sign - 1)
		if magnitude != 0 && magnitude < 1<<52 {
			mutant = math.Float64bits(-1023)
		}
		if mutant != c.want {
			exponentOnly++
		}
	}
	if real != 0 || exponentOnly != 2 {
		log.Fatalf("Logb mismatches: Go=%d exponent-only=%d", real, exponentOnly)
	}
	fmt.Printf("%s %s/%s: Logb Go=%d/8 exponent-only mutant=%d/8\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, exponentOnly)
}
