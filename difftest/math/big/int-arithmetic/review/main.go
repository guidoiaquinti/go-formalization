package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	x, y   string
	shared bool
	want   string
}

func main() {
	cases := []fixture{
		{"-7", "3", false, "-7 7 7 -4 -10 -21 -1 1"},
		{"-9", "-9", true, "-9 9 9 -18 0 81 0 0"},
		{"18446744073709551616", "-1", false, "18446744073709551616 18446744073709551616 -18446744073709551616 18446744073709551615 18446744073709551617 -18446744073709551616 1 1"},
		{"0", "0", true, "0 0 0 0 0 0 0 0"},
		{"-7", "7", false, "-7 7 7 0 -14 -49 -1 0"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.x, c.y)
		want := strings.Fields(c.want)
		for column, run := range []func(*big.Int, *big.Int) *big.Int{(*big.Int).Set, (*big.Int).Abs, (*big.Int).Neg} {
			for mode := range 2 {
				x, _ := operands(c)
				dst := big.NewInt(-999)
				if mode == 1 {
					dst = x
				}
				if result := run(dst, x); result != dst || dst.String() != want[column] ||
					x != dst && x.String() != c.x {
					log.Fatalf("%+v unary column %d mode %d did not preserve result/identity", c, column, mode)
				}
			}
		}
		for column, run := range []func(*big.Int, *big.Int, *big.Int) *big.Int{(*big.Int).Add, (*big.Int).Sub, (*big.Int).Mul} {
			for mode := range 3 {
				x, y := operands(c)
				dst := []*big.Int{big.NewInt(-999), x, y}[mode]
				if result := run(dst, x, y); result != dst || dst.String() != want[column+3] ||
					x != dst && x.String() != c.x || y != dst && y.String() != c.y {
					log.Fatalf("%+v binary column %d mode %d did not preserve result/identity", c, column, mode)
				}
			}
		}
		x, y := operands(c)
		if fmt.Sprint(x.Cmp(y)) != want[6] || fmt.Sprint(x.CmpAbs(y)) != want[7] ||
			x.String() != c.x || y.String() != c.y {
			log.Fatalf("%+v wrong comparison or mutated operand", c)
		}
	}
	rows := check.LeanRows("godiff-math-big-int-arithmetic", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].want {
			log.Fatalf("Lean row %d got %q, want %q", i, row, cases[i].want)
		}
	}
	for _, m := range []struct {
		row, col int
		value    string
	}{
		{0, 0, "7"},   // discard copied sign
		{0, 1, "-7"},  // retain sign in Abs
		{3, 2, "-0"},  // manufacture negative zero
		{0, 3, "-10"}, // add magnitudes for mixed signs
		{0, 4, "10"},  // reverse subtraction operands
		{1, 5, "-81"}, // keep negative sign for a square
		{2, 5, "0"},   // truncate product to a machine word
		{0, 6, "1"},   // compare signed operands by magnitude
		{4, 7, "-1"},  // compare absolute operands by sign
	} {
		if strings.Fields(rows[m.row])[m.col] == m.value {
			log.Fatalf("row %d column %d accepted mutant %q", m.row, m.col, m.value)
		}
	}
	log.Print("five independent native and Lean cases passed and rejected nine wrong observations")
}

func operands(c fixture) (*big.Int, *big.Int) {
	x, ok := new(big.Int).SetString(c.x, 10)
	if !ok {
		log.Fatalf("invalid fixture %q", c.x)
	}
	if c.shared {
		return x, x
	}
	y, ok := new(big.Int).SetString(c.y, 10)
	if !ok {
		log.Fatalf("invalid fixture %q", c.y)
	}
	return x, y
}
