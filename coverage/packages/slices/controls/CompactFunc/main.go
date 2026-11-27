package main

import (
	"fmt"
	"slices"
)

func main() {
	data := []int{2, 2, 1, 3}
	var pairs [][2]int
	out := slices.CompactFunc(data, func(current, previous int) bool {
		pairs = append(pairs, [2]int{current, previous})
		return current >= previous
	})
	if !slices.Equal(out, []int{2, 1}) || !slices.Equal(data, []int{2, 1, 0, 0}) ||
		!slices.Equal(pairs, [][2]int{{2, 2}, {1, 2}, {3, 1}}) {
		panic("unexpected successful compaction or callback trace")
	}
	reversed := slices.CompactFunc([]int{2, 2, 1, 3}, func(current, previous int) bool {
		return previous >= current
	})
	if slices.Equal(reversed, out) {
		panic("reversed comparator arguments were accepted")
	}

	original := []int{1, 1, 2, 3, 9}
	data = slices.Clone(original)
	token := new(int)
	actual := compactWithPanic(data, token)
	if actual != token || !slices.Equal(data, []int{1, 2, 2, 3, 9}) {
		panic("wrong panic payload or lost partial write")
	}
	buffered := slices.Clone(original)
	if compactWithPanic(slices.Clone(buffered), token) != token || slices.Equal(buffered, data) {
		panic("buffered-write mutant matched")
	}
	fmt.Println("CompactFunc: Go passed success and panic fixtures; reversed arguments and buffered writes rejected")
}

func compactWithPanic(data []int, token *int) (failure any) {
	defer func() { failure = recover() }()
	slices.CompactFunc(data, func(current, previous int) bool {
		if current == 3 {
			panic(token)
		}
		return current == previous
	})
	return nil
}
