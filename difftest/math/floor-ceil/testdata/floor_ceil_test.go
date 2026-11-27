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

func TestGoLeanFloorCeil(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_FLOOR_CEIL_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	input := bufio.NewReader(f)
	for {
		var word uint64
		var floor, ceil string
		n, err := fmt.Fscan(input, &word, &floor, &ceil)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		x := math.Float64frombits(word)
		if a, b := goLeanObserveFloorCeil(math.GoLeanFloor(x)), goLeanObserveFloorCeil(math.GoLeanCeil(x)); a != floor || b != ceil {
			t.Fatalf("input %#016x: Go=(%s %s) want=(%s %s)", word, a, b, floor, ceil)
		}
	}
}

func goLeanObserveFloorCeil(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(math.Float64bits(x), 10)
}
