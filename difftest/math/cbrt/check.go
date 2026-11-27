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

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if runtime.GOARCH != "arm64" && runtime.GOARCH != "amd64" {
		return fmt.Errorf("unmodeled Cbrt architecture %s", runtime.GOARCH)
	}
	fused := 0
	if runtime.GOARCH == "arm64" {
		fused = 1
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("cannot determine Cbrt build profile")
	}
	for _, setting := range info.Settings {
		if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
			fused = 1
		}
	}
	edges := binary64.Boundaries()
	for _, value := range []float64{0x1p-1022, 0x1p-1018, 0.125, 0.5, 1, 2, 4, 8, 27, 64, 125} {
		for _, x := range []float64{math.Nextafter(value, 0), value, math.Nextafter(value, math.Inf(1))} {
			edges = append(edges, math.Float64bits(x), math.Float64bits(-x))
		}
	}

	var maxDistance uint64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-cbrt", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", fused, word)
	}, func(word uint64) string {
		return strconv.FormatUint(binary64.Load(math.Cbrt(binary64.Store(word))), 10)
	}, func(word uint64) string {
		x := binary64.Store(word)
		if x == 0 || math.IsNaN(x) || math.IsInf(x, 0) {
			fmt.Fprintln(&fixture, word, word)
			return strconv.FormatUint(word, 10)
		}
		got := binary64.Load(math.Cbrt(x))
		want := binary64.Load(idealCbrt(x))
		if opposite := binary64.Load(math.Cbrt(-x)); opposite != got^(1<<63) {
			log.Fatalf("Cbrt is not odd for input %#016x", word)
		}
		distance := max(got, want) - min(got, want)
		if distance > 1 {
			log.Fatalf("Cbrt(%#016x) bits %#016x differ from once-rounded %#016x by %d steps", word, got, want, distance)
		}
		maxDistance = max(maxDistance, distance)
		// On these profiles the public entry point selects the generic helper.
		// The fixture result also exact-matches this row's Lean result.
		fmt.Fprintln(&fixture, word, got)
		return strconv.FormatUint(got, 10)
	})
	fmt.Printf("Cbrt once-rounded oracle: maximum encoded-word distance %d.\n", maxDistance)

	dir, err := os.MkdirTemp("", "golean-math-cbrt-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_CBRT_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/cbrt/testdata/"
	return check.TestInPackage("math", "TestGoLeanCbrt", fixtures+"bridge_test.go", fixtures+"cbrt_test.go")
}

func sample(r *rand.Rand, i int) uint64 {
	word := r.Uint64()
	switch i % 4 {
	case 0:
		return word &^ (uint64(2047) << 52)
	case 1:
		exponent := r.Intn(699) - 350
		cube := math.Ldexp(1, 3*exponent)
		return word>>63<<63 | uint64(int64(math.Float64bits(cube))+int64(r.Intn(5)-2))
	case 2:
		return word&((1<<52)-1) | uint64(2043+r.Intn(4))<<52 | word>>63<<63
	default:
		return word
	}
}

// Compute the closest binary64 cube root by comparing the exact input
// with the cube of the midpoint between the two neighboring roots.
func idealCbrt(x float64) float64 {
	f := big.NewFloat(math.Abs(x))
	units, _ := f.SetMantExp(f, 1074).Int(nil)
	target := units.Lsh(units, 2148)
	root := integerCubeRoot(target)
	shift := uint(max(0, root.BitLen()-53))
	significand := new(big.Int).Rsh(root, shift)
	mid := new(big.Int).Lsh(significand, 1)
	mid.Add(mid, big.NewInt(1)).Lsh(mid, shift)
	cubed := new(big.Int).Mul(mid, mid)
	cubed.Mul(cubed, mid)
	scaled := new(big.Int).Lsh(target, 3)
	if order := scaled.Cmp(cubed); order > 0 || order == 0 && significand.Bit(0) == 1 {
		significand.Add(significand, big.NewInt(1))
	}
	rounded := new(big.Float).SetInt(significand)
	result, _ := rounded.SetMantExp(rounded, int(shift)-1074).Float64()
	return math.Copysign(result, x)
}

func integerCubeRoot(target *big.Int) *big.Int {
	root := new(big.Int).Lsh(big.NewInt(1), uint((target.BitLen()+2)/3))
	three := big.NewInt(3)
	var square, next big.Int
	for {
		square.Mul(root, root)
		next.Quo(target, &square)
		next.Add(&next, square.Lsh(root, 1)).Quo(&next, three)
		if next.Cmp(root) >= 0 {
			return root
		}
		root.Set(&next)
	}
}
