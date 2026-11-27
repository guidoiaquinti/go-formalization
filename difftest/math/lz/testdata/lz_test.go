package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanLz(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_LZ_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewReader(file)
	for {
		var high, low uint64
		var want int32
		n, err := fmt.Fscan(input, &high, &low, &want)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		if got := math.GoLeanLz(high, low); got != want {
			t.Fatalf("lz(%#016x, %#016x) = %d, want %d", high, low, got, want)
		}
	}
}
