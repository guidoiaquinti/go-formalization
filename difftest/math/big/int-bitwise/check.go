package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"math/bits"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	x, y   *big.Int
	shift  uint
	index  int
	bit    uint
	shared bool
}

func main() {
	var fixed []fixture
	edge := []string{"0", "1", "-1", "2", "-2", "127", "-128", "255", "-256",
		"4294967296", "-4294967297", "18446744073709551616", "-18446744073709551617",
		"340282366920938463463374607431768211455", "-340282366920938463463374607431768211456",
		"-1" + strings.Repeat("0", 300)}
	scenarios := []struct {
		shift uint
		index int
		bit   uint
	}{
		{0, 0, 0}, {1, 0, 1}, {7, 7, 1}, {8, 8, 0}, {31, 32, 2}, {32, 31, 3}, {63, 64, 1}, {64, 63, 0},
		{65, 128, 1}, {257, 511, 0}, {511, 513, 1}, {1, -1, 2}, {0, -2, 0}, {8, 3, ^uint(0)},
	}
	for _, a := range edge {
		x := parse(a)
		for _, b := range edge {
			y := parse(b)
			for _, s := range scenarios {
				fixed = append(fixed, fixture{x, y, s.shift, s.index, s.bit, a == b})
			}
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				x, y := operand(rng), operand(rng)
				shared := i%11 == 0
				if shared {
					y = x
				}
				rows[i] = fixture{x, y, uint(rng.Intn(512)), rng.Intn(540) - 10, uint(rng.Intn(4)), shared}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), rows)
		}
	}
}

func parse(text string) *big.Int {
	v, ok := new(big.Int).SetString(text, 10)
	if !ok {
		log.Fatal("invalid fixture integer")
	}
	return v
}

func operand(rng *rand.Rand) *big.Int {
	raw := make([]byte, rng.Intn(48))
	rng.Read(raw)
	v := new(big.Int).SetBytes(raw)
	if rng.Intn(2) == 0 {
		v.Neg(v)
	}
	return v
}

func compare(label string, rows []fixture) {
	var input bytes.Buffer
	for _, f := range rows {
		fmt.Fprintln(&input, f.x, f.y, f.shift, f.index, f.bit)
	}
	for i, lean := range check.LeanRows("godiff-math-big-int-bitwise", len(rows), &input) {
		f := rows[i]
		want := oracle(f)
		if lean != want {
			log.Fatalf("%s row %d: Lean %q want %q", label, i, lean, want)
		}
		for mode := range 3 {
			if goResult := native(f, mode); goResult != want {
				log.Fatalf("%s row %d profile %d: native %q want %q", label, i, mode, goResult, want)
			}
		}
	}
	fmt.Printf("%s: %d pinned Go and Lean cases across receiver profiles\n", label, len(rows))
}

type operation func(z, x, y *big.Int) *big.Int

func native(f fixture, mode int) string {
	first, second := f.x.String(), f.y.String()
	run := func(op operation) string {
		x, y := new(big.Int).Set(f.x), new(big.Int).Set(f.y)
		if f.shared {
			y = x
		}
		z := []*big.Int{big.NewInt(-119), x, y}[mode]
		return observe(func() string {
			if op(z, x, y) != z || x != z && x.String() != first || y != z && y.String() != second {
				log.Fatal("wrong receiver or changed separate operand")
			}
			return z.String()
		})
	}
	x := new(big.Int).Set(f.x)
	b := observe(func() string { return strconv.FormatUint(uint64(x.Bit(f.index)), 10) })
	z := strconv.FormatUint(uint64(x.TrailingZeroBits()), 10)
	if x.String() != first {
		log.Fatal("read-only bit method changed input")
	}
	return strings.Join([]string{
		run((*big.Int).And), run((*big.Int).AndNot), run((*big.Int).Or), run((*big.Int).Xor),
		run(func(z, x, _ *big.Int) *big.Int { return z.Not(x) }),
		run(func(z, x, _ *big.Int) *big.Int { return z.Lsh(x, f.shift) }),
		run(func(z, x, _ *big.Int) *big.Int { return z.Rsh(x, f.shift) }),
		b, run(func(z, x, _ *big.Int) *big.Int { return z.SetBit(x, f.index, f.bit) }), z}, "|")
}

func observe(f func() string) (result string) {
	defer func() {
		if p := recover(); p != nil {
			switch p {
			case "negative bit index", "set bit is not 0 or 1":
				result = p.(string)
			default:
				panic(p)
			}
		}
	}()
	return f()
}

func oracle(f fixture) string {
	// Enough bytes for all shifted or changed bits plus the original sign.
	width := max(len(f.x.Bytes())+int(f.shift/8), len(f.y.Bytes()), f.index/8) + 2
	a, b := twos(f.x, width), twos(f.y, width)
	and, andNot, or, xor, not, left, right := make([]byte, width), make([]byte, width), make([]byte, width),
		make([]byte, width), make([]byte, width), make([]byte, width), make([]byte, width)
	for i := range a {
		and[i], andNot[i], or[i], xor[i], not[i] = a[i]&b[i], a[i]&^b[i], a[i]|b[i], a[i]^b[i], ^a[i]
	}
	shift := int(f.shift)
	for i := range width {
		j := i - shift/8
		k := i + shift/8
		offset := uint(shift % 8)
		left[i] = load(a, j, 0)<<offset | load(a, j-1, 0)>>(8-offset)
		pad := byte(0)
		if a[width-1] >= 128 {
			pad = 255
		}
		right[i] = load(a, k, pad)>>offset | load(a, k+1, pad)<<(8-offset)
	}
	bit, setBit := "negative bit index", "negative bit index"
	if f.index >= 0 {
		pos, mask := f.index/8, byte(1<<uint(f.index%8))
		bit = "0"
		if a[pos]&mask != 0 {
			bit = "1"
		}
		setBit = "set bit is not 0 or 1"
		if f.bit <= 1 {
			updated := bytes.Clone(a)
			updated[pos] &^= mask
			if f.bit == 1 {
				updated[pos] |= mask
			}
			setBit = signed(updated)
		}
	}
	out := []string{signed(and), signed(andNot), signed(or), signed(xor), signed(not), signed(left), signed(right),
		bit, setBit, strconv.Itoa(trailing(f.x.Bytes()))}
	return strings.Join(out, "|")
}

func twos(v *big.Int, width int) []byte {
	raw := v.Bytes()
	out := make([]byte, width)
	for i, b := range raw {
		out[len(raw)-1-i] = b
	}
	if v.Sign() < 0 {
		negate(out)
	}
	return out
}

func signed(encoded []byte) string {
	copy := bytes.Clone(encoded)
	neg := copy[len(copy)-1] >= 128
	if neg {
		negate(copy)
	}
	for i, j := 0, len(copy)-1; i < j; i, j = i+1, j-1 {
		copy[i], copy[j] = copy[j], copy[i]
	}
	text := new(big.Int).SetBytes(copy).String()
	if neg {
		text = "-" + text
	}
	return text
}

func negate(encoded []byte) {
	carry := uint16(1)
	for i, b := range encoded {
		sum := uint16(^b) + carry
		encoded[i] = byte(sum)
		carry = sum >> 8
	}
}

func load(raw []byte, i int, pad byte) byte {
	if i < 0 {
		return 0
	}
	if i >= len(raw) {
		return pad
	}
	return raw[i]
}

func trailing(magnitude []byte) int {
	for i := len(magnitude) - 1; i >= 0; i-- {
		if b := magnitude[i]; b != 0 {
			return (len(magnitude)-1-i)*8 + bits.TrailingZeros8(b)
		}
	}
	return 0
}
