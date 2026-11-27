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

func TestGoLeanAtan(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ATAN_CASES"))
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
		got := math.GoLeanAtan(math.Float64frombits(word))
		bits := math.Float64bits(got)
		result := "NaN"
		if !math.IsNaN(got) {
			result = strconv.FormatUint(bits, 10)
		}
		if result != model {
			t.Fatalf("input %#016x: private Go=%s Lean=%s", word, result, model)
		}
		if !exact {
			opposite := math.Float64bits(math.GoLeanAtan(math.Float64frombits(word ^ (1 << 63))))
			if opposite != bits^(1<<63) {
				t.Fatalf("input %#016x: private atan sign reflection failed", word)
			}
			distance := max(bits, reference) - min(bits, reference)
			if distance > 4 {
				t.Fatalf("input %#016x: private Go=%#016x reference=%#016x encoding distance=%d", word, bits, reference, distance)
			}
			maxDistance = max(maxDistance, distance)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private atan comparison rows")
	}
	t.Logf("private atan matched %d Lean rows; maximum reference distance %d encoded steps", count, maxDistance)
}
