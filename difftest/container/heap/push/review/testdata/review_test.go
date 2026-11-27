package heap

import (
	"fmt"
	"slices"
	"testing"
)

type pushReviewHeap struct {
	values []int
	trace  []string
}

func (h *pushReviewHeap) Len() int {
	h.trace = append(h.trace, fmt.Sprintf("Len:%d", len(h.values)))
	return len(h.values)
}
func (h *pushReviewHeap) Push(x any) {
	h.trace = append(h.trace, fmt.Sprintf("Push:%d", x))
	h.values = append(h.values, x.(int))
}
func (h *pushReviewHeap) Pop() any { panic("unexpected Pop") }
func (h *pushReviewHeap) Less(i, j int) bool {
	h.trace = append(h.trace, fmt.Sprintf("Less:%d:%d", i, j))
	return h.values[i] < h.values[j]
}
func (h *pushReviewHeap) Swap(i, j int) {
	h.trace = append(h.trace, fmt.Sprintf("Swap:%d:%d", i, j))
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

func pushReviewEarlyLength(h Interface, x any) {
	n := h.Len()
	h.Push(x)
	up(h, n-1)
}

func pushReviewNoSift(h Interface, x any) {
	h.Push(x)
	h.Len()
}

func TestGoLeanHeapPushReview(t *testing.T) {
	fixtures := []struct {
		before []int
		value  int
		trace  []string
		after  []int
	}{
		{nil, 5, []string{"Push:5", "Len:1"}, []int{5}},
		{[]int{5}, 7, []string{"Push:7", "Len:2", "Less:1:0"}, []int{5, 7}},
		{[]int{5, 9, 7}, 3, []string{"Push:3", "Len:4", "Less:3:1", "Swap:1:3", "Less:1:0", "Swap:0:1"}, []int{3, 5, 7, 9}},
		{[]int{5, 9, 7}, 6, []string{"Push:6", "Len:4", "Less:3:1", "Swap:1:3", "Less:1:0"}, []int{5, 6, 7, 9}},
	}
	for _, c := range []struct {
		name string
		push func(Interface, any)
		want int
	}{
		{"Go", Push, 0},
		{"Read length before insertion", pushReviewEarlyLength, 4},
		{"Omit upward sift", pushReviewNoSift, 3},
	} {
		bad := 0
		for _, f := range fixtures {
			h := &pushReviewHeap{values: slices.Clone(f.before)}
			c.push(h, f.value)
			if !slices.Equal(h.trace, f.trace) || !slices.Equal(h.values, f.after) {
				bad++
				if c.want == 0 {
					t.Errorf("Push(%v,%d) got %v and %v", f.before, f.value, h.trace, h.values)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
