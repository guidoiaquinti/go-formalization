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

func TestGoLeanTan(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_TAN_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var count int
	var maxError float64
	for {
		var word, reference uint64
		var model string
		var exact bool
		n, err := fmt.Fscan(reader, &word, &model, &reference, &exact)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got := math.GoLeanTan(math.Float64frombits(word))
		bits := math.Float64bits(got)
		result := "NaN"
		if !math.IsNaN(got) {
			result = strconv.FormatUint(bits, 10)
		}
		if result != model {
			t.Fatalf("input %#016x: private Go=%s Lean=%s", word, result, model)
		}
		if !exact {
			opposite := math.Float64bits(math.GoLeanTan(math.Float64frombits(word ^ (1 << 63))))
			if opposite != bits^(1<<63) {
				t.Fatalf("input %#016x: private tangent sign reflection failed", word)
			}
			want := math.Float64frombits(reference)
			delta := math.Abs(got-want) / max(1, math.Abs(want))
			if math.IsNaN(delta) || delta > 0x1p-35 || math.Abs(want) <= 2 && delta > 0x1p-48 {
				t.Fatalf("input %#016x: private Go=%.17g reference=%.17g scaled absolute error=%g", word, got, want, delta)
			}
			maxError = max(maxError, delta)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private tan comparison rows")
	}
	t.Logf("private tan matched %d Lean rows; maximum scaled absolute reference error %g", count, maxError)
}
