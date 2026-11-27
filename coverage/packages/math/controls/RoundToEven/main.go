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
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}

	cases := []struct {
		input, want float64
	}{
		{0.5, 0},
		{-0.5, math.Copysign(0, -1)},
		{1.5, 2},
		{-1.5, -2},
		{2.5, 2},
		{-2.5, -2},
		{math.Nextafter(2.5, 2), 2},
		{math.Nextafter(2.5, 3), 3},
		{math.Nextafter(1<<52, 0), 1 << 52},
		{1<<52 + 1, 1<<52 + 1},
		{math.Inf(1), math.Inf(1)},
		{math.Inf(-1), math.Inf(-1)},
	}
	var awayFailures, truncationFailures int
	for _, test := range cases {
		want := math.Float64bits(test.want)
		if got := math.Float64bits(math.RoundToEven(test.input)); got != want {
			log.Fatalf("RoundToEven(%#016x) produced %#016x, want %#016x", math.Float64bits(test.input), got, want)
		}
		if math.Float64bits(math.Round(test.input)) != want {
			awayFailures++
		}
		// Round to nearest with every exact half rounded toward zero.
		magnitude := math.Abs(test.input)
		whole, fraction := math.Modf(magnitude)
		if fraction > 0.5 {
			whole++
		}
		if math.Float64bits(math.Copysign(whole, test.input)) != want {
			truncationFailures++
		}
	}
	if awayFailures != 4 || truncationFailures != 3 {
		log.Fatalf("tie mutant failures away=%d toward-zero=%d, want 4 and 3", awayFailures, truncationFailures)
	}
	for _, word := range []uint64{0x7ff0000000000001, 0xfff8000000000123} {
		if !math.IsNaN(math.RoundToEven(math.Float64frombits(word))) {
			log.Fatalf("RoundToEven(%#016x) did not return NaN", word)
		}
	}
	fmt.Printf("%s %s/%s: RoundToEven passed %d cases, tie mutants failed %d and %d of %d exact-word cases\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases)+2, awayFailures, truncationFailures, len(cases))
}
