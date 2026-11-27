package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanFrexpGeneric(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_FREXP_GENERIC_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := bufio.NewReader(f)
	for {
		var word, wantFraction uint64
		var wantExponent int
		n, err := fmt.Fscan(in, &word, &wantFraction, &wantExponent)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		fraction, exponent := math.GoLeanFrexpGeneric(math.Float64frombits(word))
		bits := math.Float64bits(fraction)
		if bits != wantFraction || exponent != wantExponent {
			t.Fatalf("frexp(%#016x) = (%#016x, %d), want (%#016x, %d)", word, bits, exponent, wantFraction, wantExponent)
		}
	}
}
