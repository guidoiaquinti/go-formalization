package main

import (
	"bytes"
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

func operand(text string) *big.Int {
	n, ok := new(big.Int).SetString(text, 10)
	if !ok {
		log.Fatalf("invalid integer %q", text)
	}
	return n
}

func oracle(f fixture) string {
	x, y := operand(f.x), operand(f.y)
	if y.Sign() == 0 {
		return "panic|panic|panic|panic|panic|panic"
	}
	dividend, divisor := new(big.Int).Abs(x), new(big.Int).Abs(y)
	q, r := new(big.Int), new(big.Int)
	for bit := dividend.BitLen() - 1; bit >= 0; bit-- {
		r.Lsh(r, 1).SetBit(r, 0, dividend.Bit(bit))
		if r.Cmp(divisor) >= 0 {
			r.Sub(r, divisor)
			q.SetBit(q, bit, 1)
		}
	}
	if x.Sign() < 0 {
		r.Neg(r)
	}
	if x.Sign()*y.Sign() < 0 {
		q.Neg(q)
	}
	// An independent Euclidean result uses |y| for the remainder adjustment
	// and reconstructs the quotient by changing one unit in the opposite direction.
	e, m := new(big.Int).Set(q), new(big.Int).Set(r)
	if m.Sign() < 0 {
		m.Add(m, divisor)
		e.Sub(e, big.NewInt(int64(y.Sign())))
	}
	return fmt.Sprintf("%s|%s|%s,%s|%s|%s|%s,%s", q, r, q, r, e, m, e, m)
}

func native(f fixture, mode int) string {
	observe := func(call func(z, x, y, remainder *big.Int) string) (result string) {
		x, y := operand(f.x), operand(f.y)
		if f.shared {
			y = x
		}
		z := []*big.Int{big.NewInt(219), x, y}[mode]
		r := big.NewInt(-119)
		beforeZ, beforeR := z.String(), r.String()
		defer func() {
			if failure := recover(); failure != nil {
				if failure != "division by zero" {
					panic(failure)
				}
				if x.String() != f.x || y.String() != f.y || z.String() != beforeZ || r.String() != beforeR {
					log.Fatal("division by zero changed an operand or destination")
				}
				result = "panic"
			}
		}()
		result = call(z, x, y, r)
		if x != z && x.String() != f.x || y != z && y.String() != f.y {
			log.Fatal("operation modified a distinct operand")
		}
		return result
	}
	single := func(run func(*big.Int, *big.Int, *big.Int) *big.Int) string {
		return observe(func(z, x, y, _ *big.Int) string {
			if run(z, x, y) != z {
				log.Fatal("wrong result receiver")
			}
			return z.String()
		})
	}
	double := func(run func(*big.Int, *big.Int, *big.Int, *big.Int) (*big.Int, *big.Int)) string {
		return observe(func(z, x, y, r *big.Int) string {
			qr, rr := run(z, x, y, r)
			if qr != z || rr != r {
				log.Fatal("wrong quotient or remainder receiver")
			}
			return z.String() + "," + r.String()
		})
	}
	return strings.Join([]string{
		single((*big.Int).Quo), single((*big.Int).Rem), double((*big.Int).QuoRem),
		single((*big.Int).Div), single((*big.Int).Mod), double((*big.Int).DivMod),
	}, "|")
}

func compare(label string, cases []fixture) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.x, c.y)
	}
	for i, row := range check.LeanRows("godiff-math-big-int-division", len(cases), &input) {
		c := cases[i]
		want := oracle(c)
		if row != want {
			log.Fatalf("%s %+v: Lean got %q, want %q", label, c, row, want)
		}
		for mode := range 3 {
			if got := native(c, mode); got != want {
				log.Fatalf("%s %+v profile %d: Go got %q, want %q", label, c, mode, got, want)
			}
		}
	}
	fmt.Printf("%s: %d pairs agree under all applicable receiver alias profiles\n", label, len(cases))
}

func sample(r *rand.Rand) string {
	raw := make([]byte, r.Intn(25))
	r.Read(raw)
	value := new(big.Int).SetBytes(raw)
	if r.Intn(2) == 0 {
		value.Neg(value)
	}
	return value.String()
}

func main() {
	values := []string{"0", "1", "-1", "2", "-2", "3", "-3", "7", "-7",
		"18446744073709551615", "-18446744073709551616",
		"340282366920938463463374607431768211456",
		"-9" + strings.Repeat("9", 159), "1" + strings.Repeat("0", 159)}
	var fixed []fixture
	for _, x := range values {
		for _, y := range values {
			fixed = append(fixed, fixture{x, y, false})
		}
		fixed = append(fixed, fixture{x, x, true})
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for range 4 {
			cases := make([]fixture, 5_000)
			for i := range cases {
				x, y := sample(r), sample(r)
				shared := i%11 == 0
				if shared {
					y = x
				} else if i%7 == 0 {
					y = values[r.Intn(len(values))]
				}
				cases[i] = fixture{x, y, shared}
			}
			compare(fmt.Sprintf("seed=%d", seed), cases)
		}
	}
}
