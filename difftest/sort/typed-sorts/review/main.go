package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind   int
	Ints   []int
	Floats []uint64
	Words  [][]byte
	Want   string
}

func floats(words []uint64) string {
	parts := make([]string, len(words))
	for i, v := range words {
		parts[i] = strconv.FormatUint(v, 10)
	}
	return "values:" + strings.Join(parts, ",")
}

func words(values [][]byte) string {
	parts := make([]string, len(values))
	for i, word := range values {
		fields := make([]string, len(word))
		for j, b := range word {
			fields[j] = strconv.Itoa(int(b))
		}
		parts[i] = fmt.Sprintf("%d:%s", len(word), strings.Join(fields, ","))
	}
	return "values:" + strings.Join(parts, ";")
}

func main() {
	const pz = 0
	const nz = 0x8000000000000000
	const qi = 0xfff8000000000012
	const si = 0x7ff0000000000007
	const ns = 0xfff0000000000009
	const pq = 0x7ff800000000000b
	const inf = 0x7ff0000000000000
	const ninf = 0xfff0000000000000
	floatInput := []uint64{pz, qi, nz, si, inf, ninf, nz | 1, 1, ns, pq}
	sortedFloats := []uint64{qi, si, ns, pq, ninf, nz | 1, pz, nz, 1, inf}
	stringInput := [][]byte{{255}, {128, 0}, {'A', 0}, {239, 191, 189}, {}, {128}, {255, 0}, {0}, {'A'}}
	sortedWords := [][]byte{{}, {0}, {'A'}, {'A', 0}, {128}, {128, 0}, {239, 191, 189}, {255}, {255, 0}}
	cases := []fixture{
		{Kind: 0, Ints: []int{math.MaxInt, -1, 0, math.MinInt, -1}, Want: fmt.Sprintf("values:%d,-1,-1,0,%d", math.MinInt, math.MaxInt)},
		{Kind: 1, Floats: floatInput, Want: floats(sortedFloats)},
		{Kind: 2, Floats: floatInput, Want: floats(sortedFloats)},
		{Kind: 3, Words: stringInput, Want: words(sortedWords)},
		{Kind: 4, Words: stringInput, Want: words(sortedWords)},
		{Kind: 5, Floats: []uint64{si}, Want: "true"},
		{Kind: 5, Floats: []uint64{ns}, Want: "true"},
		{Kind: 5, Floats: []uint64{ninf}, Want: "false"},
		{Kind: 5, Floats: []uint64{nz}, Want: "false"},
	}
	unsignedInts := slices.Clone(cases[0].Ints)
	slices.SortFunc(unsignedInts, func(a, b int) int {
		if uint(a) < uint(b) {
			return -1
		}
		if uint(a) > uint(b) {
			return 1
		}
		return 0
	})
	if unsignedInts[0] == math.MinInt {
		panic("unsigned integer-order mutation survived")
	}
	strictZeros := slices.Clone(sortedFloats)
	strictZeros[6], strictZeros[7] = strictZeros[7], strictZeros[6]
	if floats(strictZeros) == cases[1].Want {
		panic("signed-zero total-order mutation survived")
	}
	noNaNOrder := slices.Clone(floatInput)
	slices.SortStableFunc(noNaNOrder, func(a, b uint64) int {
		x, y := math.Float64frombits(a), math.Float64frombits(b)
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	if floats(noNaNOrder) == cases[1].Want {
		panic("raw floating comparison mutation survived")
	}
	normalizedStrings := slices.Clone(stringInput)
	slices.SortStableFunc(normalizedStrings, func(a, b []byte) int {
		return strings.Compare(string([]rune(string(a))), string([]rune(string(b))))
	})
	if words(normalizedStrings) == cases[3].Want {
		panic("UTF-8-normalized key mutation survived")
	}
	if ninf&inf != inf || cases[7].Want != "false" {
		panic("exponent-only NaN mutation survived")
	}
	fmt.Println("detected strict signed zero order, raw floating compare, rune string keys and exponent-only NaN")

	file, err := os.CreateTemp("", "review-typed-sorts-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	var protocol bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&protocol, c.Kind, " ", strconv.IntSize)
		for _, v := range c.Ints {
			fmt.Fprint(&protocol, " ", v)
		}
		for _, v := range c.Floats {
			fmt.Fprint(&protocol, " ", v)
		}
		for _, data := range c.Words {
			fmt.Fprint(&protocol, " ", len(data))
			for _, b := range data {
				fmt.Fprint(&protocol, " ", b)
			}
		}
		fmt.Fprintln(&protocol)
		must(encoder.Encode(c))
	}
	must(file.Close())
	for i, row := range check.LeanRows("godiff-sort-typed", len(cases), &protocol) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("row %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_TYPED_SORT_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanTypedSorts", "difftest/sort/typed-sorts/testdata/typed_sorts_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
