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

func TestGoLeanArchDim(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ARCH_DIM_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewReader(file)
	count := 0
	arithmetic := runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64"
	for {
		var a, b uint64
		var expected string
		n, err := fmt.Fscanf(rows, "%d %d %q\n", &a, &b, &expected)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		x, y := math.Float64frombits(a), math.Float64frombits(b)
		maxOK, maximum, maxFailure := math.GoLeanArchMax(x, y)
		minOK, minimum, minFailure := math.GoLeanArchMin(x, y)
		if maxOK != arithmetic || minOK != arithmetic {
			t.Fatalf("archMax/Min availability mismatch: max=%t min=%t arithmetic=%t", maxOK, minOK, arithmetic)
		}
		var got string
		if arithmetic {
			if maxFailure != nil || minFailure != nil {
				t.Fatalf("archMax/Min(%#016x, %#016x) panic: max=%#v min=%#v", a, b, maxFailure, minFailure)
			}
			got = fmt.Sprintf("%d,%d", math.Float64bits(maximum), math.Float64bits(minimum))
		} else {
			maxMessage, maxString := maxFailure.(string)
			minMessage, minString := minFailure.(string)
			if !maxString || !minString || maxMessage != minMessage {
				t.Fatalf("archMax/Min(%#016x, %#016x) unequal string panics: max=%#v min=%#v", a, b, maxFailure, minFailure)
			}
			got = "panic:" + maxMessage
		}
		if got != expected {
			t.Fatalf("archMax/Min(%#016x, %#016x) got %q want %q", a, b, got, expected)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no archMax/Min cases")
	}
	t.Logf("archMax and archMin each matched %d direct native rows", count)
}
