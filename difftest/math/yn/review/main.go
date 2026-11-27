package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

// Generate the two central solutions from a convergent series and its
// derivative, then propagate the exact order recurrence in float64.
func series(n int, x float64) float64 {
	const euler = 0.57721566490153286061
	size := -x * x / 4
	term, j, harmonic, correction, dj, dc := 1.0, 1.0, 0.0, 0.0, 0.0, 0.0
	for k := 1.0; k <= 40; k++ {
		term *= size / (k * k)
		harmonic += 1 / k
		j += term
		correction -= harmonic * term
		slope := 2 * k * term / x
		dj += slope
		dc -= harmonic * slope
	}
	logTerm := math.Log(x/2) + euler
	y0 := 2 / math.Pi * (logTerm*j + correction)
	y1 := -2 / math.Pi * (j/x + logTerm*dj + dc)
	for order := 0; order < n; order++ {
		y0, y1 = y1, float64(2*(order+1))/x*y1-y0
	}
	return y0
}

func nativeProfile() int {
	if runtime.GOARCH == "arm64" {
		return 2
	}
	if runtime.GOARCH != "amd64" {
		log.Fatalf("unsupported architecture %s", runtime.GOARCH)
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "GOAMD64" && (s.Value == "v3" || s.Value == "v4") {
				return 1
			}
		}
	}
	return 0
}

type input struct {
	order int
	word  uint64
}

func main() {
	tiny := math.Float64bits(1e-200)
	cases := []input{
		{1, 0x00028be60db93910}, {1, 0x00028be60db93911}, {2, tiny}, {3, tiny}, {-3, tiny}, {64, tiny}, {1, tiny},
		{0, 0}, {1, 1 << 63}, {3, 0x7ff0000000000000}, {3, 0xfff0000000000000},
		{3, 0x7ff8000000000000}, {-3, math.Float64bits(-1)},
	}
	const fixed = 13
	for _, n := range []int{0, 1, 2, 3, 7, 16, 32, 64} {
		for _, x := range []float64{0.25, 1.25, 3.25, 7.5} {
			cases = append(cases, input{n, math.Float64bits(x)})
		}
	}
	var in bytes.Buffer
	for profile := range 3 {
		for _, c := range cases {
			fmt.Fprintf(&in, "%d %d %d\n", profile, c.order, c.word)
		}
	}
	var native, analytic int
	for i, row := range check.LeanRows("godiff-math-yn", len(cases)*3, &in) {
		c := cases[i%len(cases)]
		profile := i / len(cases)
		x := math.Float64frombits(c.word)
		if profile == nativeProfile() {
			native++
			y := math.Yn(c.order, x)
			want := "NaN"
			if !math.IsNaN(y) {
				want = strconv.FormatUint(math.Float64bits(y), 10)
			}
			if row != want {
				log.Fatalf("n=%d x=%#x: native %s Lean %s", c.order, c.word, want, row)
			}
		}
		if c.word == tiny && (c.order >= 2 || c.order == -3) {
			sign := -1
			if c.order < 0 {
				sign = 1
			}
			want := strconv.FormatUint(math.Float64bits(math.Inf(sign)), 10)
			if row != want {
				log.Fatalf("profile %d n=%d tiny overflow: Lean %s, want %s", profile, c.order, row, want)
			}
		}
		for _, frontier := range []struct{ in, out uint64 }{
			{0x00028be60db93910, 0xfff0000000000000},
			{0x00028be60db93911, 0xffeffffffffffff8},
		} {
			if c.order == 1 && c.word == frontier.in && row != strconv.FormatUint(frontier.out, 10) {
				log.Fatalf("profile %d at order-one overflow frontier: %s", profile, row)
			}
		}
		if i%len(cases) < fixed {
			continue
		}
		bits, err := strconv.ParseUint(row, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		got, want := math.Float64frombits(bits), series(c.order, x)
		if !(math.Abs(got-want) <= math.Abs(want)*2e-10+1e-13) {
			log.Fatalf("profile %d n=%d x=%g: got %g want %g", profile, c.order, x, got, want)
		}
		analytic++
	}
	fmt.Printf("Yn: %d native, %d analytic and %d total Lean observations passed.\n", native, analytic, len(cases)*3)
}
