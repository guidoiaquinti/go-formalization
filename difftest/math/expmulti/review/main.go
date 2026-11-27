package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	hi, lo float64
	k      int
}

func main() {
	native := nativeProfile()
	var cases []input
	for k := -1075; k <= 1024; k++ {
		cases = append(cases, input{0, 0, k})
	}
	for _, k := range []int{-1075, -1074, -1023, -1022, -1, 0, 1, 1023, 1024} {
		for _, a := range []float64{-0.34699, -0.2, -0x1p-54, 0x1p-54, 0.2, 0.34699} {
			for _, b := range []float64{-2.1e-7, 0, 2.1e-7} {
				cases = append(cases, input{a, b, k})
			}
		}
	}
	var requests bytes.Buffer
	for _, c := range cases {
		for profile := range 3 {
			fmt.Fprintln(&requests, profile, math.Float64bits(c.hi), math.Float64bits(c.lo), c.k)
		}
	}
	rows := check.LeanRows("godiff-math-expmulti", 3*len(cases), &requests)
	fixture, err := os.CreateTemp("", "golean-review-expmulti-")
	must(err)
	defer os.Remove(fixture.Name())
	fixtureRows := bufio.NewWriter(fixture)
	var signedExponent, ignoredLow, flushedUnderflow int
	for i, c := range cases {
		// Public Exp has its own reduction on the tested architectures.
		// Binary scaling uses public Ldexp on a positive finite mantissa.
		scaled := math.Exp(c.hi - c.lo)
		want := math.Float64bits(math.Ldexp(scaled, c.k))
		for profile := range 3 {
			row := rows[3*i+profile]
			got, err := strconv.ParseUint(row, 10, 64)
			must(err)
			distance := max(got, want) - min(got, want)
			if distance > 3 || (c.hi == 0 && c.lo == 0 && distance != 0) {
				log.Fatalf("case %d profile %d: got %#x, want %#x", i, profile, got, want)
			}
			if profile == native {
				_, err := fmt.Fprintln(fixtureRows, math.Float64bits(c.hi), math.Float64bits(c.lo), c.k, got, want)
				must(err)
			}
		}
		if c.k < 0 && math.Float64bits(math.Ldexp(scaled, -c.k)) != want {
			signedExponent++
		}
		omit := math.Float64bits(math.Ldexp(math.Exp(c.hi), c.k))
		if max(want, omit)-min(want, omit) > 3 {
			ignoredLow++
		}
		if 0 < want && want < 1<<52 {
			flushedUnderflow++
		}
	}
	must(fixtureRows.Flush())
	must(fixture.Close())
	must(os.Setenv("GOLEAN_MATH_EXPMULTI_CASES", fixture.Name()))
	const data = "difftest/math/expmulti/testdata/"
	must(check.TestInPackage("math", "TestGoLeanExpmulti", data+"bridge_test.go", data+"expmulti_test.go"))
	if signedExponent < 1075 || ignoredLow < 40 || flushedUnderflow < 52 {
		log.Fatalf("insufficient variant detection: signed=%d low=%d underflow=%d", signedExponent, ignoredLow, flushedUnderflow)
	}
	fmt.Printf("expmulti: %d native and %d Lean observations, incorrect signs/ignored low/flushed subnormals caught on %d/%d/%d rows\n",
		len(cases), len(rows), signedExponent, ignoredLow, flushedUnderflow)
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if ok {
			for _, setting := range info.Settings {
				if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
					return 1
				}
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled architecture %s", runtime.GOARCH)
		return -1
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
