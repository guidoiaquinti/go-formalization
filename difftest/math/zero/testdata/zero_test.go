package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanFmaZero(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_FMA_ZERO_CASES"))
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
		if got := math.GoLeanFmaZero(word); got != want {
			t.Fatalf("zero(%#016x) = %d, want %d", word, got, want)
		}
	}
}
