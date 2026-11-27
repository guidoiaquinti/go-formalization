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
	negativeZero := math.Float64frombits(1 << 63)
	cases := []struct{ x, want float64 }{
		{negativeZero, negativeZero}, {0, 0}, {-0.75, negativeZero}, {0.75, 0},
		{-1.75, -1}, {1.75, 1}, {-1, -1}, {math.Inf(1), math.Inf(1)},
	}
	var real, floor, positiveZero int
	for _, c := range cases {
		got, want := math.Float64bits(math.Trunc(c.x)), math.Float64bits(c.want)
		if got != want {
			real++
		}
		if math.Float64bits(math.Floor(c.x)) != want {
			floor++
		}
		if got == 1<<63 {
			got = 0
		}
		if got != want {
			positiveZero++
		}
	}
	if real != 0 || floor != 2 || positiveZero != 2 {
		log.Fatalf("Trunc mismatches: Go=%d floor=%d positive zero=%d", real, floor, positiveZero)
	}
	fmt.Printf("%s %s/%s: Trunc Go=%d/8 floor=%d/8 positive zero=%d/8\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, floor, positiveZero)
}
