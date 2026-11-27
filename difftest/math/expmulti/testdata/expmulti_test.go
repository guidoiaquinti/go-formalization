package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanExpmulti(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_EXPMULTI_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewReader(file)
	var count int
	var maximum uint64
	for {
		var hi, lo, model, ideal uint64
		var k int
		n, err := fmt.Fscan(rows, &hi, &lo, &k, &model, &ideal)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got := math.Float64bits(math.GoLeanExpmulti(math.Float64frombits(hi), math.Float64frombits(lo), k))
		if got != model {
			t.Fatalf("expmulti(%#016x,%#016x,%d)=%#016x, Lean=%#016x", hi, lo, k, got, model)
		}
		distance := max(got, ideal) - min(got, ideal)
		if distance > 3 {
			t.Fatalf("expmulti(%#016x,%#016x,%d)=%#016x, 256-bit reference=%#016x (%d steps)",
				hi, lo, k, got, ideal, distance)
		}
		maximum = max(maximum, distance)
		count++
	}
	if count == 0 {
		t.Fatal("no expmulti cases")
	}
	t.Logf("private expmulti matched %d Lean rows; reference maximum %d binary64 steps", count, maximum)
}
