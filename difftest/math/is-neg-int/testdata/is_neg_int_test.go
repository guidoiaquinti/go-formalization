package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanIsNegInt(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_IS_NEG_INT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewReader(file)
	for {
		var word, want uint64
		n, err := fmt.Fscan(input, &word, &want)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		if got := math.GoLeanIsNegInt(word); got != (want == 1) {
			t.Fatalf("isNegInt(%#016x) = %t, want %d", word, got, want)
		}
	}
}
