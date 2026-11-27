package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanXatan(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_XATAN_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var count int
	var maxDistance uint64
	for {
		var word, model, reference uint64
		n, err := fmt.Fscan(reader, &word, &model, &reference)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		x := math.Float64frombits(word)
		if !(math.Abs(x) <= 0.66) {
			t.Fatalf("input %#016x is outside the signed rational interval", word)
		}
		bits := math.Float64bits(math.GoLeanXatan(x))
		if bits != model {
			t.Fatalf("input %#016x: private Go=%#016x Lean=%#016x", word, bits, model)
		}
		if x != 0 && math.Float64bits(math.GoLeanXatan(-x)) != bits^(1<<63) {
			t.Fatalf("input %#016x: xatan sign reflection failed", word)
		}
		distance := max(bits, reference) - min(bits, reference)
		if distance > 4 {
			t.Fatalf("input %#016x: private Go=%#016x reference=%#016x encoding distance=%d", word, bits, reference, distance)
		}
		maxDistance = max(maxDistance, distance)
		count++
	}
	if count == 0 {
		t.Fatal("no private xatan comparison rows")
	}
	t.Logf("private xatan matched %d Lean rows; maximum reference distance %d encoded steps", count, maxDistance)
}
