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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ x, y uint64 }

func main() {
	if runtime.GOARCH == "386" {
		fmt.Println("The modeled archHypot profiles do not cover 386 assembly.")
		return
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
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
	nativeAMD64 := runtime.GOARCH == "amd64"
	var maxDistance uint64
	var overflowMismatches int
	var fixture bytes.Buffer
	compare := func(label string, cases []input) error {
		const executor = "godiff-math-arch-hypot"
		var amd64, fallback bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&amd64, 0, c.x, c.y)
			fmt.Fprintln(&fallback, 1, c.x, c.y)
		}
		stubs := check.LeanRows(executor, len(cases), &fallback)
		for i, model := range check.LeanRows(executor, len(cases), &amd64) {
			c := cases[i]
			if stubs[i] != "panic:not implemented" {
				return fmt.Errorf("archHypot stub(%#016x, %#016x): Lean=%q", c.x, c.y, stubs[i])
			}
			x, y := binary64.Store(c.x), binary64.Store(c.y)
			want := idealHypot(x, y)
			var got uint64
			if _, err := fmt.Sscan(model, &got); err != nil {
				return err
			}
			expected := binary64.Load(want)
			exact := math.IsInf(x, 0) || math.IsInf(y, 0) || math.IsNaN(x) || math.IsNaN(y) || x == 0 && y == 0
			distance := max(got, expected) - min(got, expected)
			if (exact && distance != 0) || distance > 2 {
				return fmt.Errorf("archHypot(%#016x, %#016x): Lean=%#016x reference=%#016x distance=%d", c.x, c.y, got, expected, distance)
			}
			if math.IsInf(binary64.Store(got), 1) != math.IsInf(want, 1) {
				overflowMismatches++
			}
			maxDistance = max(maxDistance, distance)
			if !nativeAMD64 {
				model = stubs[i]
			}
			fmt.Fprintf(&fixture, "%d %d %q\n", c.x, c.y, model)
		}
		fmt.Printf("archHypot %s: %d modeled rows per profile.\n", label, len(cases))
		return nil
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	fmt.Printf("archHypot amd64 model independent reference: maximum %d steps, %d finite/overflow classification differences.\n", maxDistance, overflowMismatches)
	dir, err := os.MkdirTemp("", "golean-math-arch-hypot-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_ARCH_HYPOT_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/arch-hypot/testdata/"
	return check.TestInPackage("math", "TestGoLeanArchHypot", data+"bridge_test.go", data+"arch_hypot_test.go")
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
		return math.Float64frombits(0x7ff8000000000001)
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
