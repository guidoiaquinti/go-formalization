package main

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text     string
	count    int
	nilSlice bool
}

func main() {
	var edges []input
	for _, text := range []string{"", "a", "Ab", "4BCd", "a_b.c—d", "a\u00a0b\u2028c", "ǆ&ǆ", "\xff\xe2\x82�-a", "Σ⁄ς⁄σ", "ab\x00cd", "Ω 🦊É", "Zz"} {
		for _, n := range []int{math.MinInt, -7, -1, 0, 1, 2, 5, 9} {
			edges = append(edges, input{text, n, false})
		}
		if text == "" || len(text) >= 2 {
			edges = append(edges, input{text, math.MaxInt, false})
		}
	}
	for _, n := range []int{-1, 0, 1, math.MaxInt} {
		edges = append(edges, input{"", n, true})
	}
	for b := range 256 {
		text := string([]byte{byte(b)})
		edges = append(edges, input{text, 2, false}, input{"a" + text + "b", 3, false})
	}
	for _, size := range []int{8191, 8192, 8193} {
		text := strings.Repeat("aB", size/2) + strings.Repeat("z", size%2)
		for _, n := range []int{1, 2, 3} {
			edges = append(edges, input{text, n, false})
		}
	}
	for _, n := range []int{4095, 4096, 4097} {
		edges = append(edges, input{"xZ", n, false})
	}
	for _, n := range []int{2730, 2731, 5463} {
		edges = append(edges, input{"aZ\x00", n, false})
	}
	check.CompareSamples("godiff-bytes-repeat-title", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	var text strings.Builder
	words := []string{"foo", "bar", "X", "é", "ǆ", "世", "Σ", "—", "⁄", "'", "_", ".", "\x00", "\xff", "\xe2\x82", "�", "\u00a0", "\n"}
	for range rng.Intn(24) {
		if i%2 == 0 {
			text.WriteByte(byte(rng.Intn(256)))
		} else {
			text.WriteString(words[rng.Intn(len(words))])
		}
	}
	count := rng.Intn(13) - 3
	if i%17 == 0 && text.Len() >= 2 {
		count = math.MaxInt
	}
	return input{text.String(), count, text.Len() == 0 && i%2 == 0}
}

func encode(c input) string {
	row := fmt.Appendf(nil, "%d %d", math.MaxInt, c.count)
	for i := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
	}
	return string(row)
}

func describe(output string, nilSlice bool, capacity string) string {
	tag := "slice"
	if nilSlice {
		tag = "nil"
	}
	row := []byte(tag + ":" + capacity + ":")
	for i := range len(output) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(output[i]), 10)
	}
	return string(row)
}

func repeatNative(input []byte, count int) (result string) {
	defer func() {
		if p := recover(); p != nil {
			switch p {
			case "bytes: negative Repeat count", "bytes: Repeat output length overflow":
				result = p.(string)
			default:
				panic(p)
			}
		}
	}()
	out := bytes.Repeat(input, count)
	result = describe(string(out), out == nil, strconv.Itoa(cap(out)))
	for i := range out {
		out[i] ^= 0xff
	}
	return
}

func titleCapacity(input, output int) string {
	if output > input {
		return "growth"
	}
	return strconv.Itoa(input)
}

func actual(c input) string {
	backing := make([]byte, len(c.text)+2)
	copy(backing, c.text)
	backing[len(c.text)], backing[len(c.text)+1] = 0xa5, 0x59
	source := backing[:len(c.text)]
	if c.nilSlice {
		source = nil
	}
	assertSource := func() {
		if string(backing[:len(c.text)]) != c.text || backing[len(c.text)] != 0xa5 || backing[len(c.text)+1] != 0x59 {
			panic("result shares source backing")
		}
	}
	repeated := repeatNative(source, c.count)
	assertSource()
	title := bytes.Title(source)
	capacity := titleCapacity(len(source), len(title))
	if cap(title) < len(source) || (capacity != "growth" && strconv.Itoa(cap(title)) != capacity) {
		panic("unexpected title capacity")
	}
	titled := describe(string(title), title == nil, capacity)
	for i := range title {
		title[i] ^= 0xff
	}
	assertSource()
	return repeated + "|" + titled
}

func repeatRef(c input) string {
	if c.count < 0 {
		return "bytes: negative Repeat count"
	}
	size := new(big.Int).Mul(big.NewInt(int64(len(c.text))), big.NewInt(int64(c.count)))
	if size.Cmp(big.NewInt(math.MaxInt)) > 0 {
		return "bytes: Repeat output length overflow"
	}
	var output strings.Builder
	if size.Sign() > 0 {
		for range c.count {
			output.WriteString(c.text)
		}
	}
	return describe(output.String(), false, strconv.Itoa(output.Len()))
}

func separator(r rune) bool {
	if r > 127 {
		return unicode.IsSpace(r)
	}
	return !(r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r == '_')
}

func expected(c input) string {
	var output strings.Builder
	prev := rune(' ')
	for _, current := range c.text {
		mapped := current
		if separator(prev) {
			mapped = unicode.ToTitle(current)
		}
		output.WriteRune(mapped)
		prev = current
	}
	titled := describe(output.String(), false, titleCapacity(len(c.text), output.Len()))
	return repeatRef(c) + "|" + titled
}
