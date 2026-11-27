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

func TestGoLeanMod(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_MOD_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		var x, y uint64
		var want string
		n, err := fmt.Fscan(reader, &x, &y, &want)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		result := math.GoLeanMod(math.Float64frombits(x), math.Float64frombits(y))
		got := "NaN"
		if !math.IsNaN(result) {
			got = strconv.FormatUint(math.Float64bits(result), 10)
		}
		if got != want {
			t.Fatalf("input %#016x %% %#016x: private Go=%s want=%s", x, y, got, want)
		}
	}
}
