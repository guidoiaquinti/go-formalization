package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanCbrt(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_CBRT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	count := 0
	for {
		var word, want uint64
		n, err := fmt.Fscan(reader, &word, &want)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got := math.Float64bits(math.GoLeanCbrt(math.Float64frombits(word)))
		if got != want {
			t.Fatalf("input %#016x: private Go=%#016x want=%#016x", word, got, want)
		}
		opposite := math.Float64bits(math.GoLeanCbrt(math.Float64frombits(word ^ (1 << 63))))
		if opposite != got^(1<<63) {
			t.Fatalf("input %#016x: private Go=%#016x opposite=%#016x", word, got, opposite)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private cbrt comparison rows")
	}
	t.Logf("private cbrt matched %d shared comparison rows", count)
}
