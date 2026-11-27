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

type stablePanic int

func TestGoLeanStable(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SLICES_STABLE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Words                             []int64
			Kind, A, M, B, Mode, FailAt, Mask int
			Want                              string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing partition fixtures")
			}
			t.Logf("four pinned stable-sort entries matched %d paired rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		ordered := slices.Clone(c.Words)
		first := captureStable(func() string {
			if c.Kind == 0 {
				slices.GoLeanSymMergeOrdered(ordered, c.A, c.M, c.B)
			} else {
				slices.GoLeanStableOrdered(ordered, c.B)
			}
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
					panic(stablePanic(call))
				}
				if c.Mode == 1 {
					return cmp.Compare(x>>8, y>>8)
				}
				if c.Mask&(1<<(call%12)) != 0 {
					return -1
				}
				return 0
			}
		}
		second := captureStable(func() string {
			if c.Kind == 0 {
				slices.GoLeanSymMergeCmpFunc(data, c.A, c.M, c.B, compare)
			} else {
				slices.GoLeanStableCmpFunc(data, c.B, compare)
			}
			return "ok"
		})
		verifySort(t, i, c.Words, ordered, c.A, c.B, true, false)
		verifySort(t, i, c.Words, data, c.A, c.B, c.Mode == 1 && second == "ok", true)
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

func captureStable(call func() string) (out string) {
	defer func() {
		switch failure := recover().(type) {
		case nil:
		case stablePanic:
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

func verifySort(t *testing.T, row int, old, result []int64, a, b int, checkSorted, keys bool) {
	t.Helper()
	if !slices.Equal(old[:a], result[:a]) || !slices.Equal(old[b:], result[b:]) {
		t.Fatalf("row %d wrote outside its target range", row)
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
	if !checkSorted {
		return
	}
	ordered := slices.Clone(old[a:b])
	key := func(x int64) int64 {
		if keys {
			return x >> 8
		}
		return x
	}
	// A minimum selection that removes the earliest equal value is an
	// independent stable reference and keeps the occurrence tags.
	for i := range ordered {
		smallest := i
		for j := i + 1; j < len(ordered); j++ {
			if key(ordered[j]) < key(ordered[smallest]) {
				smallest = j
			}
		}
		value := ordered[smallest]
		copy(ordered[i+1:smallest+1], ordered[i:smallest])
		ordered[i] = value
	}
	if !slices.Equal(result[a:b], ordered) {
		t.Fatalf("row %d sorted/stable mismatch: %v != %v", row, result[a:b], ordered)
	}
}
