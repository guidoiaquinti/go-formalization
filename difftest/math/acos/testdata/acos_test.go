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

func TestGoLeanAcos(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ACOS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var count int
	var maxPositiveError, maxNegativeError float64
	for {
		var word, reference uint64
		var model, kind string
		n, err := fmt.Fscan(reader, &word, &model, &reference, &kind)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		x := math.Float64frombits(word)
		got := math.GoLeanAcos(x)
		bits := math.Float64bits(got)
		result := "NaN"
		if !math.IsNaN(got) {
			result = strconv.FormatUint(bits, 10)
		}
		if result != model {
			t.Fatalf("input %#016x: private Go=%s Lean=%s", word, result, model)
		}
		switch kind {
		case "exact":
			if !math.IsNaN(got) && bits != reference {
				t.Fatalf("input %#016x: private special Go=%#016x want %#016x", word, bits, reference)
			}
		case "numeric", "upper":
			staged := math.Float64bits(math.Pi/2 - math.Asin(x))
			if bits != staged {
				t.Fatalf("input %#016x: private Go=%#016x, rounded subtraction=%#016x", word, bits, staged)
			}
			want := math.Float64frombits(reference)
			delta := math.Abs(got - want)
			bound := float64(0x1p-49)
			if kind == "upper" {
				bound = 0x1p-38
			}
			if math.IsNaN(delta) || delta > bound {
				t.Fatalf("input %#016x: private Go=%.17g reference=%.17g absolute difference=%g", word, got, want, delta)
			}
			if x < 0 {
				maxNegativeError = max(maxNegativeError, delta)
			} else {
				maxPositiveError = max(maxPositiveError, delta)
			}
		default:
			t.Fatalf("input %#016x: unknown fixture kind %q", word, kind)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private acos comparison rows")
	}
	t.Logf("private acos matched %d Lean rows; maximum absolute difference positive inputs=%g negative inputs=%g",
		count, maxPositiveError, maxNegativeError)
}
