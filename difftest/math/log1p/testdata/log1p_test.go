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

func TestGoLeanLog1p(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_LOG1P_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var count int
	var maxDistance uint64
	for {
		var word, reference uint64
		var model string
		var exact bool
		n, err := fmt.Fscan(reader, &word, &model, &reference, &exact)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got := math.GoLeanLog1p(math.Float64frombits(word))
		bits := math.Float64bits(got)
		result := "NaN"
		if !math.IsNaN(got) {
			result = strconv.FormatUint(bits, 10)
		}
		if result != model {
			t.Fatalf("input %#016x: private Go=%s Lean=%s", word, result, model)
		}
		if !exact {
			distance := max(bits, reference) - min(bits, reference)
			if distance > 1 {
				t.Fatalf("input %#016x: private Go=%#016x reference=%#016x distance=%d",
					word, bits, reference, distance)
			}
			maxDistance = max(maxDistance, distance)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private log1p comparison rows")
	}
	t.Logf("private log1p matched %d Lean rows; maximum reference distance %d encoded steps", count, maxDistance)
}
