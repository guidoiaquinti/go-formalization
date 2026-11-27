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

func TestGoLeanAsin(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ASIN_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var count int
	var maxDistance uint64
	var maxUpperError float64
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
		got := math.GoLeanAsin(x)
		bits := math.Float64bits(got)
		result := "NaN"
		if !math.IsNaN(got) {
			result = strconv.FormatUint(bits, 10)
		}
		if result != model {
			t.Fatalf("input %#016x: private Go=%s Lean=%s", word, result, model)
		}
		if kind != "exact" && math.Float64bits(math.GoLeanAsin(-x)) != bits^(1<<63) {
			t.Fatalf("input %#016x: private asin sign reflection failed", word)
		}
		switch kind {
		case "exact":
			if !math.IsNaN(got) && bits != reference {
				t.Fatalf("input %#016x: private special result=%#016x want %#016x", word, bits, reference)
			}
		case "numeric":
			distance := max(bits, reference) - min(bits, reference)
			if distance > 4 {
				t.Fatalf("input %#016x: private Go=%#016x reference=%#016x encoding distance=%d", word, bits, reference, distance)
			}
			maxDistance = max(maxDistance, distance)
		case "upper":
			want := math.Float64frombits(reference)
			delta := math.Abs(got - want)
			if math.IsNaN(delta) || delta > 0x1p-38 {
				t.Fatalf("input %#016x: private Go=%.17g reference=%.17g absolute error=%g", word, got, want, delta)
			}
			maxUpperError = max(maxUpperError, delta)
		default:
			t.Fatalf("input %#016x: unknown fixture kind %q", word, kind)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private asin comparison rows")
	}
	t.Logf("private asin matched %d Lean rows; maximum distance %d encoded steps; separately rounded upper absolute error %g",
		count, maxDistance, maxUpperError)
}
