package main

import (
	"bytes"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s     string
	count int
}

func main() {
	var edges []input
	for _, s := range []string{"", "a", " ", "-", "0", "=", "\t", "  x", "é🚀", "\xff\x00", "abc", "abcde"} {
		for _, n := range []int{math.MinInt, -1, 0, 1, 2, 3, 63, 64, 65, 127, 128, 129} {
			edges = append(edges, input{s, n})
		}
		if len(s) != 0 {
			for _, n := range []int{8192/len(s) - 1, 8192 / len(s), 8192/len(s) + 1, 2*8192/len(s) + 1} {
				edges = append(edges, input{s, n})
			}
		}
		if len(s) != 1 {
			edges = append(edges, input{s, math.MaxInt})
		}
		if len(s) > 1 {
			edges = append(edges, input{s, math.MaxInt/len(s) + 1})
		}
	}
	for _, size := range []int{8191, 8192, 8193} {
		edges = append(edges, input{strings.Repeat("x", size), 3})
	}
	for value := range 256 {
		edges = append(edges, input{string([]byte{byte(value)}), 2})
	}
	check.CompareSamples("godiff-strings-repeat", edges, sample, encodeInput, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	raw := make([]byte, rng.Intn(49))
	rng.Read(raw)
	if i%5 == 0 {
		return input{string(raw), -1 - rng.Intn(99)}
	}
	if i%5 == 1 && len(raw) != 1 {
		return input{string(raw), math.MaxInt}
	}
	return input{string(raw), rng.Intn(13)}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(math.MaxInt), 10)
	row = strconv.AppendInt(append(row, ' '), int64(c.count), 10)
	for i := range len(c.s) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.s[i]), 10)
	}
	return string(row)
}

func describe(s string) string {
	row := []byte("bytes")
	for i := range len(s) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(s[i]), 10)
	}
	return string(row)
}

func actual(c input) (result string) {
	defer func() {
		if p := recover(); p != nil {
			switch p {
			case "strings: negative Repeat count":
				result = "negative"
			case "strings: Repeat output length overflow":
				result = "overflow"
			default:
				panic(p)
			}
		}
	}()
	return describe(strings.Repeat(c.s, c.count))
}

func expected(c input) string {
	if c.count < 0 {
		return "negative"
	}
	if c.count != 0 && len(c.s) > math.MaxInt/c.count {
		return "overflow"
	}
	if c.s == "" {
		return "bytes"
	}
	var out bytes.Buffer
	for range c.count {
		out.WriteString(c.s)
	}
	return describe(out.String())
}
