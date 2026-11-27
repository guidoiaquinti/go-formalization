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

func TestGoLeanArchFloor(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_ARCH_FLOOR_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	input := bufio.NewReader(f)
	for {
		var word uint64
		var floor string
		n, err := fmt.Fscan(input, &word, &floor)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		x := math.Float64frombits(word)
		if got := goLeanObserveArchFloor(math.GoLeanArchFloor(x)); got != floor {
			t.Fatalf("input %#016x: Go=%s want=%s", word, got, floor)
		}
	}
}

func goLeanObserveArchFloor(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(math.Float64bits(x), 10)
}
