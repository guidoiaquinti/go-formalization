package slices_test

import (
	"runtime"
	"slices"
	"testing"
)

type reviewSlice []int

func TestGoLeanRotateLeftReview(t *testing.T) {
	cases := []struct {
		length, count int
		want          []int
		panics        bool
	}{
		{4, 1, []int{19, 2, 3, 4, 1, 5}, false},
		{4, 4, []int{19, 1, 2, 3, 4, 5}, false},
		{3, 4, []int{19, 4, 3, 2, 1, 5}, true},
		{0, 3, []int{19, 3, 2, 1, 4, 5}, true},
		{3, 6, []int{19, 1, 2, 3, 4, 5}, true},
		{3, -1, []int{19, 1, 2, 3, 4, 5}, true},
	}
	for _, candidate := range []struct {
		name string
		run  func(reviewSlice, int)
		want int
	}{
		{"Go", slices.GoLeanRotateLeft[reviewSlice], 0},
		{"Prevalidate count against length", prevalidate, 2},
		{"Rotate right on valid input", wrongDirection, 1},
	} {
		failures := 0
		for i, c := range cases {
			backing := reviewSlice{19, 1, 2, 3, 4, 5}
			view := backing[1 : 1+c.length]
			failure := observe(candidate.run, view, c.count)
			_, bounds := failure.(runtime.Error)
			if c.panics && !bounds || !c.panics && failure != nil || !slices.Equal(backing, c.want) {
				failures++
				if candidate.want == 0 {
					t.Fatalf("Go fixture %d: backing %v panic %v", i, backing, failure)
				}
			}
		}
		if failures != candidate.want {
			t.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, failures, candidate.want)
		}
		t.Logf("%s disagreed on %d of %d fixtures", candidate.name, failures, len(cases))
	}
}

func observe(run func(reviewSlice, int), view reviewSlice, count int) (failure any) {
	defer func() { failure = recover() }()
	run(view, count)
	return nil
}

func prevalidate(view reviewSlice, count int) {
	_ = view[count:]
	slices.GoLeanRotateLeft(view, count)
}

func wrongDirection(view reviewSlice, count int) {
	if count >= 0 && count <= len(view) {
		count = len(view) - count
	}
	slices.GoLeanRotateLeft(view, count)
}
