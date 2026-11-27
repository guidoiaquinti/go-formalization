package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		Backing            []int
		A, B, Length, Mode int
		Want               string
	}{
		{nil, 0, 0, 0, 1, "nil"},
		{[]int{}, 0, 0, 0, 0, "slice:"},
		{[]int{91, 7, 8, 9, 42, 92}, 1, 4, 5, 1, "slice:91,9,8,7,42,92"},
		{[]int{13, 1, 2, 3, 4, 99}, 1, 5, 5, 0, "slice:13,4,3,2,1,99"},
		{[]int{2, 3, 4}, 1, 2, 2, 1, "slice:2,3,4"},
		{[]int{3, 2, 1}, 2, 2, 2, 1, "slice:3,2,1"},
	}
	file, err := os.CreateTemp("", "review-reverse-range-")
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
		fmt.Fprint(&input, tag, " ", c.A, " ", c.B, " ", c.Length, " ", c.Mode)
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
	cmd := exec.Command("lean/.lake/build/bin/godiff-slices-reverse-range")
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
	if err := os.Setenv("GOLEAN_SLICE_REVERSE_RANGE_CASES", file.Name()); err != nil {
		panic(err)
	}
	const data = "difftest/slices/reverse-range/testdata/"
	if err := check.TestInPackage("slices", "TestGoLeanReverseRange", data+"bridge_test.go", data+"reverse_range_test.go"); err != nil {
		panic(err)
	}
}
