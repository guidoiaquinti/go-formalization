package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ x, y uint64 }

type oracleRow struct {
	input
	ideal string
	exact bool
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var fused int
	if runtime.GOARCH == "arm64" {
		fused = 1
	}
	if runtime.GOARCH != "arm64" && runtime.GOARCH != "amd64" {
		return fmt.Errorf("unmodeled Hypot architecture %s", runtime.GOARCH)
	}
	words := binary64.Boundaries()
	for _, value := range []float64{
		2 * math.SmallestNonzeroFloat64, 3 * math.SmallestNonzeroFloat64,
		0x1p-1022, 0x1p-27, 1, 3, 4, 5,
		math.MaxFloat64 / math.Sqrt2, math.Nextafter(math.MaxFloat64/math.Sqrt2, math.Inf(1)),
	} {
		words = append(words, math.Float64bits(value), math.Float64bits(-value))
	}
	var edges []input
	for _, x := range words {
		for _, y := range words {
			edges = append(edges, input{x, y})
		}
	}
	var maxDistance, overflowMismatches uint64
	var cases []oracleRow
	check.CompareSamples("godiff-math-hypot", edges, sample, func(c input) string {
		return fmt.Sprintf("%d %d %d", fused, c.x, c.y)
	}, func(c input) string {
		x, y := binary64.Store(c.x), binary64.Store(c.y)
		got := math.Hypot(x, y)
		if swapped := math.Hypot(y, x); observe(got) != observe(swapped) {
			log.Fatalf("asymmetric Hypot for %#016x and %#016x", c.x, c.y)
		}
		return observe(got)
	}, func(c input) string {
		x, y := binary64.Store(c.x), binary64.Store(c.y)
		got := math.Hypot(x, y)
		ideal := idealHypot(x, y)
		exact := math.IsInf(x, 0) || math.IsInf(y, 0) || math.IsNaN(x) || math.IsNaN(y) || x == 0 && y == 0
		row := oracleRow{input: c, ideal: observe(ideal), exact: exact}
		cases = append(cases, row)
		if exact {
			return row.ideal
		}
		a, b := binary64.Load(got), binary64.Load(ideal)
		if math.IsInf(got, 1) != math.IsInf(ideal, 1) {
			overflowMismatches++
			fmt.Printf("Hypot near overflow (%#016x, %#016x): Go=%#016x oracle=%#016x\n", c.x, c.y, a, b)
		}
		distance := max(a, b) - min(a, b)
		if distance > 2 {
			log.Fatalf("Hypot(%#016x, %#016x) bits %#016x differ from once-rounded %#016x by %d steps", c.x, c.y, a, b, distance)
		}
		maxDistance = max(maxDistance, distance)
		return observe(got)
	})
	fmt.Printf("Hypot once-rounded oracle: maximum distance %d word steps; %d finite/overflow classification differences.\n", maxDistance, overflowMismatches)
	return compareGeneric(cases)
}

func compareGeneric(cases []oracleRow) error {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("cannot determine generic Hypot build profile")
	}
	fused := 0
	if runtime.GOARCH == "arm64" {
		fused = 1
	}
	for _, setting := range info.Settings {
		if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
			fused = 1
		}
	}

	var modelInput, fixture bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&modelInput, fused, c.x, c.y)
	}
	for i, actual := range check.LeanRows("godiff-math-hypot", len(cases), &modelInput) {
		c := cases[i]
		fmt.Fprintln(&fixture, c.x, c.y, actual, c.ideal, c.exact)
	}
	dir, err := os.MkdirTemp("", "golean-math-hypot-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_HYPOT_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/hypot/testdata/"
	return check.TestInPackage("math", "TestGoLeanHypot", fixtures+"bridge_test.go", fixtures+"hypot_test.go")
}

func sample(r *rand.Rand, i int) input {
	x, y := r.Uint64(), r.Uint64()
	switch i % 5 {
	case 0:
		x = x&^(uint64(2047)<<52) | y&(uint64(2047)<<52)
	case 1:
		x &= 1<<52 - 1
		y &= 1<<52 - 1
	case 2:
		y = x
	case 3:
		x = x&^(uint64(2047)<<52) | uint64(2046)<<52
		y = y&^(uint64(2047)<<52) | uint64(2046)<<52
	}
	return input{x, y}
}

// Compute the nearest binary64 hypotenuse directly. Compare squares of
// exact integer unit counts with the squared midpoint of adjacent floats.
func idealHypot(x, y float64) float64 {
	switch {
	case math.IsInf(x, 0) || math.IsInf(y, 0):
		return math.Inf(1)
	case math.IsNaN(x) || math.IsNaN(y):
		return math.NaN()
	}
	units := func(value float64) *big.Int {
		f := big.NewFloat(math.Abs(value))
		n, _ := f.SetMantExp(f, 1074).Int(nil)
		return n
	}
	a, b := units(x), units(y)
	square := new(big.Int).Add(a.Mul(a, a), b.Mul(b, b))
	root := new(big.Int).Sqrt(square)
	shift := uint(max(0, root.BitLen()-53))
	significand := new(big.Int).Rsh(root, shift)
	midpoint := new(big.Int).Lsh(significand, 1)
	midpoint.Add(midpoint, big.NewInt(1)).Lsh(midpoint, shift)
	midpoint.Mul(midpoint, midpoint)
	twiceRootSquared := new(big.Int).Lsh(square, 2)
	if order := twiceRootSquared.Cmp(midpoint); order > 0 || order == 0 && significand.Bit(0) == 1 {
		significand.Add(significand, big.NewInt(1))
	}
	rounded := new(big.Float).SetInt(significand)
	result, _ := rounded.SetMantExp(rounded, int(shift)-1074).Float64()
	return result
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}
