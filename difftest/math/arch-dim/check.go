package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct{ a, b uint64 }

const (
	executor = "godiff-math-arch-dim"
	sign     = uint64(1 << 63)
	quiet    = uint64(1 << 51)
	inf      = uint64(0x7ff0000000000000)
	panicRow = "panic:not implemented"
)

func main() {
	switch runtime.GOARCH {
	case "loong64", "riscv64", "s390x":
		fmt.Printf("Architecture %s assembly is outside the model.\n", runtime.GOARCH)
		return
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	native := 2
	switch runtime.GOARCH {
	case "amd64":
		native = 0
	case "arm64":
		native = 1
	}
	words := append(binary64.Boundaries(), 0x7ff0000000000042, 0xfff0000000000042,
		0x7ff80000000000a5, 0xfff80000000000b6)
	var boundaries []pair
	for _, a := range words {
		for _, b := range words {
			boundaries = append(boundaries, pair{a, b})
		}
	}
	var fixture bytes.Buffer
	compare := func(label string, cases []pair) error {
		for profile := range 3 {
			var input bytes.Buffer
			for _, p := range cases {
				fmt.Fprintln(&input, profile, p.a, p.b)
			}
			for i, result := range check.LeanRows(executor, len(cases), &input) {
				p := cases[i]
				want := expected(profile, p)
				if result != want {
					return fmt.Errorf("profile %d (%#016x, %#016x): Lean=%q want %q", profile, p.a, p.b, result, want)
				}
				if profile == native {
					fmt.Fprintf(&fixture, "%d %d %q\n", p.a, p.b, want)
				}
			}
		}
		fmt.Printf("archMax/Min %s: %d Lean rows per profile.\n", label, len(cases))
		return nil
	}
	if err := compare("boundaries", boundaries); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]pair, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-arch-dim-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_ARCH_DIM_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/arch-dim/testdata/"
	return check.TestInPackage("math", "TestGoLeanArchDim", data+"bridge_test.go", data+"arch_dim_test.go")
}

func sample(rng *rand.Rand, i int) pair {
	a, b := rng.Uint64(), rng.Uint64()
	switch i % 7 {
	case 1:
		b = a
	case 2:
		b = a ^ sign
	case 3:
		a = inf | (a & sign)
	case 4:
		b = inf | (b & sign)
	case 5:
		a |= inf | 1
	case 6:
		a |= inf | 1
		b |= inf | 1
	}
	return pair{a, b}
}

func expected(profile int, p pair) string {
	if profile == 2 {
		return panicRow
	}
	a, b := p.a, p.b
	if isNaN(a) || isNaN(b) {
		nan := inf | quiet | 1
		if profile == 1 {
			// Highest-ranked NaN wins. The second input is first in
			// the instruction's operand order and wins equal priority.
			winner := b
			if nanPriority(a) > nanPriority(b) {
				winner = a
			}
			nan = winner | quiet
		}
		maximum, minimum := nan, nan
		if a == inf || b == inf {
			maximum = inf
		}
		if a == sign|inf || b == sign|inf {
			minimum = sign | inf
		}
		return fmt.Sprintf("%d,%d", maximum, minimum)
	}
	if (a|b)&^sign == 0 {
		return fmt.Sprintf("%d,%d", a&b, a|b)
	}
	if binary64.Store(a) < binary64.Store(b) {
		return fmt.Sprintf("%d,%d", b, a)
	}
	return fmt.Sprintf("%d,%d", a, b)
}

func isNaN(word uint64) bool {
	return word&^sign > inf
}

func nanPriority(word uint64) int {
	if !isNaN(word) {
		return 0
	}
	if word&quiet == 0 {
		return 2
	}
	return 1
}
