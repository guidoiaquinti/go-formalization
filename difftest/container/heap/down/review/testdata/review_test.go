package heap

import (
	"math"
	"slices"
	"testing"
)

type downReviewCall struct {
	kind          byte
	first, second int
}

type downReviewHeap struct {
	decisions uint
	trace     []downReviewCall
}

func (h *downReviewHeap) Len() int { panic("unexpected Len") }
func (h *downReviewHeap) Push(any) { panic("unexpected Push") }
func (h *downReviewHeap) Pop() any { panic("unexpected Pop") }
func (h *downReviewHeap) Less(i, j int) bool {
	h.trace = append(h.trace, downReviewCall{'L', i, j})
	decision := h.decisions&1 != 0
	h.decisions >>= 1
	return decision
}
func (h *downReviewHeap) Swap(i, j int) {
	h.trace = append(h.trace, downReviewCall{'S', i, j})
}

func downReviewWithoutOverflowGuard(h Interface, start, n int) bool {
	i := start
	for {
		left := 2*i + 1
		if left >= n {
			break
		}
		child := left
		if right := left + 1; right < n && h.Less(right, left) {
			child = right
		}
		if !h.Less(child, i) {
			break
		}
		h.Swap(i, child)
		i = child
	}
	return i > start
}

func downReviewLeftOnly(h Interface, start, n int) bool {
	i := start
	for child := 2*i + 1; child >= 0 && child < n; child = 2*i + 1 {
		if !h.Less(child, i) {
			break
		}
		h.Swap(i, child)
		i = child
	}
	return i > start
}

func TestGoLeanHeapDownReview(t *testing.T) {
	fixtures := []struct {
		i, n      int
		decisions uint
		moved     bool
		trace     []downReviewCall
	}{
		{0, 0, 0, false, nil},
		{0, 2, 0, false, []downReviewCall{{'L', 1, 0}}},
		{0, 2, 1, true, []downReviewCall{{'L', 1, 0}, {'S', 0, 1}}},
		{0, 3, 1, false, []downReviewCall{{'L', 2, 1}, {'L', 2, 0}}},
		{0, 3, 3, true, []downReviewCall{{'L', 2, 1}, {'L', 2, 0}, {'S', 0, 2}}},
		{0, 5, 6, true, []downReviewCall{{'L', 2, 1}, {'L', 1, 0}, {'S', 0, 1}, {'L', 4, 3}, {'L', 4, 1}}},
		{math.MaxInt / 2, math.MaxInt, 0, false, nil},
		{math.MaxInt/2 + 1, math.MaxInt, 0, false, nil},
		{math.MaxInt, math.MaxInt, 0, false, nil},
	}
	for _, c := range []struct {
		name string
		down func(Interface, int, int) bool
		want int
	}{
		{"Go", down, 0},
		{"Ignore signed overflow", downReviewWithoutOverflowGuard, 2},
		{"Use only the left child", downReviewLeftOnly, 3},
		{"Report a move after any visit", func(h Interface, i, n int) bool {
			return down(h, i, n) || (i >= 0 && i < n/2)
		}, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			h := &downReviewHeap{decisions: f.decisions}
			got := c.down(h, f.i, f.n)
			if got != f.moved || !slices.Equal(h.trace, f.trace) {
				bad++
				if c.want == 0 {
					t.Errorf("down(%d,%d,%d): got %t %v, want %t %v", f.i, f.n, f.decisions, got, h.trace, f.moved, f.trace)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
