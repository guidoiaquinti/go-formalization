package heap

import (
	"fmt"
	"slices"
	"testing"
)

type initReviewHeap struct {
	values    []int
	trace     []string
	scripted  bool
	decisions uint
}

func (h *initReviewHeap) Len() int {
	h.trace = append(h.trace, "Len")
	return len(h.values)
}
func (h *initReviewHeap) Push(any) { panic("unexpected Push") }
func (h *initReviewHeap) Pop() any { panic("unexpected Pop") }
func (h *initReviewHeap) Less(i, j int) bool {
	h.trace = append(h.trace, fmt.Sprintf("L:%d,%d", i, j))
	if h.scripted {
		yes := h.decisions&1 != 0
		h.decisions >>= 1
		return yes
	}
	return h.values[i] < h.values[j]
}
func (h *initReviewHeap) Swap(i, j int) {
	h.trace = append(h.trace, fmt.Sprintf("S:%d,%d", i, j))
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

func initReviewAscending(h *initReviewHeap) {
	n := h.Len()
	for parent := 0; parent < n/2; parent++ {
		down(h, parent, n)
	}
}

func initReviewResetScript(h *initReviewHeap) {
	n := h.Len()
	decisions := h.decisions
	for parent := n/2 - 1; parent >= 0; parent-- {
		h.decisions = decisions
		down(h, parent, n)
	}
}

func TestGoLeanHeapInitReview(t *testing.T) {
	input := []int{9, 7, 5, 3, 1}
	for _, c := range []struct {
		name string
		init func(*initReviewHeap)
		want int
	}{
		{"Go", func(h *initReviewHeap) { Init(h) }, 0},
		{"Visit parents from the root first", initReviewAscending, 2},
		{"Reset comparison replies before each sift", initReviewResetScript, 1},
	} {
		bad := 0
		for _, f := range []struct {
			before    []int
			scripted  bool
			decisions uint
			trace     []string
			after     []int
		}{
			{nil, false, 0, []string{"Len"}, nil},
			{[]int{2}, false, 0, []string{"Len"}, []int{2}},
			{input, false, 0, []string{"Len", "L:4,3", "L:4,1", "S:1,4", "L:2,1", "L:1,0", "S:0,1", "L:4,3", "L:3,1", "S:1,3"}, []int{1, 3, 5, 9, 7}},
			{input, true, 25, []string{"Len", "L:4,3", "L:4,1", "L:2,1", "L:1,0", "S:0,1", "L:4,3", "L:4,1"}, []int{7, 9, 5, 3, 1}},
		} {
			h := &initReviewHeap{values: slices.Clone(f.before), scripted: f.scripted, decisions: f.decisions}
			c.init(h)
			if !slices.Equal(h.trace, f.trace) || !slices.Equal(h.values, f.after) {
				bad++
				if c.want == 0 {
					t.Errorf("Init(%v, scripted=%t) got %v and %v", f.before, f.scripted, h.trace, h.values)
				}
			}
		}
		t.Logf("%s disagreed on %d of 4", c.name, bad)
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
