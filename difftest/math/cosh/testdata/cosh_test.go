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

func TestGoLeanCosh(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_COSH_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var count, premature int
	var maximum uint64
	for {
		var word uint64
		var model, reference, kind string
		n, err := fmt.Fscan(reader, &word, &model, &reference, &kind)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		x := math.Float64frombits(word)
		got := math.GoLeanCosh(x)
		bits := math.Float64bits(got)
		result := "NaN"
		if !math.IsNaN(got) {
			if reflected := math.GoLeanCosh(-x); math.Float64bits(reflected) != bits {
				t.Fatalf("input %#016x: private Go=%g reflected=%g", word, got, reflected)
			}
			result = strconv.FormatUint(bits, 10)
		}
		if result != model {
			t.Fatalf("input %#016x: private Go=%s Lean=%s", word, result, model)
		}
		switch kind {
		case "exact":
			if result != reference {
				t.Fatalf("input %#016x: private Go=%s exact reference=%s", word, result, reference)
			}
		case "staged":
			want, err := strconv.ParseUint(reference, 10, 64)
			if err != nil {
				t.Fatal(err)
			}
			expected := 0.5 * math.Exp(math.Abs(x))
			if math.Abs(x) <= 21 || !math.IsInf(got, 0) || math.IsInf(math.Float64frombits(want), 0) ||
				bits != math.Float64bits(expected) {
				t.Fatalf("input %#016x: appearance of staged overflow differs: private Go=%g, stage=%g, reference=%#016x", word, got, expected, want)
			}
			premature++
		case "numeric":
			want, err := strconv.ParseUint(reference, 10, 64)
			if err != nil {
				t.Fatal(err)
			}
			distance := max(bits, want) - min(bits, want)
			if math.IsNaN(got) || math.IsInf(got, 0) != math.IsInf(math.Float64frombits(want), 0) || distance > 5 {
				t.Fatalf("input %#016x: private Go=%#016x reference=%#016x distance=%d", word, bits, want, distance)
			}
			maximum = max(maximum, distance)
		default:
			t.Fatalf("input %#016x: unknown fixture kind %q", word, kind)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private cosh comparison rows")
	}
	t.Logf("private cosh matched %d Lean rows; numerical reference maximum %d binary64 steps, %d staged overflows with finite reference", count, maximum, premature)
}
