package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanFmaSplit(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_FMA_SPLIT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewReader(file)
	for {
		var word, wantSignificand uint64
		var wantSign uint32
		var wantExponent int32
		n, err := fmt.Fscan(input, &word, &wantSign, &wantExponent, &wantSignificand)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		sign, exponent, significand := math.GoLeanFmaSplit(word)
		if sign != wantSign || exponent != wantExponent || significand != wantSignificand {
			t.Fatalf("split(%#016x) = (%d, %d, %#x), want (%d, %d, %#x)", word, sign,
				exponent, significand, wantSign, wantExponent, wantSignificand)
		}
	}
}
