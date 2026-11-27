package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type divisionPair struct{ quotient, remainder int64 }

type fixture struct {
	x, y                 int64
	truncated, euclidean divisionPair
}

func main() {
	cases := []fixture{
		{-7, 3, divisionPair{-2, -1}, divisionPair{-3, 2}},
		{-7, -3, divisionPair{2, -1}, divisionPair{3, 2}},
		{7, -3, divisionPair{-2, 1}, divisionPair{-2, 1}},
		{-1, 7, divisionPair{0, -1}, divisionPair{-1, 6}},
		{-6, -3, divisionPair{2, 0}, divisionPair{2, 0}},
		{0, 5, divisionPair{0, 0}, divisionPair{0, 0}},
		{7, 0, divisionPair{}, divisionPair{}},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.x, c.y)
	}
	var rejected [7]int
	for i, lean := range check.LeanRows("godiff-math-big-int-division", len(cases), &input) {
		c := cases[i]
		want := c.render()
		if lean != want {
			log.Fatalf("row %d got %q want %q", i, lean, want)
		}
		for profile := range 3 {
			if got := native(c, profile); got != want {
				log.Fatalf("row %d profile %d: native %q want %q", i, profile, got, want)
			}
		}
		for j, wrong := range mutants(c) {
			if wrong != lean {
				rejected[j]++
			}
		}
	}
	if want := [7]int{3, 3, 3, 1, 3, 3, 1}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	fmt.Printf("Int division: %d pinned native and Lean fixtures in three profiles, mutation counts %v\n", len(cases), rejected)
}

func (c fixture) render() string {
	if c.y == 0 {
		return "panic|panic|panic|panic|panic|panic"
	}
	q, r, e, m := c.truncated.quotient, c.truncated.remainder, c.euclidean.quotient, c.euclidean.remainder
	return fmt.Sprintf("%d|%d|%d,%d|%d|%d|%d,%d", q, r, q, r, e, m, e, m)
}

func mutants(c fixture) [7]string {
	baseline := c.render()
	out := [7]string{baseline, baseline, baseline, baseline, baseline, baseline, baseline}
	if c.y == 0 {
		out[6] = "0|0|0,0|0|0|0,0"
		return out
	}
	q, r, e, m := c.truncated.quotient, c.truncated.remainder, c.euclidean.quotient, c.euclidean.remainder
	absR, absQ := r, q
	if absR < 0 {
		absR = -absR
	}
	if absQ < 0 {
		absQ = -absQ
	}
	wrongQ, wrongR := e, m
	if r < 0 {
		wrongQ = q - 1
		wrongR = r + absQ
	}
	fields := []string{fmt.Sprint(e), fmt.Sprint(absR), fmt.Sprintf("%d,%d", q, absR),
		fmt.Sprint(wrongQ), fmt.Sprint(r), fmt.Sprintf("%d,%d", e, wrongR)}
	for i, value := range fields {
		observations := strings.Split(baseline, "|")
		observations[i] = value
		out[i] = strings.Join(observations, "|")
	}
	return out
}

func native(c fixture, profile int) string {
	type operation func(z, x, y, r *big.Int) string
	single := func(call func(*big.Int, *big.Int, *big.Int) *big.Int) operation {
		return func(z, x, y, _ *big.Int) string {
			if call(z, x, y) != z {
				log.Fatal("wrong quotient or remainder pointer")
			}
			return z.String()
		}
	}
	pair := func(call func(*big.Int, *big.Int, *big.Int, *big.Int) (*big.Int, *big.Int)) operation {
		return func(z, x, y, r *big.Int) string {
			q, rem := call(z, x, y, r)
			if q != z || rem != r {
				log.Fatal("wrong result pointers")
			}
			return q.String() + "," + rem.String()
		}
	}
	run := func(call operation) (result string) {
		x, y := big.NewInt(c.x), big.NewInt(c.y)
		z, r := []*big.Int{big.NewInt(18), x, y}[profile], big.NewInt(-19)
		saved := z.Int64()
		defer func() {
			if p := recover(); p != nil {
				if p != "division by zero" || x.Int64() != c.x || y.Int64() != c.y || z.Int64() != saved || r.Int64() != -19 {
					log.Fatalf("unexpected panic or mutation: %v, operands (%s,%s), destinations (%s,%s)", p, x, y, z, r)
				}
				result = "panic"
			}
		}()
		result = call(z, x, y, r)
		if x != z && x.Int64() != c.x || y != z && y.Int64() != c.y {
			log.Fatal("an independent operand changed")
		}
		return result
	}
	return strings.Join([]string{run(single((*big.Int).Quo)), run(single((*big.Int).Rem)), run(pair((*big.Int).QuoRem)),
		run(single((*big.Int).Div)), run(single((*big.Int).Mod)), run(pair((*big.Int).DivMod))}, "|")
}
