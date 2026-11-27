package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

func main() {
	backing := []int{1, 2, 3, 99}
	captured := backing[:3]
	seq := slices.Values(captured)
	captured = append(captured, 88)
	captured = append(captured, 77)
	captured[0] = 100
	backing[1] = 20
	if got := slices.Collect(seq); !slices.Equal(got, []int{1, 20, 3}) {
		panic(fmt.Sprintf("captured view changed: %v", got))
	}
	var first []int
	seq(func(v int) bool {
		first = append(first, v)
		backing[1] = 30
		return false
	})
	if !slices.Equal(first, []int{1}) || !slices.Equal(slices.Collect(seq), []int{1, 30, 3}) {
		panic("false did not stop or restart")
	}
	marker := new(int)
	func() {
		defer func() {
			if recover() != marker {
				panic("callback panic lost identity")
			}
		}()
		seq(func(v int) bool { backing[1] = 40; panic(marker) })
	}()
	if got := slices.Collect(seq); !slices.Equal(got, []int{1, 40, 3}) {
		panic(fmt.Sprintf("panic lost edits: %v", got))
	}
	slices.Values([]int{})(nil)

	cases := []struct{ input, want string }{
		{"0 -1 5 1 2 3", "ok 1 7 8 trace 1 7 8 replay 1 7 8"},
		{"1 0 5 1 2 3", "ok 1 7 3 trace 1 replay 1 7 3"},
		{"2 1 5 1 2 3", "panic 2 1 7 8 trace 1 7 replay 1 7 8"},
		{"3 -1 0", "ok trace replay"},
		{"3 -1 0 9", "nil-panic 9 trace replay 9"},
	}
	var model bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&model, c.input)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-slices-values")
	cmd.Stdin = &model
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		panic("wrong result count")
	}
	for i, c := range cases {
		if rows[i] != c.want {
			panic(fmt.Sprintf("row %d got %q want %q", i, rows[i], c.want))
		}
	}
	fmt.Println("slices.Values passed five independent model rows and native capture, stop, panic and replay checks")
}
