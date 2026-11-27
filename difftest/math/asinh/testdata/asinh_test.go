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

func TestGoLeanAsinh(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ASINH_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var count int
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
		got := math.GoLeanAsinh(math.Float64frombits(word))
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
			if result != reference {
				t.Fatalf("input %#016x: private Go=%s exact reference=%s", word, result, reference)
			}
		case "numeric":
			want, err := strconv.ParseUint(reference, 10, 64)
			if err != nil {
				t.Fatal(err)
			}
			distance := max(bits, want) - min(bits, want)
			if math.IsNaN(got) || distance > 4 {
				t.Fatalf("input %#016x: private Go=%#016x reference=%#016x distance=%d", word, bits, want, distance)
			}
			maximum = max(maximum, distance)
		default:
			t.Fatalf("input %#016x: unknown fixture kind %q", word, kind)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private asinh comparison rows")
	}
	t.Logf("private asinh matched %d Lean rows; numerical reference maximum %d binary64 steps", count, maximum)
}
