package main

import (
	"bytes"
	"cmp"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	x, y   string
	shared bool
}

func main() {
	edge := []string{"0", "1", "-1", "9", "-10", "999", "-1000",
		"4294967295", "-4294967296", "18446744073709551615", "-18446744073709551616",
		"340282366920938463463374607431768211455", "-340282366920938463463374607431768211456",
		strings.Repeat("9", 720), "-1" + strings.Repeat("0", 720)}
	var fixed []fixture
	for _, x := range edge {
		for _, y := range edge {
			fixed = append(fixed, fixture{x, y, false})
		}
		fixed = append(fixed, fixture{x, x, true})
	}
	run("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				x, y := sample(rng), sample(rng)
				share := i%11 == 0
				switch {
				case share || i%11 == 1:
					y = x
				case i%11 == 2:
					y = flip(x)
				}
				rows[i] = fixture{x, y, share}
			}
			run(fmt.Sprintf("seed=%d block=%d", seed, block), rows)
		}
	}
}

func sample(rng *rand.Rand) string {
	digits := make([]byte, rng.Intn(65)+1)
	for i := range digits {
		digits[i] = byte(rng.Intn(10)) + '0'
	}
	value := normalized(string(digits))
	if rng.Intn(2) == 0 {
		return flip(value)
	}
	return value
}

func run(label string, rows []fixture) {
	var input bytes.Buffer
	for _, row := range rows {
		fmt.Fprintln(&input, row.x, row.y)
	}
	for i, got := range check.LeanRows("godiff-math-big-int-arithmetic", len(rows), &input) {
		f := rows[i]
		want := oracle(f.x, f.y)
		if got != want {
			log.Fatalf("%s row %d: Lean %q want %q", label, i, got, want)
		}
		for mode := range 3 {
			if native := observe(f, mode); native != want {
				log.Fatalf("%s row %d alias=%d: Go %q want %q", label, i, mode, native, want)
			}
		}
	}
	fmt.Printf("%s: %d pinned Go and Lean cases under each receiver alias profile\n", label, len(rows))
}

func observe(f fixture, mode int) string {
	x, y := inputs(f)
	cmp, cmpAbs := x.Cmp(y), x.CmpAbs(y)
	unary := []func(*big.Int, *big.Int) *big.Int{(*big.Int).Set, (*big.Int).Abs, (*big.Int).Neg}
	binary := []func(*big.Int, *big.Int, *big.Int) *big.Int{(*big.Int).Add, (*big.Int).Sub, (*big.Int).Mul}
	var out []string
	for _, operation := range unary {
		a, _ := inputs(f)
		dst := big.NewInt(-117)
		if mode == 1 {
			dst = a
		}
		if operation(dst, a) != dst || (mode != 1 && a.String() != f.x) {
			log.Fatal("unary operation changed input or returned a different receiver")
		}
		out = append(out, dst.String())
	}
	for _, operation := range binary {
		a, b := inputs(f)
		dst := big.NewInt(-117)
		switch mode {
		case 1:
			dst = a
		case 2:
			dst = b
		}
		if operation(dst, a, b) != dst || a != dst && a.String() != f.x || b != dst && b.String() != f.y {
			log.Fatal("binary operation changed input or returned a different receiver")
		}
		out = append(out, dst.String())
	}
	if x.String() != f.x || y.String() != f.y {
		log.Fatal("comparison mutated an operand")
	}
	return fmt.Sprintf("%s %d %d", strings.Join(out, " "), cmp, cmpAbs)
}

func inputs(f fixture) (x, y *big.Int) {
	x, okX := new(big.Int).SetString(f.x, 10)
	y, okY := new(big.Int).SetString(f.y, 10)
	if !okX || !okY {
		log.Fatal("invalid decimal input")
	}
	if f.shared {
		y = x
	}
	return
}

func oracle(x, y string) string {
	a, negA := magnitude(x)
	b, negB := magnitude(y)
	absoluteOrder := compareDigits(a, b)
	order := absoluteOrder
	if negA != negB {
		if negA {
			order = -1
		} else {
			order = 1
		}
	} else if negA {
		order = -absoluteOrder
	}
	product := withSign(multiply(a, b), negA != negB)
	return fmt.Sprintf("%s %s %s %s %s %s %d %d", x, a, flip(x),
		sum(x, y), sum(x, flip(y)), product, order, absoluteOrder)
}

func magnitude(value string) (string, bool) {
	return strings.CutPrefix(value, "-")
}

func normalized(value string) string {
	value = strings.TrimLeft(value, "0")
	if value == "" {
		return "0"
	}
	return value
}

func withSign(value string, negative bool) string {
	if negative && value != "0" {
		return "-" + value
	}
	return value
}

func flip(value string) string {
	digits, neg := magnitude(value)
	return withSign(digits, !neg)
}

func compareDigits(a, b string) int {
	if order := cmp.Compare(len(a), len(b)); order != 0 {
		return order
	}
	return strings.Compare(a, b)
}

func sum(x, y string) string {
	a, negA := magnitude(x)
	b, negB := magnitude(y)
	if negA == negB {
		return withSign(columnSum(a, b, false), negA)
	}
	if compareDigits(a, b) < 0 {
		return withSign(columnSum(b, a, true), negB)
	}
	return withSign(columnSum(a, b, true), negA)
}

// columnSum subtracts only when a is at least b. In subtraction, the
// previous column's carry is a borrow represented by -1.
func columnSum(a, b string, subtract bool) string {
	n := max(len(a), len(b))
	digits := make([]byte, n+1)
	carry := 0
	for pos := range n {
		value := carry + digitFromRight(a, pos)
		right := digitFromRight(b, pos)
		if subtract {
			value -= right
		} else {
			value += right
		}
		if value < 0 {
			value, carry = value+10, -1
		} else {
			value, carry = value%10, value/10
		}
		digits[n-pos] = byte(value) + '0'
	}
	digits[0] = byte(carry) + '0'
	return normalized(string(digits))
}

func digitFromRight(s string, pos int) int {
	if pos >= len(s) {
		return 0
	}
	return int(s[len(s)-1-pos] - '0')
}

func multiply(a, b string) string {
	columns := make([]int, len(a)+len(b))
	for i := range len(a) {
		for j := range len(b) {
			columns[i+j] += digitFromRight(a, i) * digitFromRight(b, j)
		}
	}
	out := make([]byte, len(columns))
	carry := 0
	for i, n := range columns {
		n += carry
		out[len(out)-1-i] = byte(n%10) + '0'
		carry = n / 10
	}
	return normalized(string(out))
}
