package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"strconv"
	"testing"
)

func TestGoLeanPone(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_PONE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewReader(file)
	var count, panics int
	var maximum uint64
	for {
		var input uint64
		var model, ideal string
		n, err := fmt.Fscan(rows, &input, &model, &ideal)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got, panicked := callPone(math.Float64frombits(input))
		if model == "panic" || ideal == "panic" || panicked {
			if model != "panic" || ideal != "panic" || !panicked {
				t.Fatalf("pone(%#016x)=%#016x panic=%t, Lean=%s reference=%s", input, got, panicked, model, ideal)
			}
			panics++
		} else {
			if strconv.FormatUint(got, 10) != model {
				t.Fatalf("pone(%#016x)=%#016x, Lean=%s", input, got, model)
			}
			want, err := strconv.ParseUint(ideal, 10, 64)
			if err != nil {
				t.Fatal(err)
			}
			distance := max(got, want) - min(got, want)
			if distance > 3 {
				t.Fatalf("pone(%#016x)=%#016x, independent=%#016x (%d steps)", input, got, want, distance)
			}
			maximum = max(maximum, distance)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no pone cases")
	}
	t.Logf("private pone matched %d Lean rows (%d panics); reference maximum %d binary64 steps", count, panics, maximum)
}

func callPone(x float64) (word uint64, panicked bool) {
	defer func() {
		if err := recover(); err != nil {
			fault, ok := err.(runtime.Error)
			if !ok || fault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(err)
			}
			panicked = true
		}
	}()
	return math.Float64bits(math.GoLeanPone(x)), false
}
