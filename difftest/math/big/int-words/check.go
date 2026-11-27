package main

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	value   *big.Int
	word    uint64
	present bool
}

func words(value *big.Int, width int) []uint64 {
	magnitude := value.Bytes()
	var out []uint64
	for len(magnitude) > 0 {
		start := max(0, len(magnitude)-width/8)
		var word uint64
		for _, digit := range magnitude[start:] {
			word = word<<8 | uint64(digit)
		}
		out = append(out, word)
		magnitude = magnitude[:start]
	}
	return out
}

func sign(value *big.Int) int {
	text := value.String()
	if text == "0" {
		return 0
	}
	if text[0] == '-' {
		return -1
	}
	return 1
}

func wordParts(word uint64, width int) []uint64 {
	if word == 0 {
		return nil
	}
	if width == 32 && word>>32 != 0 {
		return []uint64{uint64(uint32(word)), word >> 32}
	}
	return []uint64{word}
}

func wire(negative bool, values []uint64) string {
	fields := make([]string, len(values))
	for i, v := range values {
		fields[i] = strconv.FormatUint(v, 10)
	}
	return fmt.Sprintf("%t:%s", negative, strings.Join(fields, ","))
}

func model(c fixture, width int) string {
	negative := int64(c.word) < 0
	magnitude := c.word
	if negative {
		magnitude = -magnitude
	}
	s := wire(negative, wordParts(magnitude, width))
	if !c.present {
		return s + strings.Repeat("|panic", 8)
	}
	u := wire(false, wordParts(c.word, width))
	mag := c.value.Bytes()
	var length int
	if len(mag) > 0 {
		length = (len(mag)-1)*8 + bits.Len8(mag[0])
	}
	var low uint64
	for _, digit := range mag[max(0, len(mag)-8):] {
		low = low<<8 | uint64(digit)
	}
	direction := sign(c.value)
	signedLow := low
	if direction < 0 {
		signedLow = -signedLow
	}
	fitsSigned := length < 64 || length == 64 && direction < 0 && low == 1<<63
	fitsUnsigned := direction >= 0 && length <= 64
	return fmt.Sprintf("%s|%s|%s|%d|%d|%d|%t|%t|%d", s, s, u, direction, int64(signedLow), low, fitsSigned, fitsUnsigned, length)
}

func actual(c fixture) string {
	show := func(x *big.Int) string {
		ws := x.Bits()
		out := make([]uint64, len(ws))
		for i, w := range ws {
			out[i] = uint64(w)
		}
		return wire(x.Sign() < 0, out)
	}
	observe := func(f func() string) (res string) {
		defer func() {
			if p := recover(); p != nil {
				if err, ok := p.(runtime.Error); !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
					panic(p)
				}
				res = "panic"
			}
		}()
		return f()
	}
	fresh := big.NewInt(int64(c.word))
	another := big.NewInt(int64(c.word))
	if fresh == another {
		panic("NewInt returned the same receiver")
	}
	before := show(fresh)
	another.SetInt64(123)
	if before != show(fresh) {
		panic("NewInt backing aliases another instance")
	}
	var a, b, x *big.Int
	if c.present {
		a = new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 137))
		b = new(big.Int).Set(a)
		x = c.value
	}
	return strings.Join([]string{before,
		observe(func() string {
			if a.SetInt64(int64(c.word)) != a {
				panic("wrong receiver")
			}
			return show(a)
		}),
		observe(func() string {
			if b.SetUint64(c.word) != b {
				panic("wrong receiver")
			}
			return show(b)
		}),
		observe(func() string { return fmt.Sprint(x.Sign()) }),
		observe(func() string { return fmt.Sprint(x.Int64()) }),
		observe(func() string { return fmt.Sprint(x.Uint64()) }),
		observe(func() string { return fmt.Sprint(x.IsInt64()) }),
		observe(func() string { return fmt.Sprint(x.IsUint64()) }),
		observe(func() string { return fmt.Sprint(x.BitLen()) })}, "|")
}

func compare(label string, cases []fixture) {
	for _, width := range []int{32, 64} {
		var rows bytes.Buffer
		for _, c := range cases {
			live, negative := 0, 0
			if c.present {
				live = 1
			}
			if sign(c.value) < 0 {
				negative = 1
			}
			fmt.Fprint(&rows, width, " ", live, " ", negative, " ", c.word)
			for _, w := range words(c.value, width) {
				fmt.Fprint(&rows, " ", w)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-math-big-int-words", len(cases), &rows) {
			c := cases[i]
			want := model(c, width)
			if got != want {
				panic(fmt.Sprintf("width %d %s row %d %v Lean=%q want=%q", width, label, i, c, got, want))
			}
			if width == strconv.IntSize {
				if native := actual(c); native != want {
					panic(fmt.Sprintf("width %d %s row %d %v Go=%q want=%q", width, label, i, c, native, want))
				}
			}
		}
		fmt.Printf("math/big Int words width=%d %s: %d cases agree\n", width, label, len(cases))
	}
}

func main() {
	var edges []fixture
	for _, bit := range []uint{0, 1, 7, 8, 31, 32, 33, 62, 63, 64, 65, 95, 96, 127, 128, 256, 511} {
		for _, delta := range []int64{-1, 0, 1} {
			n := new(big.Int).Lsh(big.NewInt(1), bit)
			n.Add(n, big.NewInt(delta))
			for _, sign := range []int64{-1, 1} {
				v := new(big.Int).Mul(n, big.NewInt(sign))
				for _, word := range []uint64{0, 1, math.MaxUint32, 1 << 32, math.MaxInt64, 1 << 63, math.MaxUint64} {
					edges = append(edges, fixture{v, word, true}, fixture{v, word, false})
				}
			}
		}
	}
	compare("boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20000)
		for i := range cases {
			raw := make([]byte, rng.Intn(65))
			rng.Read(raw)
			n := new(big.Int).SetBytes(raw)
			if i%2 == 0 {
				n.Neg(n)
			}
			cases[i] = fixture{n, rng.Uint64(), i%7 != 0}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
}
