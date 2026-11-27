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

type pdqsortPanic int

func TestGoLeanPdqsort(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SLICES_PDQSORT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Words                           []int64
			A, B, Limit, Mode, FailAt, Mask int
			Want                            string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing partition fixtures")
			}
			t.Logf("two pinned pdqsort entries matched %d paired rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		ordered := slices.Clone(c.Words)
		first := capturePdqsort(func() string {
			slices.GoLeanPdqsortOrdered(ordered, c.A, c.B, c.Limit)
			return "ok"
		})
		data := slices.Clone(c.Words)
		var calls []string
		var compare func(int64, int64) int
		if c.Mode != 0 {
			compare = func(x, y int64) int {
				call := len(calls)
				calls = append(calls, fmt.Sprintf("%d,%d", x, y))
				if call == c.FailAt {
					panic(pdqsortPanic(call))
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
		second := capturePdqsort(func() string {
			slices.GoLeanPdqsortCmpFunc(data, c.A, c.B, c.Limit, compare)
			return "ok"
		})
		verifyPermutation(t, i, c.Words, ordered, c.A, c.B, true)
		verifyPermutation(t, i, c.Words, data, c.A, c.B, c.Mode == 1 && second == "ok")
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

func capturePdqsort(call func() string) (out string) {
	defer func() {
		switch failure := recover().(type) {
		case nil:
		case pdqsortPanic:
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

func verifyPermutation(t *testing.T, row int, old, result []int64, a, b int, sorted bool) {
	t.Helper()
	if !slices.Equal(old[b:], result[b:]) {
		t.Fatalf("row %d wrote after its target range", row)
	}
	if sorted && (!slices.Equal(old[:a], result[:a]) || !slices.IsSorted(result[a:b])) {
		t.Fatalf("row %d failed sorted subrange or prefix preservation", row)
	}
	frequencies := make(map[int64]int)
	for _, v := range old {
		frequencies[v]++
	}
	for _, v := range result {
		frequencies[v]--
	}
	for v, count := range frequencies {
		if count != 0 {
			t.Fatalf("row %d changed total frequency of %d", row, v)
		}
	}
}
