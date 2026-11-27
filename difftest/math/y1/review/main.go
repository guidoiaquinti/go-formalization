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

// Differentiate the convergent Y0 series and use Y1 = -Y0'.
func central(x float64) float64 {
	const euler = 0.57721566490153286061
	square := -0.25 * x * x
	term, j, derivative, correctionDerivative, harmonic := 1.0, 1.0, 0.0, 0.0, 0.0
	for k := 1.0; k <= 40; k++ {
		term *= square / (k * k)
		harmonic += 1 / k
		j += term
		slope := (2 * k / x) * term
		derivative += slope
		correctionDerivative -= harmonic * slope
	}
	return -2 / math.Pi * (j/x + (math.Log(x/2)+euler)*derivative + correctionDerivative)
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
					return 1
				}
			}
		}
		return 0
	default:
		log.Fatal("Y1 reviewer requires a supported architecture")
		return 0
	}
}

func main() {
	inputs := []uint64{0, 1 << 63, 1, 1 | 1<<63, 0x7ff0000000000000, 0xfff0000000000000,
		0x7ff0000000000101, 0xfff8000000000042, 0x0010000000000000, 0x000fffffffffffff, 0x00028be60db93910, 0x00028be60db93911}
	for _, x := range []float64{0x1p-27, 0x1p-13, 1, 2, 2.404825557695773, 3 * math.Pi / 4,
		5.520078110286311, 7 * math.Pi / 4, 17.8, 0x1p129, math.MaxFloat64 / 2, math.MaxFloat64} {
		for _, nearby := range []float64{math.Nextafter(x, 0), x, math.Nextafter(x, math.Inf(1))} {
			inputs = append(inputs, math.Float64bits(nearby), math.Float64bits(-nearby))
		}
	}
	var source bytes.Buffer
	for profile := range 3 {
		for _, word := range inputs {
			fmt.Fprintf(&source, "%d %d\n", profile, word)
		}
	}
	profile := nativeProfile()
	var exact, analytic int
	rows := check.LeanRows("godiff-math-y1", 3*len(inputs), &source)
	for i, row := range rows {
		word := inputs[i%len(inputs)]
		x := math.Float64frombits(word)
		p := i / len(inputs)
		host := math.Y1(x)
		if p == profile {
			exact++
			want := "NaN"
			if !math.IsNaN(host) {
				want = strconv.FormatUint(math.Float64bits(host), 10)
			}
			if row != want {
				log.Fatalf("native mismatch at %016x: Lean %q Go %q", word, row, want)
			}
		}
		if math.IsNaN(x) || x < 0 {
			if row != "NaN" {
				log.Fatalf("NaN word %016x produced %q", word, row)
			}
			continue
		}
		bits, err := strconv.ParseUint(row, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		if math.IsInf(x, 0) && bits != 0 {
			log.Fatalf("infinite word %016x produced %016x", word, bits)
		}
		if x == 0 && bits != math.Float64bits(math.Inf(-1)) {
			log.Fatalf("zero word %016x produced %016x", word, bits)
		}
		if word == 0x00028be60db93910 && bits != math.Float64bits(math.Inf(-1)) {
			log.Fatalf("profile %d last infinite result: %#x", p, bits)
		}
		if word == 0x00028be60db93911 && bits != 0xffeffffffffffff8 {
			log.Fatalf("profile %d first finite result: %#x", p, bits)
		}
		magnitude := math.Abs(x)
		if magnitude >= 0x1p-25 && magnitude <= 8 {
			analytic++
			error := math.Abs(math.Float64frombits(bits)-central(magnitude)) * math.Sqrt(max(1, magnitude))
			if !(error <= 0x1p-36*max(1, math.Abs(central(magnitude)))) {
				log.Fatalf("profile %d input %016x central series discrepancy %g", p, word, error)
			}
		}
	}
	fmt.Printf("Y1 review: %d exact native rows, %d analytic profile rows and %d total Lean rows passed.\n",
		exact, analytic, len(rows))
}
