package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const inf = 0x7ff0000000000000
	rows := []struct{ word, want uint64 }{
		{0, 0x3ff0000000000000}, {0x8000000000000000, 0x3ff0000000000000},
		{inf, inf}, {0xfff0000000000000, 0},
		{0x7ff0000000000042, 0x7ff0000000000042}, {0xfff8000000000042, 0xfff8000000000042},
		{math.Float64bits(0x1p-52), 0x3ff0000000000001},
		{math.Float64bits(-745), 1}, {math.Float64bits(-746), 0},
		{math.Float64bits(710), inf},
	}
	const nearHigh = 709.5
	var input bytes.Buffer
	for profile := range 3 {
		for _, c := range rows {
			fmt.Fprintln(&input, profile, c.word)
		}
		fmt.Fprintln(&input, profile, math.Float64bits(nearHigh))
	}
	native, err := os.CreateTemp("", "golean-review-arch-exp-")
	must(err)
	defer os.Remove(native.Name())
	stride := len(rows) + 1
	expectations := check.LeanRows("godiff-math-arch-exp", 3*stride, &input)
	for i, got := range expectations {
		word, want := math.Float64bits(nearHigh), uint64(inf)
		profile, idx := i/stride, i%stride
		exact := true
		if idx < len(rows) {
			word, want = rows[idx].word, rows[idx].want
			if profile == 2 && word == math.Float64bits(-745) {
				want = 0
			}
		} else if profile == 2 {
			want = math.Float64bits(math.Exp(nearHigh))
			exact = false
		}
		value, err := strconv.ParseUint(got, 10, 64)
		must(err)
		if idx == len(rows) && profile == 2 {
			if value >= inf {
				log.Fatal("arm64 scaling returned early infinity")
			}
		} else if value != want {
			log.Fatalf("profile %d word %#x: got %#x want %#x", profile, word, value, want)
		}
		if profile == 2 && runtime.GOARCH == "arm64" {
			_, err := fmt.Fprintln(native, word, value, want, exact)
			must(err)
		}
	}
	must(native.Close())
	if runtime.GOARCH == "arm64" {
		must(os.Setenv("GOLEAN_MATH_ARCH_EXP_CASES", native.Name()))
		const dir = "difftest/math/arch-exp/testdata/"
		must(check.TestInPackage("math", "TestGoLeanArchExp", dir+"bridge_test.go", dir+"arch_exp_test.go"))
	}
	fmt.Printf("archExp passed %d fixed profile rows, with native arm64 rows replayed on arm64\n", len(expectations))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
