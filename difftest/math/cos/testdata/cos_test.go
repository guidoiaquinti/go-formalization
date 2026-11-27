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

func TestGoLeanCos(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_COS_CASES"))
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
		got := math.GoLeanCos(math.Float64frombits(word))
		bits := math.Float64bits(got)
		result := "NaN"
		if !math.IsNaN(got) {
			result = strconv.FormatUint(bits, 10)
		}
		if result != model {
			t.Fatalf("input %#016x: private Go=%s Lean=%s", word, result, model)
		}
		if !exact {
			opposite := math.Float64bits(math.GoLeanCos(math.Float64frombits(word ^ (1 << 63))))
			if opposite != bits {
				t.Fatalf("input %#016x: private cosine evenness failed", word)
			}
			want := math.Float64frombits(reference)
			delta := math.Abs(got - want)
			if math.IsNaN(delta) || delta > 0x1p-50 {
				t.Fatalf("input %#016x: private Go=%.17g reference=%.17g absolute error=%g", word, got, want, delta)
			}
			maxError = max(maxError, delta)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private cos comparison rows")
	}
	t.Logf("private cos matched %d Lean rows; maximum absolute reference difference %g", count, maxError)
}
