package main

import (
	"fmt"
	"log"
	"slices"
)

type ints []int

func main() {
	backing := ints{90, 10, 20, 40, 50}
	offset := backing[1:3]
	inputs := []ints{nil, {}, backing[2:2], offset, offset[:len(offset):len(offset)]}
	implementations := []struct {
		name          string
		clip          func(ints) ints
		discrepancies int
	}{
		{"slices.Clip", slices.Clip[ints], 0},
		{"preserve old capacity", func(input ints) ints { return input }, 2},
		{"copy instead of sharing", func(input ints) ints {
			cloned := slices.Clone(input)
			return cloned[:len(cloned):len(cloned)]
		}, 2},
		{"collapse empty to nil", func(input ints) ints {
			if len(input) == 0 {
				return nil
			}
			return slices.Clip(input)
		}, 2},
	}
	for _, candidate := range implementations {
		var mismatches int
		for _, input := range inputs {
			got := candidate.clip(input)
			if len(got) != len(input) || cap(got) != len(input) ||
				(got == nil) != (input == nil) || !slices.Equal(got, input) ||
				(len(got) > 0 && &got[0] != &input[0]) {
				mismatches++
			}
		}
		if mismatches != candidate.discrepancies {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, mismatches, candidate.discrepancies)
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, mismatches, len(inputs))
	}
}
