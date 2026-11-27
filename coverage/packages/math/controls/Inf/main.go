package main

import (
	"fmt"
	"log"
	"math"
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
	signs := []int{math.MinInt, -2, -1, 0, 1, 2, math.MaxInt}
	var real, zeroNegative, finiteExponent int
	for _, sign := range signs {
		bits := math.Float64bits(math.Inf(sign))
		if !valid(bits, sign) {
			real++
		}
		wrongZero := bits
		if sign == 0 {
			wrongZero ^= 1 << 63
		}
		if !valid(wrongZero, sign) {
			zeroNegative++
		}
		if !valid(math.Float64bits(math.Inf(sign))&^(1<<52), sign) {
			finiteExponent++
		}
	}
	if real != 0 || zeroNegative != 1 || finiteExponent != len(signs) {
		log.Fatalf("Inf failures: Go=%d, zero input negative=%d, finite exponent=%d", real, zeroNegative, finiteExponent)
	}
	fmt.Printf("%s %s/%s: Inf Go=%d/%d, zero input negative mutant=%d/%d, finite exponent mutant=%d/%d\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, len(signs), zeroNegative, len(signs),
		finiteExponent, len(signs))
}

func valid(bits uint64, sign int) bool {
	return (bits>>63 != 0) == (sign < 0) && (bits>>52)&0x7ff == 0x7ff && bits&((1<<52)-1) == 0
}
