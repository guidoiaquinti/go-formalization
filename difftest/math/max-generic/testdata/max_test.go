package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanMaxGeneric(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_MAX_GENERIC_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := bufio.NewReader(f)
	for {
		var a, b, want uint64
		n, err := fmt.Fscan(in, &a, &b, &want)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		got := math.Float64bits(math.GoLeanMaxGeneric(math.Float64frombits(a), math.Float64frombits(b)))
		if got != want {
			t.Fatalf("max(%#016x, %#016x) has bits %#016x, want %#016x", a, b, got, want)
		}
	}
}
