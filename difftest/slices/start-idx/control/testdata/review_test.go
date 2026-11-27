package slices_test

import (
	"runtime"
	"slices"
	"testing"
)

func TestGoLeanStartIdxReview(t *testing.T) {
	backing := []uint64{7, 7, 7, 7, 7, 7}
	separate := []uint64{7, 7, 7}
	tests := []struct {
		hay, needle []uint64
		want        int
		status      string
	}{
		{backing[1:5], backing[3:4], 2, "ok"},
		{backing[1:4], backing[3:6], 2, "ok"},
		{backing[1:5], backing[:2], 0, "missing"},
		{backing[1:5], separate, 0, "missing"},
		{backing[1:1], backing[1:2], 0, "missing"},
		{backing[1:5], backing[2:2], 0, "bounds"},
		{nil, nil, 0, "bounds"},
	}
	for _, candidate := range []struct {
		name string
		run  func([]uint64, []uint64) int
		want int
	}{
		{"Go", slices.GoLeanStartIdx[uint64], 0},
		{"Compare values", byValue, 4},
		{"Require the entire needle to fit", requireFit, 1},
		{"Return zero for an empty needle", acceptEmpty, 2},
	} {
		failures := 0
		for i, c := range tests {
			index, status := observe(candidate.run, c.hay, c.needle)
			if status != c.status || (status == "ok" && index != c.want) {
				failures++
				if candidate.want == 0 {
					t.Fatalf("Go fixture %d: got (%d,%s), want (%d,%s)", i, index, status, c.want, c.status)
				}
			}
		}
		if failures != candidate.want {
			t.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, failures, candidate.want)
		}
		t.Logf("%s disagreed on %d of %d fixtures", candidate.name, failures, len(tests))
	}
}

func observe(run func([]uint64, []uint64) int, a, b []uint64) (index int, status string) {
	defer func() {
		caught := recover()
		switch caught.(type) {
		case nil:
		case runtime.Error:
			status = "bounds"
		default:
			if caught == "needle not found" {
				status = "missing"
			} else {
				status = "unexpected panic"
			}
		}
	}()
	return run(a, b), "ok"
}

func byValue(a, b []uint64) int {
	first := b[0]
	for i, v := range a {
		if v == first {
			return i
		}
	}
	panic("needle not found")
}

func requireFit(a, b []uint64) int {
	i := slices.GoLeanStartIdx(a, b)
	if len(b) > len(a)-i {
		panic("needle not found")
	}
	return i
}

func acceptEmpty(a, b []uint64) int {
	if len(b) == 0 {
		return 0
	}
	return slices.GoLeanStartIdx(a, b)
}
