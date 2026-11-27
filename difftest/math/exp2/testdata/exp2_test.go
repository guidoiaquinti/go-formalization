package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanExp2(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_EXP2_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var count int
	var maxDistance uint64
	for {
		var word, model, reference uint64
		var exact bool
		n, err := fmt.Fscan(reader, &word, &model, &reference, &exact)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got := math.Float64bits(math.GoLeanExp2(math.Float64frombits(word)))
		if got != model {
			t.Fatalf("input %#016x: private Go=%#016x Lean=%#016x", word, got, model)
		}
		distance := max(got, reference) - min(got, reference)
		if distance > 1 || exact && distance != 0 {
			t.Fatalf("input %#016x: private Go=%#016x reference=%#016x distance=%d exact=%t",
				word, got, reference, distance, exact)
		}
		maxDistance = max(maxDistance, distance)
		count++
	}
	if count == 0 {
		t.Fatal("no private exp2 comparison rows")
	}
	t.Logf("private exp2 matched %d Lean rows; maximum reference distance %d encoded steps", count, maxDistance)
}
