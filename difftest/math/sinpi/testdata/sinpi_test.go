package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanSinPi(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_SINPI_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewReader(file)
	var count int
	var maximum uint64
	for {
		var input, model, ideal uint64
		n, err := fmt.Fscan(rows, &input, &model, &ideal)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got := math.Float64bits(math.GoLeanSinPi(math.Float64frombits(input)))
		if got != model {
			t.Fatalf("sinPi(%#016x)=%#016x, Lean=%#016x", input, got, model)
		}
		distance := max(got, ideal) - min(got, ideal)
		if distance > 4 || (ideal<<1 == 0 && got != ideal) {
			t.Fatalf("sinPi(%#016x)=%#016x, independent reference=%#016x (%d steps)", input, got, ideal, distance)
		}
		maximum = max(maximum, distance)
		count++
	}
	if count == 0 {
		t.Fatal("no sinPi cases")
	}
	t.Logf("private sinPi matched %d Lean rows, independent reference maximum %d binary64 steps", count, maximum)
}
