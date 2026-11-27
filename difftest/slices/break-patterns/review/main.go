package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		Backing            []int
		A, B, Length, Mode int
		Want               string
	}{
		{nil, 1, 0, 0, 1, "ok:nil"},
		{[]int{}, 0, 8, 0, 0, "panic:slice:"},
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1, 8, 10, 0, "ok:slice:0,1,2,3,4,5,6,7,8,9,10"},
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1, 9, 10, 1, "ok:slice:0,6,2,3,1,5,4,7,8,9,10"},
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, 2, 11, 12, 1, "ok:slice:0,1,5,7,4,2,6,3,8,9,10,11,12"},
		{[]int{0, 1, 2, 3, 4}, 0, 8, 4, 1, "panic:slice:3,1,2,0,4"},
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 0, 15, 15, 0,
			"ok:slice:5,1,2,3,4,0,8,7,6,9,10,11,12,13,14,15"},
	}
	threshold := cases[3]
	if threshold.B-threshold.A > 8 || slices.Equal(threshold.Backing, []int{0, 6, 2, 3, 1, 5, 4, 7, 8, 9, 10}) {
		panic("length-greater-than-eight mutation survived")
	}
	nine := cases[4]
	candidate := slices.Clone(nine.Backing)
	word := uint64(nine.B - nine.A)
	for idx := 5; idx <= 7; idx++ {
		word ^= word << 13
		word ^= word >> 7
		word ^= word << 17
		other := nine.A + int(uint(word)%uint(nine.B-nine.A))
		candidate[idx], candidate[other] = candidate[other], candidate[idx]
	}
	if slices.Equal(candidate, []int{0, 1, 5, 7, 4, 2, 6, 3, 8, 9, 10, 11, 12}) {
		panic(fmt.Sprintf("unexpected modulo mutation backing %v", candidate))
	}
	partial := cases[5]
	if partial.Length > 4 || slices.Equal(partial.Backing, []int{3, 1, 2, 0, 4}) {
		panic("preflight mutation survived")
	}
	fmt.Printf("rejected threshold, unsigned modulo, and preflight mutations; modulo backing %v\n", candidate)

	file, err := os.CreateTemp("", "review-break-patterns-")
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
		fmt.Fprint(&input, tag, " ", strconv.IntSize, " ", c.A, " ", c.B, " ", c.Length, " ", c.Mode)
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
	cmd := exec.Command("lean/.lake/build/bin/godiff-slices-break-patterns")
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
	if err := os.Setenv("GOLEAN_SLICES_BREAK_PATTERNS_CASES", file.Name()); err != nil {
		panic(err)
	}
	const data = "difftest/slices/break-patterns/testdata/"
	if err := check.TestInPackage("slices", "TestGoLeanBreakPatterns", data+"bridge_test.go", data+"break_patterns_test.go"); err != nil {
		panic(err)
	}
}
