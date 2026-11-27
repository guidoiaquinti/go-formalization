package slices_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"strings"
	"testing"
)

type pivotPanic int

func TestGoLeanChoosePivot(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SLICES_PIVOT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Words                    []int64
			A, B, Mode, FailAt, Mask int
			Want                     string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing choosePivot fixtures")
			}
			t.Logf("both pinned choosePivot entries matched %d rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		ordered := slices.Clone(c.Words)
		first := capturePivot(func() (int, int) { return slices.GoLeanChoosePivotOrdered(ordered, c.A, c.B) })
		data := slices.Clone(c.Words)
		var calls []string
		var compare func(int64, int64) int
		if c.Mode != 0 {
			compare = func(x, y int64) int {
				call := len(calls)
				calls = append(calls, fmt.Sprintf("%d,%d", x, y))
				if call == c.FailAt {
					panic(pivotPanic(call))
				}
				if c.Mode == 1 {
					return cmp.Compare(x, y)
				}
				if c.Mask&(1<<call) != 0 {
					return -1
				}
				return 0
			}
		}
		second := capturePivot(func() (int, int) { return slices.GoLeanChoosePivotCmpFunc(data, c.A, c.B, compare) })
		if !slices.Equal(ordered, c.Words) || !slices.Equal(data, c.Words) {
			t.Fatalf("row %d mutated its input", i)
		}
		if got := first + "/|" + second + "/" + strings.Join(calls, ";"); got != c.Want {
			t.Fatalf("row %d got %q want %q", i, got, c.Want)
		}
	}
}

func capturePivot(call func() (int, int)) (out string) {
	defer func() {
		switch failure := recover().(type) {
		case nil:
		case pivotPanic:
			out = fmt.Sprintf("panic:%d", failure)
		case runtime.Error:
			if failure.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			out = "nil"
		default:
			panic(failure)
		}
	}()
	pivot, hint := call()
	return fmt.Sprintf("ok:%d,%d", pivot, hint)
}
