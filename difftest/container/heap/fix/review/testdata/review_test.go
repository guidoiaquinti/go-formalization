package heap

import (
	"fmt"
	"slices"
	"testing"
)

type fixReviewHeap struct {
	values    []int
	trace     []string
	scripted  bool
	decisions uint
}

func (h *fixReviewHeap) Len() int {
	h.trace = append(h.trace, "Len")
	return len(h.values)
}
func (h *fixReviewHeap) Push(any) { panic("unexpected Push") }
func (h *fixReviewHeap) Pop() any { panic("unexpected Pop") }
func (h *fixReviewHeap) Less(i, j int) bool {
	h.trace = append(h.trace, fmt.Sprintf("L:%d,%d", i, j))
	if h.scripted {
		yes := h.decisions&1 != 0
		h.decisions >>= 1
		return yes
	}
	return h.values[i] < h.values[j]
}
func (h *fixReviewHeap) Swap(i, j int) {
	h.trace = append(h.trace, fmt.Sprintf("S:%d,%d", i, j))
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

func alwaysUp(h *fixReviewHeap, i int) {
	down(h, i, h.Len())
	up(h, i)
}

func rewindScript(h *fixReviewHeap, i int) {
	decisions := h.decisions
	if !down(h, i, h.Len()) {
		h.decisions = decisions
		up(h, i)
	}
}

func TestGoLeanHeapFixReview(t *testing.T) {
	fixtures := []struct {
		before    []int
		index     int
		scripted  bool
		decisions uint
		trace     []string
		after     []int
	}{
		{[]int{1}, 0, false, 0, []string{"Len"}, []int{1}},
		{[]int{1, 9, 4, 3, 5}, 1, false, 0, []string{"Len", "L:4,3", "L:3,1", "S:1,3"}, []int{1, 3, 4, 9, 5}},
		{[]int{5, 1, 7, 8, 9}, 1, false, 0, []string{"Len", "L:4,3", "L:3,1", "L:1,0", "S:0,1"}, []int{1, 5, 7, 8, 9}},
		{[]int{8, 4, 3, 2, 7}, 1, true, 4, []string{"Len", "L:4,3", "L:3,1", "L:1,0", "S:0,1"}, []int{4, 8, 3, 2, 7}},
		{[]int{1, 3, 4, 9, 5}, 1, false, 0, []string{"Len", "L:4,3", "L:4,1", "L:1,0"}, []int{1, 3, 4, 9, 5}},
	}
	for _, c := range []struct {
		name string
		fix  func(*fixReviewHeap, int)
		want int
	}{
		{"Go", func(h *fixReviewHeap, i int) { Fix(h, i) }, 0},
		{"Always sift up after sifting down", alwaysUp, 1},
		{"Reset Less replies before sifting up", rewindScript, 1},
	} {
		bad := 0
		for _, f := range fixtures {
			h := &fixReviewHeap{values: slices.Clone(f.before), scripted: f.scripted, decisions: f.decisions}
			c.fix(h, f.index)
			if !slices.Equal(h.trace, f.trace) || !slices.Equal(h.values, f.after) {
				bad++
				if c.want == 0 {
					t.Errorf("Fix(%v, %d) got %v and %v", f.before, f.index, h.trace, h.values)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
