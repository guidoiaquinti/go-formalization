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
	{[]float64{-1.78381727510958865572e-07, -1.02517042607985553460e-01, -2.75220568278187460720e+00, -1.96636162643703720221e+01, -4.23253133372830490089e+01, -2.13719211703704061733e+01}, []float64{2.95333629060523854548e+01, 2.52981549982190529136e+02, 7.57502834868645436472e+02, 7.39393205320467245656e+02, 1.55949003336666123687e+02, -4.95949898822628210127e+00}},
	{[]float64{-5.07831226461766561369e-09, -1.02537829820837089745e-01, -4.61011581139473403113e+00, -5.78472216562783643212e+01, -2.28244540737631695038e+02, -2.19210128478909325622e+02}, []float64{4.76651550323729509273e+01, 6.73865112676699709482e+02, 3.38015286679526343505e+03, 5.54772909720722782367e+03, 1.90311919338810798763e+03, -1.35201191444307340817e+02}},
	{[]float64{-2.08979931141764104297e-11, -1.02539050241375426231e-01, -8.05644828123936029840e+00, -1.83669607474888380239e+02, -1.37319376065508163265e+03, -2.61244440453215656817e+03}, []float64{8.12765501384335777857e+01, 1.99179873460485964642e+03, 1.74684851924908907677e+04, 4.98514270910352279316e+04, 2.79480751638918118260e+04, -4.71918354795128470869e+03}},
	{[]float64{0.00000000000000000000e+00, -1.02539062499992714161e-01, -1.62717534544589987888e+01, -7.59601722513950107896e+02, -1.18498066702429587167e+04, -4.84385124285750353010e+04}, []float64{1.61395369700722909556e+02, 7.82538599923348465381e+03, 1.33875336287249578163e+05, 7.19657723683240939863e+05, 6.66601232617776375264e+05, -2.94490264303834643215e+05}},
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
	return (0.375 + horner(table.p)/mulAdd(z, horner(table.q), 1, profile != 0)) / x
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
	if observation(math.Inf(1), 2) != "9223372036854775808" {
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
	file, err := os.CreateTemp("", "golean-qone-review-")
	must(err)
	defer os.Remove(file.Name())
	fixtures := bufio.NewWriter(file)
	for i, got := range check.LeanRows("godiff-math-qone", 3*len(inputs), &rows) {
		x, profile := inputs[i/3], i%3
		want := observation(x, profile)
		if got != want {
			log.Fatalf("qone(%x), profile %d: Lean=%s oracle=%s", x, profile, got, want)
		}
		if profile == nativeProfile() {
			fmt.Fprintln(fixtures, math.Float64bits(x), got, want)
		}
	}
	must(fixtures.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_MATH_QONE_CASES", file.Name()))
	const data = "difftest/math/qone/testdata/"
	must(check.TestInPackage("math", "TestGoLeanQone", data+"bridge_test.go", data+"qone_test.go"))
	if differences[0] == 0 || differences[1] == 0 || wrongBand != 4 || lostSign != 1 {
		log.Fatalf("invalid variant strength: fusion=%v wrong band=%d", differences, wrongBand)
	}
	fmt.Printf("qone: %d native and %d model observations; omitting denominator/Horner contraction, swapping coefficient bands, and reversing the infinity-result zero sign failed %d/%d/%d/%d inputs\n",
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
