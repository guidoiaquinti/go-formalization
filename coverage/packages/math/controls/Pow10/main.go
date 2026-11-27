package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"runtime"
	"strings"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	cases := []struct {
		exponent int
		want     uint64
	}{
		{math.MinInt, 0},
		{-324, 0},
		{-323, 2},
		{-307, 0x0031fa182c40c60e},
		{0, 0x3ff0000000000000},
		{33, 0x46c8a6e32246c99d},
		{309, 0x7ff0000000000000},
		{math.MaxInt, 0x7ff0000000000000},
	}
	for _, c := range cases {
		if got := math.Float64bits(math.Pow10(c.exponent)); got != c.want {
			log.Fatalf("Pow10(%d) = %#016x, want %#016x", c.exponent, got, c.want)
		}
	}
	// Round the ideal decimal once, ignoring the stored operands.
	for _, c := range []struct {
		exponent int
		want     uint64
	}{
		{-307, 0x0031fa182c40c60d},
		{33, 0x46c8a6e32246c99c},
	} {
		magnitude := max(c.exponent, -c.exponent)
		exact := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(magnitude)), nil))
		if c.exponent < 0 {
			exact.Inv(exact)
		}
		rounded, _ := exact.Float64()
		if got := math.Float64bits(rounded); got != c.want || got == math.Float64bits(math.Pow10(c.exponent)) {
			log.Fatalf("single-rounding mutant for %d produced %#016x, want mutant word %#016x", c.exponent, got, c.want)
		}
	}
	fmt.Printf("%s %s/%s: Pow10 passed %d cases, single-rounding mutant rejected twice\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases))
}
