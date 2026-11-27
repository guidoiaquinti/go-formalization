package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		Backing                    []int
		A, Middle, B, Length, Mode int
		Want                       string
	}{
		{nil, 0, 0, 0, 0, 1, "nil"},
		{[]int{}, 0, 0, 0, 0, 0, "slice:"},
		{[]int{90, 1, 2, 3, 4, 91, 99}, 1, 3, 5, 6, 1, "slice:90,3,4,1,2,91,99"},
		{[]int{90, 1, 2, 3, 4, 5, 91, 99}, 1, 2, 6, 7, 1, "slice:90,2,3,4,5,1,91,99"},
		{[]int{90, 1, 2, 3, 4, 5, 91, 99}, 1, 5, 6, 7, 0, "slice:90,5,1,2,3,4,91,99"},
		{[]int{90, 1, 2, 3, 4, 5, 6, 7, 91, 99}, 1, 4, 8, 9, 1, "slice:90,4,5,6,7,1,2,3,91,99"},
		{[]int{90, 1, 91, 99}, 2, 2, 2, 3, 1, "slice:90,1,91,99"},
	}
	for _, mutation := range []struct {
		name                       string
		dropFinal, wrongRightBlock bool
		fixture                    int
		want                       []int
	}{
		{"omit final swap", true, false, 2, []int{90, 3, 4, 1, 2, 91, 99}},
		{"swap the start of the longer right block", false, true, 3, []int{90, 2, 3, 4, 5, 1, 91, 99}},
	} {
		c := cases[mutation.fixture]
		cells := slices.Clone(c.Backing)
		swap := func(a, b, count int) {
			for i := range count {
				cells[a+i], cells[b+i] = cells[b+i], cells[a+i]
			}
		}
		i, j := c.Middle-c.A, c.B-c.Middle
		for i != j {
			if i > j {
				swap(c.Middle-i, c.Middle, j)
				i -= j
			} else {
				right := c.Middle + j - i
				if mutation.wrongRightBlock {
					right = c.Middle
				}
				swap(c.Middle-i, right, i)
				j -= i
			}
		}
		if !mutation.dropFinal {
			swap(c.Middle-i, c.Middle, i)
		}
		if slices.Equal(cells, mutation.want) {
			panic("mutation survived: " + mutation.name)
		}
		fmt.Printf("%s rejected with backing %v\n", mutation.name, cells)
	}

	file, err := os.CreateTemp("", "review-rotate-range-")
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
		fmt.Fprint(&input, tag, " ", c.A, " ", c.Middle, " ", c.B, " ", c.Length, " ", c.Mode)
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
	cmd := exec.Command("lean/.lake/build/bin/godiff-slices-rotate-range")
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
	if err := os.Setenv("GOLEAN_SLICE_ROTATE_RANGE_CASES", file.Name()); err != nil {
		panic(err)
	}
	const data = "difftest/slices/rotate-range/testdata/"
	if err := check.TestInPackage("slices", "TestGoLeanRotateRange", data+"bridge_test.go", data+"rotate_range_test.go"); err != nil {
		panic(err)
	}
}
