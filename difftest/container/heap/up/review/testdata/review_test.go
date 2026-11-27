package heap

import (
	"math"
	"slices"
	"testing"
)

type upReviewCall struct {
	kind          byte
	first, second int
}

type upReviewInterface struct {
	decisions uint
	trace     []upReviewCall
}

func (h *upReviewInterface) Len() int { panic("unexpected Len") }
func (h *upReviewInterface) Push(any) { panic("unexpected Push") }
func (h *upReviewInterface) Pop() any { panic("unexpected Pop") }

func (h *upReviewInterface) Less(i, j int) bool {
	h.trace = append(h.trace, upReviewCall{'L', i, j})
	answer := h.decisions&1 != 0
	h.decisions >>= 1
	return answer
}

func (h *upReviewInterface) Swap(i, j int) {
	h.trace = append(h.trace, upReviewCall{'S', i, j})
}

func upReviewWrongParent(h Interface, j int) {
	for j != 0 {
		i := j / 2
		if !h.Less(j, i) {
			return
		}
		h.Swap(i, j)
		j = i
	}
}

func upReviewRootCallback(h Interface, j int) {
	for {
		i := (j - 1) / 2
		if !h.Less(j, i) || i == j {
			return
		}
		h.Swap(i, j)
		j = i
	}
}

func TestGoLeanHeapUpReview(t *testing.T) {
	fixtures := []struct {
		child     int
		decisions uint
		want      []upReviewCall
	}{
		{0, 1, nil},
		{1, 1, []upReviewCall{{'L', 1, 0}, {'S', 0, 1}}},
		{2, 0, []upReviewCall{{'L', 2, 0}}},
		{6, 1, []upReviewCall{{'L', 6, 2}, {'S', 2, 6}, {'L', 2, 0}}},
		{14, 7, []upReviewCall{{'L', 14, 6}, {'S', 6, 14}, {'L', 6, 2}, {'S', 2, 6}, {'L', 2, 0}, {'S', 0, 2}}},
		{math.MaxInt, 0, []upReviewCall{{'L', math.MaxInt, (math.MaxInt - 1) / 2}}},
	}
	for _, c := range []struct {
		name string
		up   func(Interface, int)
		want int
	}{
		{"Go", up, 0},
		{"Use child divided by two", upReviewWrongParent, 3},
		{"Compare before checking for root", upReviewRootCallback, 3},
	} {
		mismatches := 0
		for _, f := range fixtures {
			h := &upReviewInterface{decisions: f.decisions}
			c.up(h, f.child)
			if !slices.Equal(h.trace, f.want) {
				mismatches++
				if c.want == 0 {
					t.Errorf("up(%d) trace=%v, want=%v", f.child, h.trace, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, mismatches, len(fixtures))
		if mismatches != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
