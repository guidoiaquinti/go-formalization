package math_test

import (
	"encoding/json"
	"math"
	"os"
	"testing"
)

func TestGoLeanReviewStirling(t *testing.T) {
	data, err := os.ReadFile(os.Getenv("GOLEAN_REVIEW_STIRLING_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Word  uint64
		First float64
		Gamma float64
		Mode  string
	}
	if err := json.Unmarshal(data, &cases); err != nil || len(cases) != 13 {
		t.Fatalf("invalid review fixtures: %v", err)
	}
	for _, c := range cases {
		first, second := math.GoLeanStirling(math.Float64frombits(c.Word))
		switch c.Mode {
		case "saturated":
			if !math.IsInf(first, 1) || second != 1 {
				t.Fatalf("input %#x: wanted infinity and one, got %g and %g", c.Word, first, second)
			}
			continue
		case "ordinary":
			if !(2.5 < second && second < 2.6) {
				t.Fatalf("input %#x: second factor %g outside ordinary correction range", c.Word, second)
			}
		case "split":
			if !(second > 1e50 && !math.IsInf(second, 0)) {
				t.Fatalf("input %#x: second factor %g outside split range", c.Word, second)
			}
		default:
			t.Fatalf("unknown fixture mode %q", c.Mode)
		}
		if !(math.Abs(first/c.First-1) <= 0x1p-40) {
			t.Fatalf("input %#x: first factor %g, expected near %g", c.Word, first, c.First)
		}
		if c.Gamma != 0 && !(math.Abs(first*second/c.Gamma-1) <= 0x1p-40) {
			t.Fatalf("input %#x: product %g, expected factorial %g", c.Word, first*second, c.Gamma)
		}
	}
}
