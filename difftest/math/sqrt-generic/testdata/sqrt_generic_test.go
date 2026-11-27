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

func TestGoLeanSqrtGeneric(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_SQRT_GENERIC_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewReader(file)
	for {
		var word uint64
		var want string
		n, err := fmt.Fscan(input, &word, &want)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		value := math.GoLeanSqrtGeneric(word)
		got := "NaN"
		if !math.IsNaN(value) {
			got = strconv.FormatUint(math.Float64bits(value), 10)
		}
		if got != want {
			t.Fatalf("sqrt(%#016x) = %s, want %s", word, got, want)
		}
	}
}
