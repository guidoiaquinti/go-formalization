package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanLdexpGeneric(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_LDEXP_GENERIC_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := bufio.NewReader(f)
	for {
		var word, want uint64
		var exponent int
		n, err := fmt.Fscan(in, &word, &exponent, &want)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		got := math.Float64bits(math.GoLeanLdexpGeneric(math.Float64frombits(word), exponent))
		if got != want {
			t.Fatalf("ldexp(%#016x, %d) = %#016x, want %#016x", word, exponent, got, want)
		}
	}
}
