package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"testing"
)

func TestGoLeanArchExp(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ARCH_EXP_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewReader(file)
	count, earlyOverflows := 0, 0
	var maximum uint64
	for {
		var word, model, ideal uint64
		var exact bool
		n, err := fmt.Fscan(rows, &word, &model, &ideal, &exact)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		x := math.Float64frombits(word)
		got := math.Float64bits(math.GoLeanArchExp(x))
		if got != model {
			t.Fatalf("archExp(%#016x): native=%#016x Lean=%#016x", word, got, model)
		}
		count++
		if exact && got != ideal {
			t.Fatalf("archExp(%#016x): native=%#016x exact outcome=%#016x", word, got, ideal)
		}
		if runtime.GOARCH == "amd64" && got == 0x7ff0000000000000 && ideal < got &&
			math.RoundToEven(x*math.Log2E) == 1024 {
			earlyOverflows++
			continue
		}
		distance := max(got, ideal) - min(got, ideal)
		if distance > 2 {
			t.Fatalf("archExp(%#016x): native=%#016x independent=%#016x differ by %d steps", word, got, ideal, distance)
		}
		maximum = max(maximum, distance)
	}
	if count == 0 {
		t.Fatal("no archExp cases")
	}
	t.Logf("archExp matched %d direct native rows; independent maximum %d steps, %d amd64 early overflows", count, maximum, earlyOverflows)
}
