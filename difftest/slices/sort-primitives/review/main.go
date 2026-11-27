package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Length  int
	Word    uint64
	Present bool
	Want    string
}

func main() {
	cases := []fixture{
		{0, 0, true, "0,0/1/false/false/false"},
		{32, 1, true, "1082269761,1082269761/64/false/false/false"},
		{-1, 0x8000000000000000, true, "9295429630892703744,9295429630892703744/0/false/false/false"},
		{math.MaxInt, 0, false, fmt.Sprintf("panic/%d/false/false/false", uint(1)<<(strconv.IntSize-1))},
		{31, 0x7f800001, true, "2296608287824355393,2296608287824355393/32/true/false/false"},
		{1, 0x7ff0000000000000, true, "13767468876499517440,13767468876499517440/2/false/false/false"},
		{2, 0xfff0000000000001, true, "4472039246689083457,4472039246689083457/4/false/true/false"},
		{3, 0xff800000, true, "4584140164581752832,4584140164581752832/4/false/false/false"},
	}
	inputWord := cases[1].Word
	parallelStages := inputWord ^ (inputWord << 13) ^ (inputWord >> 7) ^ (inputWord << 17)
	if parallelStages != 139265 || parallelStages == 1082269761 {
		panic("parallel xorshift mutation survived")
	}
	inclusivePower := uint(1)
	for inclusivePower < uint(cases[1].Length) {
		inclusivePower <<= 1
	}
	if inclusivePower != 32 || inclusivePower == 64 {
		panic("inclusive power-of-two mutation survived")
	}
	infinity := cases[5].Word
	exponentOnly := (infinity>>52)&0x7ff == 0x7ff
	if !exponentOnly || math.IsNaN(math.Float64frombits(infinity)) {
		panic("exponent-only NaN mutation survived")
	}
	fmt.Println("rejected parallel xorshift stages, inclusive power of two and exponent-only NaN")
	file, err := os.CreateTemp("", "review-sort-primitives-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	var input bytes.Buffer
	for _, c := range cases {
		state := 0
		if c.Present {
			state = 1
		}
		fmt.Fprintln(&input, strconv.IntSize, c.Length, state, c.Word)
		must(encoder.Encode(c))
	}
	must(file.Close())
	for i, got := range check.LeanRows("godiff-slices-sort-primitives", len(cases), &input) {
		if got != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, got, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_SLICES_SORT_PRIMITIVES_CASES", file.Name()))
	const data = "difftest/slices/sort-primitives/testdata/"
	must(check.TestInPackage("slices", "TestGoLeanSortPrimitives", data+"bridge_test.go", data+"primitives_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
