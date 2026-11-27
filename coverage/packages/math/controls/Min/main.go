package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
)

const (
	positiveZero uint64 = 0
	negativeZero uint64 = 1 << 63
	positiveInf  uint64 = 0x7ff0000000000000
	negativeInf  uint64 = negativeZero | positiveInf
	quietNaN     uint64 = 0x7ff8000000000001
	signalingNaN uint64 = 0x7ff0000000000001
	positiveOne  uint64 = 0x3ff0000000000000
	negativeOne  uint64 = negativeZero | positiveOne
	negativeTwo  uint64 = 0xc000000000000000
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected Go version %s", runtime.Version())
	}

	cases := []struct{ x, y, want uint64 }{
		{negativeInf, quietNaN, negativeInf},
		{signalingNaN, negativeInf, negativeInf},
		{positiveInf, quietNaN, quietNaN},
		{signalingNaN, positiveOne, quietNaN},
		{positiveZero, negativeZero, negativeZero},
		{negativeZero, positiveZero, negativeZero},
		{negativeZero, negativeZero, negativeZero},
		{positiveZero, positiveZero, positiveZero},
		{negativeOne, negativeTwo, negativeTwo},
		{positiveOne, negativeOne, negativeOne},
		{positiveInf, positiveOne, positiveOne},
	}
	var real, naNFirst, positiveZeroOnly int
	for _, c := range cases {
		got := math.Float64bits(math.Min(math.Float64frombits(c.x), math.Float64frombits(c.y)))
		if !matches(got, c.want) {
			real++
		}
		withNaNFirst := got
		if isNaN(c.x) || isNaN(c.y) {
			withNaNFirst = quietNaN
		}
		if !matches(withNaNFirst, c.want) {
			naNFirst++
		}
		withPositiveZero := got
		if got == negativeZero {
			withPositiveZero = positiveZero
		}
		if !matches(withPositiveZero, c.want) {
			positiveZeroOnly++
		}
	}
	if real != 0 || naNFirst != 2 || positiveZeroOnly != 3 {
		log.Fatalf("Min failures: Go=%d NaN-first=%d positive-zero-only=%d", real, naNFirst, positiveZeroOnly)
	}
	fmt.Printf("%s %s/%s: Min Go=%d/%d NaN-first=%d/%d positive-zero-only=%d/%d\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, len(cases),
		naNFirst, len(cases), positiveZeroOnly, len(cases))
}

func isNaN(word uint64) bool {
	return word&^negativeZero > positiveInf
}

func matches(got, want uint64) bool {
	if isNaN(want) {
		return isNaN(got)
	}
	return got == want
}
