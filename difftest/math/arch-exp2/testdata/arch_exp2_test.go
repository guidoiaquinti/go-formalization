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

func TestGoLeanArchExp2(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ARCH_EXP2_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewReader(file)
	count, earlyUnderflows := 0, 0
	var maximum uint64
	for {
		var word, ideal uint64
		var model string
		var exact bool
		n, err := fmt.Fscanf(rows, "%d %q %d %t\n", &word, &model, &ideal, &exact)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		count++
		x := math.Float64frombits(word)
		available, output, failure := math.GoLeanArchExp2(x)
		if runtime.GOARCH != "arm64" {
			message, ok := failure.(string)
			if available || !ok || message != "not implemented" || model != "panic:not implemented" {
				t.Fatalf("archExp2 stub(%#016x): available=%t panic=%#v Lean=%s", word, available, failure, model)
			}
			continue
		}
		got := math.Float64bits(output)
		if !available || failure != nil || strconv.FormatUint(got, 10) != model {
			t.Fatalf("archExp2(%#016x): available=%t panic=%#v native=%#016x Lean=%s", word, available, failure, got, model)
		}
		if exact && got != ideal {
			t.Fatalf("archExp2(%#016x): native=%#016x exact reference=%#016x", word, got, ideal)
		}
		if -1075 < x && x < -1074 {
			earlyUnderflows++
		}
		distance := max(got, ideal) - min(got, ideal)
		if distance > 1 {
			t.Fatalf("archExp2(%#016x): native=%#016x independent=%#016x differ by %d steps", word, got, ideal, distance)
		}
		maximum = max(maximum, distance)
	}
	if count == 0 {
		t.Fatal("no archExp2 cases")
	}
	t.Logf("archExp2 matched %d direct native rows; independent maximum %d steps, %d early underflows", count, maximum, earlyUnderflows)
}
