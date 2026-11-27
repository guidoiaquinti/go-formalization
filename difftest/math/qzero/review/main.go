package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

// Coefficients copied in ascending order from pinned Go src/math/j0.go.
var bands = []struct{ p, q []float64 }{
	{[]float64{1.50444444886983272379e-07, 7.32234265963079278272e-02, 1.99819174093815998816e+00, 1.44956029347885735348e+01, 3.16662317504781540833e+01, 1.62527075710929267416e+01}, []float64{3.03655848355219184498e+01, 2.69348118608049844624e+02, 8.44783757595320139444e+02, 8.82935845112488550512e+02, 2.12666388511798828631e+02, -5.31095493882666946917e+00}},
	{[]float64{4.37741014089738620906e-09, 7.32411180042911447163e-02, 3.34423137516170720929e+00, 4.26218440745412650017e+01, 1.70808091340565596283e+02, 1.66733948696651168575e+02}, []float64{4.87588729724587182091e+01, 7.09689221056606015736e+02, 3.70414822620111362994e+03, 6.46042516752568917582e+03, 2.51633368920368957333e+03, -1.49247451836156386662e+02}},
	{[]float64{1.84085963594515531381e-11, 7.32421766612684765896e-02, 5.83563508962056953777e+00, 1.35111577286449829671e+02, 1.02724376596164097464e+03, 1.98997785864605384631e+03}, []float64{8.27766102236537761883e+01, 2.07781416421392987104e+03, 1.88472887785718085070e+04, 5.67511122894947329769e+04, 3.59767538425114471465e+04, -5.35434275601944773371e+03}},
	{[]float64{0.00000000000000000000e+00, 7.32421874999935051953e-02, 1.17682064682252693899e+01, 5.57673380256401856059e+02, 8.85919720756468632317e+03, 3.70146267776887834771e+04}, []float64{1.63776026895689824414e+02, 8.09834494656449805916e+03, 1.42538291419120476348e+05, 8.03309257119514397345e+05, 8.40501579819060512818e+05, -3.43899293537866615225e+05}},
}

func profileOutput(x float64, profile, band int) float64 {
	z := 1 / float64(x*x)
	mulAdd := func(a, b, c float64, fuse bool) float64 {
		if fuse {
			return math.FMA(a, b, c)
		}
		return float64(a*b) + c
	}
	horner := func(coefficients []float64) float64 {
		value := coefficients[len(coefficients)-1]
		for i := len(coefficients) - 2; i >= 0; i-- {
			value = mulAdd(z, value, coefficients[i], profile == 2)
		}
		return value
	}
	table := bands[band]
	return (-0.125 + horner(table.p)/mulAdd(z, horner(table.q), 1, profile != 0)) / x
}

func bandOf(x float64) int {
	switch {
	case x >= 8:
		return 3
	case x >= 4.5454:
		return 2
	case x >= 2.8571:
		return 1
	case x >= 2:
		return 0
	default:
		return -1
	}
}

func observation(x float64, profile int) string {
	band := bandOf(x)
	if band < 0 {
		return "panic"
	}
	return strconv.FormatUint(math.Float64bits(profileOutput(x, profile, band)), 10)
}

func main() {
	inputs := []float64{0, math.Copysign(0, -1), -2, math.Inf(-1), math.Inf(1), math.NaN(),
		math.Float64frombits(0x7ff0000000000001), math.Float64frombits(0xfff8000000000042), math.MaxFloat64}
	for _, cut := range []float64{2, 2.8571, 4.5454, 8, 0x1p512} {
		inputs = append(inputs, math.Nextafter(cut, 0), cut, math.Nextafter(cut, math.Inf(1)))
	}
	// Search directly in each normal coefficient band for a word whose
	// final rounding differs when one or both fused stages are disabled.
	rng := rand.New(rand.NewSource(65537))
	var differences [2]int
	for range 200000 {
		x := 2 + 9*rng.Float64()
		split, fused, arm := observation(x, 0), observation(x, 1), observation(x, 2)
		if split != fused && differences[0] < 8 {
			differences[0]++
			inputs = append(inputs, x)
		}
		if fused != arm && differences[1] < 8 {
			differences[1]++
			inputs = append(inputs, x)
		}
		if differences == [2]int{8, 8} {
			break
		}
	}
	var wrongBand, lostSign int
	if observation(math.Inf(1), 2) != "0" {
		lostSign++
	}
	for i, x := range []float64{2.2, 3.3, 6, 12} {
		inputs = append(inputs, x)
		if profileOutput(x, 2, i) != profileOutput(x, 2, (i+1)%4) {
			wrongBand++
		}
	}
	var rows bytes.Buffer
	for _, x := range inputs {
		for profile := range 3 {
			fmt.Fprintln(&rows, profile, math.Float64bits(x))
		}
	}
	file, err := os.CreateTemp("", "golean-qzero-review-")
	must(err)
	defer os.Remove(file.Name())
	fixtures := bufio.NewWriter(file)
	for i, got := range check.LeanRows("godiff-math-qzero", 3*len(inputs), &rows) {
		x, profile := inputs[i/3], i%3
		want := observation(x, profile)
		if got != want {
			log.Fatalf("qzero(%x), profile %d: Lean=%s oracle=%s", x, profile, got, want)
		}
		if profile == nativeProfile() {
			fmt.Fprintln(fixtures, math.Float64bits(x), got, want)
		}
	}
	must(fixtures.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_MATH_QZERO_CASES", file.Name()))
	const data = "difftest/math/qzero/testdata/"
	must(check.TestInPackage("math", "TestGoLeanQzero", data+"bridge_test.go", data+"qzero_test.go"))
	if differences[0] == 0 || differences[1] == 0 || wrongBand != 4 || lostSign != 1 {
		log.Fatalf("invalid variant strength: fusion=%v wrong band=%d", differences, wrongBand)
	}
	fmt.Printf("qzero: %d native and %d model observations; omitting denominator/Horner contraction, swapping coefficient bands, and canonicalizing zero failed %d/%d/%d/%d inputs\n",
		len(inputs), 3*len(inputs), differences[0], differences[1], wrongBand, lostSign)
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

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
