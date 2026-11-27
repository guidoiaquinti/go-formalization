package main

import (
	"fmt"
	"log"
	"slices"
)

type ints []int

func main() {
	implementations := []struct {
		name          string
		clone         func(ints) ints
		discrepancies int
	}{
		{"slices.Clone", slices.Clone[ints], 0},
		{"return original descriptor", func(input ints) ints { return input }, 2},
		{"collapse empty to nil", func(input ints) ints { return append(ints(nil), input...) }, 2},
		{"convert nil to nonnil", func(input ints) ints { return append(ints{}, input...) }, 1},
		{"drop initial elements", func(input ints) ints {
			if len(input) > 0 {
				return slices.Clone(input[1:])
			}
			return slices.Clone(input)
		}, 2},
	}
	for _, candidate := range implementations {
		backing := ints{90, 10, 20, 40, 50}
		inputs := []ints{nil, {}, backing[2:2], backing[1:3], ints{7}}
		var mismatches int
		for _, input := range inputs {
			if !preservesAndCopies(candidate.clone, input) {
				mismatches++
			}
		}
		if mismatches != candidate.discrepancies {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, mismatches, candidate.discrepancies)
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, mismatches, len(inputs))
	}
}

func preservesAndCopies(clone func(ints) ints, input ints) bool {
	want := append(ints(nil), input...)
	result := clone(input)
	if !slices.Equal(result, want) || !slices.Equal(input, want) ||
		(result == nil) != (input == nil) {
		return false
	}
	for i := range result {
		result[i] ^= 1
		if !slices.Equal(input, want) {
			return false
		}
	}
	return true
}
