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
	{[]float64{-8.87534333032526411254e-08, -7.03030995483624743247e-02, -1.45073846780952986357e+00, -7.63569613823527770791e+00, -1.11931668860356747786e+01, -3.23364579351335335033e+00}, []float64{2.22202997532088808441e+01, 1.36206794218215208048e+02, 2.70470278658083486789e+02, 1.53875394208320329881e+02, 1.46576176948256193810e+01}},
	{[]float64{-2.54704601771951915620e-09, -7.03119616381481654654e-02, -2.40903221549529611423e+00, -2.19659774734883086467e+01, -5.80791704701737572236e+01, -3.14479470594888503854e+01}, []float64{3.58560338055209726349e+01, 3.61513983050303863820e+02, 1.19360783792111533330e+03, 1.12799679856907414432e+03, 1.73580930813335754692e+02}},
	{[]float64{-1.14125464691894502584e-11, -7.03124940873599280078e-02, -4.15961064470587782438e+00, -6.76747652265167261021e+01, -3.31231299649172967747e+02, -3.46433388365604912451e+02}, []float64{6.07539382692300335975e+01, 1.05125230595704579173e+03, 5.97897094333855784498e+03, 9.62544514357774460223e+03, 2.40605815922939109441e+03}},
	{[]float64{0.00000000000000000000e+00, -7.03124999999900357484e-02, -8.08167041275349795626e+00, -2.57063105679704847262e+02, -2.48521641009428822144e+03, -5.25304380490729545272e+03}, []float64{1.16534364619668181717e+02, 3.83374475364121826715e+03, 4.05978572648472545552e+04, 1.16752972564375915681e+05, 4.76277284146730962675e+04}},
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
	return 1 + horner(table.p)/mulAdd(z, horner(table.q), 1, profile != 0)
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
	var wrongBand int
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
	file, err := os.CreateTemp("", "golean-pzero-review-")
	must(err)
	defer os.Remove(file.Name())
	fixtures := bufio.NewWriter(file)
	for i, got := range check.LeanRows("godiff-math-pzero", 3*len(inputs), &rows) {
		x, profile := inputs[i/3], i%3
		want := observation(x, profile)
		if got != want {
			log.Fatalf("pzero(%x), profile %d: Lean=%s oracle=%s", x, profile, got, want)
		}
		if profile == nativeProfile() {
			fmt.Fprintln(fixtures, math.Float64bits(x), got, want)
		}
	}
	must(fixtures.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_MATH_PZERO_CASES", file.Name()))
	const data = "difftest/math/pzero/testdata/"
	must(check.TestInPackage("math", "TestGoLeanPzero", data+"bridge_test.go", data+"pzero_test.go"))
	if differences[0] == 0 || differences[1] == 0 || wrongBand != 4 {
		log.Fatalf("invalid variant strength: fusion=%v wrong band=%d", differences, wrongBand)
	}
	fmt.Printf("pzero: %d native and %d model observations; omitting denominator/Horner contraction and swapping coefficient bands failed %d/%d/%d inputs\n",
		len(inputs), 3*len(inputs), differences[0], differences[1], wrongBand)
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
