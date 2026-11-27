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

func TestGoLeanArchLog(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ARCH_LOG_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewReader(file)
	count := 0
	for {
		var word uint64
		var model string
		n, err := fmt.Fscanf(rows, "%d %q\n", &word, &model)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		count++
		available, output, failure := math.GoLeanArchLog(math.Float64frombits(word))
		if runtime.GOARCH != "amd64" {
			message, ok := failure.(string)
			if available || !ok || message != "not implemented" || model != "panic:not implemented" {
				t.Fatalf("archLog stub(%#016x): available=%t panic=%#v Lean=%s", word, available, failure, model)
			}
			continue
		}
		got := "NaN"
		if !math.IsNaN(output) {
			got = strconv.FormatUint(math.Float64bits(output), 10)
		}
		if !available || failure != nil || got != model {
			t.Fatalf("archLog(%#016x): available=%t panic=%#v native=%s Lean=%s", word, available, failure, got, model)
		}
	}
	if count == 0 {
		t.Fatal("no archLog cases")
	}
	t.Logf("archLog matched %d direct native rows", count)
}
