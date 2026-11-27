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

func TestGoLeanLog2(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_LOG2_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var count int
	var maxScaledError float64
	var maxStepsAwayFromZero uint64
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
		got := math.GoLeanLog2(math.Float64frombits(word))
		bits := math.Float64bits(got)
		result := "NaN"
		if !math.IsNaN(got) {
			result = strconv.FormatUint(bits, 10)
		}
		if result != model {
			t.Fatalf("input %#016x: private Go=%s Lean=%s", word, result, model)
		}
		if !exact {
			want := math.Float64frombits(reference)
			scaledError := math.Abs(got-want) / max(1, math.Abs(want))
			if math.IsNaN(scaledError) || scaledError > 0x1p-50 {
				t.Fatalf("input %#016x: private Go=%.17g reference=%.17g scaled error=%g", word, got, want, scaledError)
			}
			maxScaledError = max(maxScaledError, scaledError)
			if math.Abs(want) >= 1 {
				distance := max(bits, reference) - min(bits, reference)
				if distance > 2 {
					t.Fatalf("input %#016x: private Go=%#016x reference=%#016x distance=%d", word, bits, reference, distance)
				}
				maxStepsAwayFromZero = max(maxStepsAwayFromZero, distance)
			}
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private log2 comparison rows")
	}
	t.Logf("private log2 matched %d Lean rows; max scaled error %g; max encoded steps where |result| >= 1 is %d",
		count, maxScaledError, maxStepsAwayFromZero)
}
