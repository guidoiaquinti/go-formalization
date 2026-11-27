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

type fixture struct {
	input, want float64
}

func main() {
	tests := []fixture{
		{0, 0}, {math.Copysign(0, -1), 0},
		{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64},
		{-math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64},
		{0x1p-27, 0x1p-27}, {-0x1p-27, -0x1p-27},
		// Nontrivial values were obtained independently with platform libm.
		{0.03125, 0x1.ffd55bba97625p-6}, {-0.03125, -0x1.ffd55bba97625p-6},
		{0.125, 0x1.fd5ba9aac2f6ep-4}, {-0.125, -0x1.fd5ba9aac2f6ep-4},
		{0.25, 0x1.f5b75f92c80ddp-3}, {-0.25, -0x1.f5b75f92c80ddp-3},
		{1.0 / 3, 0x1.4978fa3269ee1p-2}, {-1.0 / 3, -0x1.4978fa3269ee1p-2},
		{0.5, 0x1.dac670561bb50p-2}, {-0.5, -0x1.dac670561bb50p-2},
		{0.6, 0x1.14b1dd5f90ce1p-1}, {-0.6, -0x1.14b1dd5f90ce1p-1},
		{0x1.51eb851eb851ep-1, 0x1.2aafdde4d0c9ep-1},
		{-0x1.51eb851eb851ep-1, -0x1.2aafdde4d0c9ep-1},
		{0.66, 0x1.2aafdde4d0c9fp-1}, {-0.66, -0x1.2aafdde4d0c9fp-1},
	}
	const profiles = 2
	var inputs, cases bytes.Buffer
	for _, c := range tests {
		for profile := range profiles {
			fmt.Fprintln(&inputs, profile, math.Float64bits(c.input))
		}
	}
	rows := check.LeanRows("godiff-math-xatan", profiles*len(tests), &inputs)
	native := nativeProfile()
	var linear, absolute int
	for i, c := range tests {
		want := math.Float64bits(c.want)
		for profile := range profiles {
			got, err := strconv.ParseUint(rows[profiles*i+profile], 10, 64)
			must(err)
			if distance(got, want) > 2 || math.Abs(c.input) <= 0x1p-27 && got != want {
				log.Fatalf("Lean profile %d xatan(%g) = %#016x, want %#016x", profile, c.input, got, want)
			}
		}
		fmt.Fprintln(&cases, math.Float64bits(c.input), rows[profiles*i+native], want)
		if distance(math.Float64bits(c.input), want) > 2 {
			linear++
		}
		if distance(math.Float64bits(math.Abs(c.want)), want) > 2 {
			absolute++
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-xatan-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_XATAN_CASES", path))
	const overlay = "difftest/math/xatan/testdata/"
	must(check.TestInPackage("math", "TestGoLeanXatan", overlay+"bridge_test.go", overlay+"xatan_test.go"))
	if linear == 0 || absolute == 0 {
		log.Fatal("independent fixtures did not distinguish linear and unsigned approximations")
	}
	fmt.Printf("Private xatan passed %d native and %d Lean fixtures; linear failed %d and unsigned failed %d\n",
		len(tests), len(rows), linear, absolute)
}

func distance(x, y uint64) uint64 {
	return max(x, y) - min(x, y)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 1
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine amd64 xatan profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled xatan architecture %s", runtime.GOARCH)
		return 0
	}
}
