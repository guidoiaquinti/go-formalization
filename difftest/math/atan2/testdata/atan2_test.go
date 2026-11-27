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

func TestGoLeanAtan2(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ATAN2_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var count, underflows int
	var maxDistance uint64
	for {
		var yWord, xWord, reference uint64
		var model, kind string
		n, err := fmt.Fscan(reader, &yWord, &xWord, &model, &reference, &kind)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		y, x := math.Float64frombits(yWord), math.Float64frombits(xWord)
		got := math.GoLeanAtan2(y, x)
		bits := math.Float64bits(got)
		result := "NaN"
		if !math.IsNaN(got) {
			result = strconv.FormatUint(bits, 10)
		}
		if result != model {
			t.Fatalf("inputs %#016x, %#016x: private Go=%s Lean=%s", yWord, xWord, result, model)
		}
		switch kind {
		case "exact":
			if !math.IsNaN(got) && bits != reference {
				t.Fatalf("inputs %#016x, %#016x: private special result=%#016x want=%#016x",
					yWord, xWord, bits, reference)
			}
		case "numeric":
			distance := max(bits, reference) - min(bits, reference)
			if distance > 4 {
				t.Fatalf("inputs %#016x, %#016x: private Go=%#016x reference=%#016x encoding distance=%d",
					yWord, xWord, bits, reference, distance)
			}
			maxDistance = max(maxDistance, distance)
		case "underflow":
			if !(y < 0 && x < 0) || math.Float64bits(y/x) != 0 || got != math.Pi || reference != math.Float64bits(-math.Pi) {
				t.Fatalf("inputs %#016x, %#016x: private lower-left zero case: ratio=%g Go=%g reference=%g",
					yWord, xWord, y/x, got, math.Float64frombits(reference))
			}
			underflows++
		default:
			t.Fatalf("inputs %#016x, %#016x: unknown fixture kind %q", yWord, xWord, kind)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private atan2 comparison rows")
	}
	t.Logf("private atan2 matched %d Lean rows, including %d lower-left zero cases; maximum distance elsewhere %d encoded steps",
		count, underflows, maxDistance)
}
