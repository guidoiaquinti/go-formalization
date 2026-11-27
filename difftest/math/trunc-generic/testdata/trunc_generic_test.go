package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanTruncGeneric(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_TRUNC_GENERIC_CASES"))
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
		if got := math.GoLeanTruncGeneric(word); got != want {
			t.Fatalf("trunc(%#016x) = %#016x, want %#016x", word, got, want)
		}
	}
}
