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
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/trig"
)

const (
	callerExponent = 1023 + 29
	lastExponent   = 2046
	mantissaMask   = (1 << 52) - 1
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var boundaries []uint64
	for exponent := uint64(0); exponent <= lastExponent; exponent++ {
		for _, fraction := range []uint64{0, 1, 1 << 51, mantissaMask} {
			if word := exponent<<52 | fraction; word != 0 {
				boundaries = append(boundaries, word)
			}
		}
	}
	for _, center := range []int64{1 << 30, 1 << 50} {
		for offset := range 16 {
			x := float64(center+int64(offset)) * (math.Pi / 4)
			boundaries = append(boundaries, binary64.Load(math.Nextafter(x, 0)), binary64.Load(x),
				binary64.Load(math.Nextafter(x, math.Inf(1))))
		}
	}
	threshold := float64(math.Pi / 4)
	boundaries = append(boundaries, binary64.Load(math.Nextafter(threshold, 0)),
		binary64.Load(threshold), binary64.Load(math.Nextafter(threshold, math.Inf(1))))
	seeds := []int64{17, 2718, 65537}
	words := append([]uint64(nil), boundaries...)
	const samplesPerSeed = 20_000
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range samplesPerSeed {
			words = append(words, sample(rng, i))
		}
	}
	var leanInput bytes.Buffer
	for _, word := range words {
		fmt.Fprintln(&leanInput, word)
	}
	halfPi := trig.HalfPi()
	quarterPi := new(big.Float).SetPrec(halfPi.Prec()).SetMantExp(halfPi, -1)
	var fixture bytes.Buffer
	for i, model := range check.LeanRows("godiff-math-trig-reduce", len(words), &leanInput) {
		word := words[i]
		if len(strings.Fields(model)) != 2 {
			return fmt.Errorf("input %#016x: invalid Lean result %q", word, model)
		}
		octant, residual := reference(binary64.Store(word), quarterPi)
		fmt.Fprintln(&fixture, word, model, octant, binary64.Load(residual))
	}
	dir, err := os.MkdirTemp("", "golean-trig-reduce-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_TRIG_REDUCE_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/trig-reduce/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanTrigReduce", data+"bridge_test.go", data+"trig_reduce_test.go"); err != nil {
		return err
	}
	fmt.Printf("trigReduce matched %d boundaries and %d samples with seeds %v.\n",
		len(boundaries), len(words)-len(boundaries), seeds)
	return nil
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		k := (1 << 30) + rng.Int63n(1<<52)
		x := float64(k) * (math.Pi / 4)
		return binary64.Load(x) + uint64(rng.Intn(9)-4)
	case 1:
		exponent := uint64(callerExponent + rng.Intn(lastExponent-callerExponent+1))
		return exponent<<52 | rng.Uint64()&mantissaMask
	case 2:
		k := 1 + rng.Intn(1<<20)
		x := float64(k) * (math.Pi / 4)
		return binary64.Load(x) + uint64(rng.Intn(9)-4)
	default:
		exponent := uint64(rng.Intn(lastExponent + 1))
		word := exponent<<52 | rng.Uint64()&mantissaMask
		return max(1, word)
	}
}

func reference(x float64, quarterPi *big.Float) (uint64, float64) {
	prec := quarterPi.Prec()
	value := new(big.Float).SetPrec(prec).SetFloat64(x)
	quotient := new(big.Float).SetPrec(prec).Quo(value, quarterPi)
	whole, _ := quotient.Int(nil)
	if whole.Bit(0) != 0 {
		whole.Add(whole, big.NewInt(1))
	}
	center := new(big.Float).SetPrec(prec).SetInt(whole)
	center.Mul(center, quarterPi)
	angle, _ := value.Sub(value, center).Float64()
	return whole.And(whole, big.NewInt(7)).Uint64(), angle
}
