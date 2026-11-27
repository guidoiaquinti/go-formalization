package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"testing"
)

func TestGoLeanHypot(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_HYPOT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	observe := func(value float64) string {
		if math.IsNaN(value) {
			return "NaN"
		}
		return strconv.FormatUint(math.Float64bits(value), 10)
	}

	reader := bufio.NewReader(f)
	var count, maxDistance, overflowMismatches uint64
	for {
		var x, y uint64
		var model, ideal string
		var exact bool
		n, err := fmt.Fscan(reader, &x, &y, &model, &ideal, &exact)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		p, q := math.Float64frombits(x), math.Float64frombits(y)
		result := math.GoLeanHypot(p, q)
		got := observe(result)
		if swapped := observe(math.GoLeanHypot(q, p)); got != swapped || got != model {
			t.Fatalf("input %#016x, %#016x: private Go=%s swapped=%s Lean=%s", x, y, got, swapped, model)
		}
		if exact {
			if got != ideal {
				t.Fatalf("input %#016x, %#016x: private Go=%s oracle=%s", x, y, got, ideal)
			}
		} else {
			want, err := strconv.ParseUint(ideal, 10, 64)
			if err != nil {
				t.Fatal(err)
			}
			word := math.Float64bits(result)
			distance := max(word, want) - min(word, want)
			if distance > 2 {
				t.Fatalf("input %#016x, %#016x: private Go=%#016x oracle=%#016x distance=%d", x, y, word, want, distance)
			}
			maxDistance = max(maxDistance, distance)
			if math.IsInf(result, 1) != math.IsInf(math.Float64frombits(want), 1) {
				overflowMismatches++
			}
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private Hypot comparison rows")
	}
	t.Logf("private hypot %d inputs agree with Lean; maximum oracle distance %d word steps; %d finite/overflow classification differences",
		count, maxDistance, overflowMismatches)
}
