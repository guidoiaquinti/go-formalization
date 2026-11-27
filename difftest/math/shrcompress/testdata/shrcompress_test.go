package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanShrCompress(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_SHRCOMPRESS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewReader(file)
	for {
		var high, low, wantHigh, wantLow uint64
		var shift uint
		n, err := fmt.Fscan(input, &high, &low, &shift, &wantHigh, &wantLow)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		if outHigh, outLow := math.GoLeanShrCompress(high, low, shift); outHigh != wantHigh || outLow != wantLow {
			t.Fatalf("shrcompress(%#016x, %#016x, %d) = (%#016x, %#016x), want (%#016x, %#016x)",
				high, low, shift, outHigh, outLow, wantHigh, wantLow)
		}
	}
}
