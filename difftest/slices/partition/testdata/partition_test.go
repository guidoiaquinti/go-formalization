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

type partitionPanic int

func TestGoLeanPartition(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SLICES_PARTITION_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Words                                 []int64
			Kind, A, B, Pivot, Mode, FailAt, Mask int
			Want                                  string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing partition fixtures")
			}
			t.Logf("four pinned partition entries matched %d paired rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		ordered := slices.Clone(c.Words)
		first := capturePartition(func() string {
			if c.Kind == 1 {
				return fmt.Sprintf("ok:%d", slices.GoLeanPartitionEqualOrdered(ordered, c.A, c.B, c.Pivot))
			}
			pivot, already := slices.GoLeanPartitionOrdered(ordered, c.A, c.B, c.Pivot)
			return fmt.Sprintf("ok:%d,%t", pivot, already)
		})
		data := slices.Clone(c.Words)
		var calls []string
		var compare func(int64, int64) int
		if c.Mode != 0 {
			compare = func(x, y int64) int {
				call := len(calls)
				calls = append(calls, fmt.Sprintf("%d,%d", x, y))
				if call == c.FailAt {
					panic(partitionPanic(call))
				}
				if c.Mode == 1 {
					return cmp.Compare(x, y)
				}
				if c.Mask&(1<<(call%12)) != 0 {
					return -1
				}
				return 0
			}
		}
		second := capturePartition(func() string {
			if c.Kind == 1 {
				return fmt.Sprintf("ok:%d", slices.GoLeanPartitionEqualCmpFunc(data, c.A, c.B, c.Pivot, compare))
			}
			pivot, already := slices.GoLeanPartitionCmpFunc(data, c.A, c.B, c.Pivot, compare)
			return fmt.Sprintf("ok:%d,%t", pivot, already)
		})
		if got := wire(first, ordered, nil) + "|" + wire(second, data, calls); got != c.Want {
			t.Fatalf("row %d got %q want %q", i, got, c.Want)
		}
	}
}

func wire(status string, data []int64, calls []string) string {
	values := make([]string, len(data))
	for i, v := range data {
		values[i] = fmt.Sprint(v)
	}
	return status + "/" + strings.Join(values, ",") + "/" + strings.Join(calls, ";")
}

func capturePartition(call func() string) (out string) {
	defer func() {
		switch failure := recover().(type) {
		case nil:
		case partitionPanic:
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
	return call()
}
