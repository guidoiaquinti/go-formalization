package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanSatan(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_SATAN_CASES"))
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
		if math.IsNaN(x) || x < 0 {
			t.Fatalf("input %#016x is outside the nonnegative reduction domain", word)
		}
		bits := math.Float64bits(math.GoLeanSatan(x))
		if bits != model {
			t.Fatalf("input %#016x: private Go=%#016x Lean=%#016x", word, bits, model)
		}
		if (x == 0 || math.IsInf(x, 1)) && bits != reference {
			t.Fatalf("input %#016x: special result=%#016x want %#016x", word, bits, reference)
		}
		distance := max(bits, reference) - min(bits, reference)
		if distance > 4 {
			t.Fatalf("input %#016x: private Go=%#016x reference=%#016x encoding distance=%d", word, bits, reference, distance)
		}
		maxDistance = max(maxDistance, distance)
		count++
	}
	if count == 0 {
		t.Fatal("no private satan comparison rows")
	}
	t.Logf("private satan matched %d Lean rows; maximum reference distance %d encoded steps", count, maxDistance)
}
