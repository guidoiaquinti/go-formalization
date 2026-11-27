package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		Backing                   []int
		A, B, Count, Length, Mode int
		Want                      string
	}{
		{nil, math.MinInt, math.MaxInt, 0, 0, 1, "ok:nil"},
		{[]int{}, 0, 0, 1, 0, 0, "panic:slice:"},
		{[]int{4, 7, 1, 8, 91}, 0, 1, 3, 4, 1, "ok:slice:7,1,8,4,91"},
		{[]int{4, 7, 1, 8, 91}, 0, 2, math.MaxInt, 4, 1, "panic:slice:1,8,4,7,91"},
		{[]int{4, 7, 1, 8, 91}, -1, 2, 1, 4, 1, "panic:slice:4,7,1,8,91"},
		{[]int{4, 7, 1, 8, 91}, 0, math.MaxInt, 1, 4, 0, "panic:slice:4,7,1,8,91"},
		{[]int{4, 7, 1, 8, 91}, 1, 1, 3, 4, 0, "ok:slice:4,7,1,8,91"},
		{[]int{4, 7, 1, 8, 91}, -1, 2, math.MinInt, 4, 1, "ok:slice:4,7,1,8,91"},
	}
	overlap := cases[2]
	backward := slices.Clone(overlap.Backing)
	for i := overlap.Count - 1; i >= 0; i-- {
		x, y := overlap.A+i, overlap.B+i
		backward[x], backward[y] = backward[y], backward[x]
	}
	if !slices.Equal(backward, []int{8, 4, 7, 1, 91}) ||
		slices.Equal(backward, []int{7, 1, 8, 4, 91}) {
		panic("descending-swap mutation survived")
	}
	partial := cases[3]
	// An atomic implementation would see the third invalid pair before doing the first two swaps.
	atomic := slices.Clone(partial.Backing)
	if partial.B+2 < partial.Length || slices.Equal(atomic, []int{1, 8, 4, 7, 91}) {
		panic("atomic-preflight mutation survived")
	}
	capacityBound := slices.Clone(partial.Backing)
	for i := 0; partial.B+i < len(capacityBound); i++ {
		x, y := partial.A+i, partial.B+i
		capacityBound[x], capacityBound[y] = capacityBound[y], capacityBound[x]
	}
	if !slices.Equal(capacityBound, []int{1, 8, 91, 7, 4}) ||
		slices.Equal(capacityBound, []int{1, 8, 4, 7, 91}) {
		panic("capacity-for-length mutation survived")
	}

	file, err := os.CreateTemp("", "review-swap-range-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())
	var input bytes.Buffer
	encoder := json.NewEncoder(file)
	for _, c := range cases {
		tag := 1
		if c.Backing == nil {
			tag = 0
		}
		fmt.Fprint(&input, tag, " ", c.A, " ", c.B, " ", c.Count, " ", c.Length, " ", c.Mode)
		for _, cell := range c.Backing {
			fmt.Fprint(&input, " ", cell)
		}
		fmt.Fprintln(&input)
		c.Want += "|" + c.Want
		if err := encoder.Encode(c); err != nil {
			panic(err)
		}
	}
	if err := file.Close(); err != nil {
		panic(err)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-slices-swap-range")
	cmd.Stdin = &input
	result, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	rows := strings.Split(strings.TrimSpace(string(result)), "\n")
	if len(rows) != len(cases) {
		panic("unexpected row count")
	}
	for i, c := range cases {
		if rows[i] != c.Want+"|"+c.Want {
			panic(fmt.Sprintf("fixture %d got %q", i, rows[i]))
		}
	}
	if err := os.Setenv("GOLEAN_SLICE_SWAP_RANGE_CASES", file.Name()); err != nil {
		panic(err)
	}
	const data = "difftest/slices/swap-range/testdata/"
	if err := check.TestInPackage("slices", "TestGoLeanSwapRange", data+"bridge_test.go", data+"swap_range_test.go"); err != nil {
		panic(err)
	}
}
