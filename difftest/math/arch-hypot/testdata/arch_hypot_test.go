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

func TestGoLeanArchHypot(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ARCH_HYPOT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewReader(file)
	count := 0
	for {
		var x, y uint64
		var model string
		n, err := fmt.Fscanf(rows, "%d %d %q\n", &x, &y, &model)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		count++
		available, output, failure := math.GoLeanArchHypot(math.Float64frombits(x), math.Float64frombits(y))
		if runtime.GOARCH != "amd64" {
			message, ok := failure.(string)
			if available || !ok || message != "not implemented" || model != "panic:not implemented" {
				t.Fatalf("archHypot stub(%#016x, %#016x): available=%t panic=%#v Lean=%s", x, y, available, failure, model)
			}
			continue
		}
		got := strconv.FormatUint(math.Float64bits(output), 10)
		if !available || failure != nil || got != model {
			t.Fatalf("archHypot(%#016x, %#016x): available=%t panic=%#v native=%s Lean=%s", x, y, available, failure, got, model)
		}
	}
	if count == 0 {
		t.Fatal("no archHypot cases")
	}
	t.Logf("archHypot matched %d direct native rows", count)
}
