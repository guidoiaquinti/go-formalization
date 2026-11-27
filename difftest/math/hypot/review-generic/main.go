package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const (
		inf  = uint64(0x7ff0000000000000)
		nan  = inf + 1
		sign = uint64(1) << 63
	)
	bits := math.Float64bits
	cases := []struct {
		x, y, fused, separate uint64
		rejectUnscaled        bool
	}{
		{bits(3), bits(-4), bits(5), bits(5), false},
		{sign, sign, 0, 0, false},
		{inf | sign, nan, inf, inf, true},
		{nan, inf, inf, inf, true},
		{nan, bits(1), nan, nan, false},
		{bits(math.MaxFloat64), 0, bits(math.MaxFloat64), bits(math.MaxFloat64), true},
		{1, 0, 1, 1, true},
		{0x3ff305267911c858, 0x3feaee7002cc1d39, 0x3ff74de5f2bff45e, 0x3ff74de5f2bff45f, false},
		{0x7fe6a09e667f3bcc, 0x7fe6a09e667f3bcc, inf, inf, false},
	}
	encode := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
	}
	var input bytes.Buffer
	var expected []string
	for _, c := range cases {
		fmt.Fprintln(&input, 0, c.x, c.y)
		fmt.Fprintln(&input, 1, c.x, c.y)
		expected = append(expected, encode(c.separate), encode(c.fused))
	}
	for i, row := range check.LeanRows("godiff-math-hypot", len(expected), &input) {
		if row != expected[i] {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/2, i%2, row, expected[i])
		}
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		log.Fatal("could not determine build profile")
	}
	nativeFused := runtime.GOARCH == "arm64"
	for _, setting := range info.Settings {
		if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
			nativeFused = true
		}
	}
	var fixture bytes.Buffer
	var unscaledFailures, exchangedProfileFailures int
	for _, c := range cases {
		x, y := math.Float64frombits(c.x), math.Float64frombits(c.y)
		want, exchanged := c.separate, c.fused
		if nativeFused {
			want, exchanged = exchanged, want
		}
		ideal, exact := encode(want), true
		if c.x == 0x7fe6a09e667f3bcc {
			ideal, exact = encode(bits(math.MaxFloat64)), false
		}
		fmt.Fprintln(&fixture, c.x, c.y, encode(want), ideal, exact)
		if want != exchanged {
			exchangedProfileFailures++
		}
		if c.rejectUnscaled && encode(bits(math.Sqrt(x*x+y*y))) != encode(want) {
			unscaledFailures++
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-hypot-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_HYPOT_CASES", path); err != nil {
		log.Fatal(err)
	}
	const data = "difftest/math/hypot/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanHypot", data+"bridge_test.go", data+"hypot_test.go"); err != nil {
		log.Fatal(err)
	}
	if unscaledFailures != 4 || exchangedProfileFailures != 1 {
		log.Fatalf("unscaled failures=%d, exchanged profile failures=%d", unscaledFailures, exchangedProfileFailures)
	}
	fmt.Printf("Private hypot passed %d fixtures; unscaled formula failed %d selected probes; exchanging profiles failed %d\n",
		len(cases), unscaledFailures, exchangedProfileFailures)
}
