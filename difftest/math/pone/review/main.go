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

// Coefficients copied in ascending order from pinned Go src/math/j1.go.
var bands = []struct{ p, q []float64 }{
	{[]float64{1.07710830106873743082e-07, 1.17176219462683348094e-01, 2.36851496667608785174e+00, 1.22426109148261232917e+01, 1.76939711271687727390e+01, 5.07352312588818499250e+00}, []float64{2.14364859363821409488e+01, 1.25290227168402751090e+02, 2.32276469057162813669e+02, 1.17679373287147100768e+02, 8.36463893371618283368e+00}},
	{[]float64{3.02503916137373618024e-09, 1.17186865567253592491e-01, 3.93297750033315640650e+00, 3.51194035591636932736e+01, 9.10550110750781271918e+01, 4.85590685197364919645e+01}, []float64{3.47913095001251519989e+01, 3.36762458747825746741e+02, 1.04687139975775130551e+03, 8.90811346398256432622e+02, 1.03787932439639277504e+02}},
	{[]float64{1.31990519556243522749e-11, 1.17187493190614097638e-01, 6.80275127868432871736e+00, 1.08308182990189109773e+02, 5.17636139533199752805e+02, 5.28715201363337541807e+02}, []float64{5.92805987221131331921e+01, 9.91401418733614377743e+02, 5.35326695291487976647e+03, 7.84469031749551231769e+03, 1.50404688810361062679e+03}},
	{[]float64{0.00000000000000000000e+00, 1.17187499999988647970e-01, 1.32394806593073575129e+01, 4.12051854307378562225e+02, 3.87474538913960532227e+03, 7.91447954031891731574e+03}, []float64{1.14207370375678408436e+02, 3.65093083420853463394e+03, 3.69562060269033463555e+04, 9.76027935934950801311e+04, 3.08042720627888811578e+04}},
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
	if want := strconv.FormatUint(math.Float64bits(1), 10); observation(math.Inf(1), 2) != want {
		log.Fatal("invalid infinity fixture")
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
	file, err := os.CreateTemp("", "golean-pone-review-")
	must(err)
	defer os.Remove(file.Name())
	fixtures := bufio.NewWriter(file)
	for i, got := range check.LeanRows("godiff-math-pone", 3*len(inputs), &rows) {
		x, profile := inputs[i/3], i%3
		want := observation(x, profile)
		if got != want {
			log.Fatalf("pone(%x), profile %d: Lean=%s oracle=%s", x, profile, got, want)
		}
		if profile == nativeProfile() {
			fmt.Fprintln(fixtures, math.Float64bits(x), got, want)
		}
	}
	must(fixtures.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_MATH_PONE_CASES", file.Name()))
	const data = "difftest/math/pone/testdata/"
	must(check.TestInPackage("math", "TestGoLeanPone", data+"bridge_test.go", data+"pone_test.go"))
	if differences[0] == 0 || differences[1] == 0 || wrongBand != 4 {
		log.Fatalf("invalid variant strength: fusion=%v wrong band=%d", differences, wrongBand)
	}
	fmt.Printf("pone: %d native and %d model observations; omitting denominator/Horner contraction, swapping coefficient bands, failed %d/%d/%d inputs\n",
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
