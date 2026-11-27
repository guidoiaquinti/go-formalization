package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixed := []struct{ row, want string }{
		{"0 -1 0 0 1", "0|0|-1|-1|-1|0|0|0|1|0"},
		{"-3 5 1 0 0", "5|-8|-3|-8|2|-6|-2|1|-4|0"},
		{"-5 -9 3 2 1", "-13|8|-1|12|4|-40|-1|0|-1|0"},
		{"-256 255 8 15 0", "0|-256|-1|-1|255|-65536|-1|1|-33024|8"},
		{"18446744073709551616 -36893488147419103233 64 64 0", "18446744073709551616|0|-36893488147419103233|-55340232221128654849|-18446744073709551617|340282366920938463463374607431768211456|1|1|0|64"},
		{"-680564733841876926926749214863536422912 36893488147419103232 3 130 0", "0|-680564733841876926926749214863536422912|-680564733841876926889855726716117319680|-680564733841876926889855726716117319680|680564733841876926926749214863536422911|-5444517870735015415413993718908291383296|-85070591730234615865843651857942052864|1|-2041694201525630780780247644590609268736|129"},
		{"-1 0 4 130 0", "0|-1|-1|-1|0|-16|-1|1|-1361129467683753853853498429727072845825|0"},
		{"1 -1 2 -1 2", "1|0|-1|-2|-2|4|0|negative bit index|negative bit index|0"},
		{"1 -1 2 3 2", "1|0|-1|-2|-2|4|0|0|set bit is not 0 or 1|0"},
		{"-5 -5 0 2 0", "-5|0|-5|0|4|-5|-5|0|-5|0"},
	}
	var input bytes.Buffer
	for _, c := range fixed {
		fmt.Fprintln(&input, c.row)
	}
	shiftMutation := 0
	for i, got := range check.LeanRows("godiff-math-big-int-bitwise", len(fixed), &input) {
		c := fixed[i]
		if got != c.want {
			log.Fatalf("row %d: Lean %q, expected %q", i, got, c.want)
		}
		checkLaws(c.row)
		var x, y big.Int
		var shift uint
		if _, err := fmt.Sscan(c.row, &x, &y, &shift); err != nil {
			log.Fatal(err)
		}
		divisor := new(big.Int).Lsh(big.NewInt(1), shift)
		if x.Quo(&x, divisor).String() != strings.Split(c.want, "|")[6] {
			shiftMutation++
		}
		for alias := range 3 {
			if nativeResult := native(c.row, alias); nativeResult != c.want {
				log.Fatalf("row %d alias %d: Go %q, expected %q", i, alias, nativeResult, c.want)
			}
		}
	}
	if shiftMutation != 3 {
		log.Fatalf("truncate-to-zero shift mutant disagrees on %d rows, want 3", shiftMutation)
	}
	fmt.Printf("%d literal Go/Lean rows with three native receiver profiles; truncated shift disagrees on %d\n", len(fixed), shiftMutation)
}

func checkLaws(row string) {
	var x big.Int
	if _, err := fmt.Sscan(row, &x); err != nil {
		log.Fatal(err)
	}
	var z, neg big.Int
	if z.And(&x, &x).Cmp(&x) != 0 || z.Or(&x, &x).Cmp(&x) != 0 ||
		z.Xor(&x, &x).Sign() != 0 || z.AndNot(&x, &x).Sign() != 0 ||
		z.Not(z.Not(&x)).Cmp(&x) != 0 || z.Lsh(&x, 0).Cmp(&x) != 0 ||
		z.Rsh(&x, 0).Cmp(&x) != 0 ||
		neg.Neg(&x).TrailingZeroBits() != x.TrailingZeroBits() {
		log.Fatal("a registered algebraic law failed on pinned Go")
	}
}

func native(row string, alias int) string {
	var a, b big.Int
	var shift, bit uint
	var index int
	if _, err := fmt.Sscan(row, &a, &b, &shift, &index, &bit); err != nil {
		log.Fatal(err)
	}
	apply := func(op func(z, x, y *big.Int) *big.Int) string {
		x, y := new(big.Int).Set(&a), new(big.Int).Set(&b)
		if a.Cmp(&b) == 0 {
			y = x
		}
		z := []*big.Int{big.NewInt(-245), x, y}[alias]
		return observe(func() string {
			if op(z, x, y) != z || (x != z && x.Cmp(&a) != 0) || (y != z && y.Cmp(&b) != 0) {
				log.Fatal("receiver identity or distinct input preservation failed")
			}
			return z.String()
		})
	}
	return strings.Join([]string{
		apply((*big.Int).And), apply((*big.Int).AndNot), apply((*big.Int).Or), apply((*big.Int).Xor),
		apply(func(z, x, y *big.Int) *big.Int { return z.Not(x) }),
		apply(func(z, x, y *big.Int) *big.Int { return z.Lsh(x, shift) }),
		apply(func(z, x, y *big.Int) *big.Int { return z.Rsh(x, shift) }),
		observe(func() string { return strconv.Itoa(int(a.Bit(index))) }),
		apply(func(z, x, y *big.Int) *big.Int { return z.SetBit(x, index, bit) }),
		strconv.Itoa(int(a.TrailingZeroBits())),
	}, "|")
}

func observe(run func() string) (out string) {
	defer func() {
		if p := recover(); p != nil {
			switch p {
			case "negative bit index", "set bit is not 0 or 1":
				out = p.(string)
			default:
				panic(p)
			}
		}
	}()
	return run()
}
