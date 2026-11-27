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
		{0x1p-27, 0x1p-27},
		// Finite golden values were obtained independently with platform libm.
		{0.25, 0x1.f5b75f92c80ddp-3},
		{0.6, 0x1.14b1dd5f90ce1p-1},
		{0x1.51eb851eb851ep-1, 0x1.2aafdde4d0c9ep-1},
		{0x1.51eb851eb851fp-1, 0x1.2aafdde4d0c9fp-1},
		{0x1.51eb851eb8520p-1, 0x1.2aafdde4d0ca0p-1},
		{0.75, 0x1.4978fa3269ee1p-1},
		{1, math.Pi / 4},
		{0x1.3504f333f9de5p+1, 0x1.2d97c7f3321d2p+0},
		{0x1.3504f333f9de6p+1, 0x1.2d97c7f3321d2p+0},
		{0x1.3504f333f9de7p+1, 0x1.2d97c7f3321d2p+0},
		{3, 0x1.3fc176b7a8560p+0},
		{8, 0x1.7249faa996a21p+0},
		{32, 0x1.8a205fd558740p+0},
		{math.MaxFloat64, math.Pi / 2}, {math.Inf(1), math.Pi / 2},
	}
	const profiles = 2
	var inputs, cases bytes.Buffer
	for _, c := range tests {
		for profile := range profiles {
			fmt.Fprintln(&inputs, profile, math.Float64bits(c.input))
		}
	}
	rows := check.LeanRows("godiff-math-satan", profiles*len(tests), &inputs)
	native := nativeProfile()
	var missingMiddle, missingLarge int
	for i, c := range tests {
		want := math.Float64bits(c.want)
		for profile := range profiles {
			got, err := strconv.ParseUint(rows[profiles*i+profile], 10, 64)
			must(err)
			if distance(got, want) > 2 || (math.Abs(c.input) <= 0x1p-27 || math.IsInf(c.input, 1)) && got != want {
				log.Fatalf("Lean profile %d satan(%g) = %#016x, want %#016x", profile, c.input, got, want)
			}
		}
		fmt.Fprintln(&cases, math.Float64bits(c.input), rows[profiles*i+native], want)
		if c.input > 0.66 && c.input <= 2.414213562373095 {
			reduced := math.Atan((c.input - 1) / (c.input + 1))
			if math.Abs(reduced-c.want) > 0x1p-51 {
				missingMiddle++
			}
		}
		if c.input > 2.414213562373095 && math.Abs(-math.Atan(1/c.input)-c.want) > 0x1p-51 {
			missingLarge++
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-satan-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_SATAN_CASES", path))
	const overlay = "difftest/math/satan/testdata/"
	must(check.TestInPackage("math", "TestGoLeanSatan", overlay+"bridge_test.go", overlay+"satan_test.go"))
	if missingMiddle == 0 || missingLarge == 0 {
		log.Fatal("independent fixtures did not distinguish omitted middle- and large-argument offsets")
	}
	fmt.Printf("Private satan passed %d native and %d Lean fixtures; missing middle offset failed %d and missing large offset failed %d\n",
		len(tests), len(rows), missingMiddle, missingLarge)
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
			log.Fatal("cannot determine amd64 satan profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled satan architecture %s", runtime.GOARCH)
		return 0
	}
}
