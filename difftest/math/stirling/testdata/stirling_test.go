package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanStirling(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_STIRLING_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewReader(file)
	var count int
	var maximum float64
	for {
		var word, modelFirst, modelSecond, idealFirst, idealSecond uint64
		n, err := fmt.Fscan(input, &word, &modelFirst, &modelSecond, &idealFirst, &idealSecond)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		x := math.Float64frombits(word)
		first, second := math.GoLeanStirling(x)
		if got1, got2 := math.Float64bits(first), math.Float64bits(second); got1 != modelFirst || got2 != modelSecond {
			t.Fatalf("stirling(%#016x)=%#016x %#016x, Lean=%#016x %#016x", word, got1, got2, modelFirst, modelSecond)
		}
		if x > 200 {
			if !math.IsInf(first, 1) || second != 1 {
				t.Fatalf("stirling(%g)=%g %g, want infinity and one", x, first, second)
			}
		} else {
			wantFirst, wantSecond := math.Float64frombits(idealFirst), math.Float64frombits(idealSecond)
			error1, error2 := math.Abs(first/wantFirst-1), math.Abs(second/wantSecond-1)
			if math.IsNaN(error1) || math.IsNaN(error2) || error1 > 0x1p-43 || error2 > 0x1p-43 {
				t.Fatalf("stirling(%g)=%g %g, high-precision=%g %g, relative errors=%g %g",
					x, first, second, wantFirst, wantSecond, error1, error2)
			}
			maximum = max(maximum, error1, error2)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no Stirling cases")
	}
	t.Logf("private stirling matched %d Lean pairs; largest relative difference %.3g", count, maximum)
}
