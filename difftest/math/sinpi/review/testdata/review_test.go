package math_test

import (
	"encoding/json"
	"math"
	"os"
	"testing"
)

func TestGoLeanReviewSinPi(t *testing.T) {
	data, err := os.ReadFile(os.Getenv("GOLEAN_REVIEW_SINPI_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct{ Input, Want, Steps uint64 }
	if err := json.Unmarshal(data, &cases); err != nil || len(cases) != 17 {
		t.Fatalf("invalid review fixtures: %v", err)
	}
	for _, c := range cases {
		got := math.Float64bits(math.GoLeanSinPi(math.Float64frombits(c.Input)))
		if max(got, c.Want)-min(got, c.Want) > c.Steps {
			t.Fatalf("input %#x: got %#x, want %#x within %d", c.Input, got, c.Want, c.Steps)
		}
	}
}
