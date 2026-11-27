package heap

import (
	"fmt"
	"slices"
	"testing"
)

type popReviewItem struct{ priority int }
type popReviewHeap struct {
	values []*popReviewItem
	probe  *popReviewItem
	trace  []string
}

func (h *popReviewHeap) Len() int {
	h.trace = append(h.trace, "Len")
	return len(h.values)
}
func (h *popReviewHeap) Swap(i, j int) {
	h.trace = append(h.trace, fmt.Sprintf("Swap:%d,%d", i, j))
	if h.probe == nil {
		h.values[i], h.values[j] = h.values[j], h.values[i]
	}
}
func (h *popReviewHeap) Less(i, j int) bool {
	h.trace = append(h.trace, fmt.Sprintf("Less:%d,%d", i, j))
	return h.values[i].priority < h.values[j].priority
}
func (h *popReviewHeap) Push(any) { panic("unexpected Push") }
func (h *popReviewHeap) Pop() any {
	h.trace = append(h.trace, "Pop")
	if h.probe != nil {
		return h.probe
	}
	last := len(h.values) - 1
	item := h.values[last]
	h.values = h.values[:last]
	return item
}

func includeFormerLast(h Interface) any {
	n := h.Len() - 1
	h.Swap(0, n)
	down(h, 0, n+1)
	return h.Pop()
}
func discardReturn(h Interface) any {
	Pop(h)
	return nil
}

func TestGoLeanHeapPopReview(t *testing.T) {
	root, middle, last, probe := &popReviewItem{1}, &popReviewItem{4}, &popReviewItem{8}, &popReviewItem{-9}
	fixtures := []struct {
		before      []*popReviewItem
		probe, want *popReviewItem
		trace       []string
		after       []*popReviewItem
	}{
		{nil, probe, probe, []string{"Len", "Swap:0,-1", "Pop"}, nil},
		{[]*popReviewItem{root}, nil, root, []string{"Len", "Swap:0,0", "Pop"}, nil},
		{[]*popReviewItem{root, middle, last}, nil, root,
			[]string{"Len", "Swap:0,2", "Less:1,0", "Swap:0,1", "Pop"}, []*popReviewItem{middle, last}},
	}
	for _, c := range []struct {
		name string
		pop  func(Interface) any
		want int
	}{
		{"Go", Pop, 0},
		{"Include the former last index during sift", includeFormerLast, 1},
		{"Discard Interface.Pop result", discardReturn, 3},
	} {
		bad := 0
		for _, f := range fixtures {
			h := &popReviewHeap{values: slices.Clone(f.before), probe: f.probe}
			got := c.pop(h)
			if got != f.want || !slices.Equal(h.trace, f.trace) || !slices.Equal(h.values, f.after) {
				bad++
				if c.want == 0 {
					t.Errorf("Pop got identity=%t trace=%v values=%v", got == f.want, h.trace, h.values)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
