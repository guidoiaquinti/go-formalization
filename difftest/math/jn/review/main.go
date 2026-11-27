package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

// The convergent defining series uses exact rational terms for binary64 x.
func series(n int, x float64) float64 {
	xRat := new(big.Rat).SetFloat64(x)
	half := new(big.Rat).Quo(xRat, big.NewRat(2, 1))
	term := big.NewRat(1, 1)
	for k := 1; k <= n; k++ {
		term.Mul(term, half).Quo(term, big.NewRat(int64(k), 1))
	}
	square := new(big.Rat).Mul(half, half)
	square.Neg(square)
	sum := new(big.Rat).Set(term)
	for k := 1; k <= 80; k++ {
		term.Mul(term, square).Quo(term, big.NewRat(int64(k*(k+n)), 1))
		sum.Add(sum, term)
	}
	rounded, _ := sum.Float64()
	return rounded
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
	const root = 0x40033d152e971b40
	cases := []input{
		{3, root - 1}, {3, root}, {3, root + 1}, {64, root}, {-64, root},
		{-3, root}, {3, root | 1<<63}, {0, root}, {1, root},
		{1, 1 | 1<<63}, {2, 0}, {3, 1 << 63}, {3, 0x7ff0000000000000},
		{3, 0x7ff8000000000000}, {-3, 0xfff0000000000000},
	}
	for _, n := range []int{0, 1, 2, 3, 7, 16, 32, 64} {
		for _, x := range []float64{0.25, 1.25, 3.25, 7.5, 15.5} {
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
	for i, row := range check.LeanRows("godiff-math-jn", len(cases)*3, &in) {
		c := cases[i%len(cases)]
		profile := i / len(cases)
		x := math.Float64frombits(c.word)
		if profile == nativeProfile() {
			native++
			y := math.Jn(c.order, x)
			want := "NaN"
			if !math.IsNaN(y) {
				want = strconv.FormatUint(math.Float64bits(y), 10)
			}
			if row != want {
				log.Fatalf("n=%d x=%#x: native %s Lean %s", c.order, c.word, want, row)
			}
		}
		if c.order == 3 && c.word == root {
			if !(math.Abs(series(c.order, x)-0.19899990535769083) < 1e-15) {
				log.Fatal("incorrect independent root witness")
			}
			if profile == 2 && row != strconv.FormatUint(math.Float64bits(math.Inf(-1)), 10) {
				log.Fatalf("root witness: Lean %s", row)
			}
		}
		if i%len(cases) < 15 {
			continue
		}
		bits, err := strconv.ParseUint(row, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		got, want := math.Float64frombits(bits), series(c.order, x)
		if !(math.Abs(got-want) <= math.Abs(want)*1e-12+1e-150) {
			log.Fatalf("profile %d n=%d x=%g: got %g want %g", profile, c.order, x, got, want)
		}
		analytic++
	}
	fmt.Printf("Jn: %d native, %d analytic and %d total Lean observations passed.\n", native, analytic, len(cases)*3)
}
