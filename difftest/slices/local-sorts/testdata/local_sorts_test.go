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

type sortPanic int

func TestGoLeanLocalSorts(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SLICES_LOCAL_SORT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Words                                []int64
			Kind, A, B, Root, Mode, FailAt, Mask int
			Want                                 string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing local-sort fixtures")
			}
			t.Logf("six pinned local-sort entries matched %d paired rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		ordered := slices.Clone(c.Words)
		first := captureSort(func() {
			switch c.Kind {
			case 0:
				slices.GoLeanInsertionSortOrdered(ordered, c.A, c.B)
			case 1:
				slices.GoLeanSiftDownOrdered(ordered, c.Root, c.B-c.A, c.A)
			case 2:
				slices.GoLeanHeapSortOrdered(ordered, c.A, c.B)
			}
		})
		data := slices.Clone(c.Words)
		var calls []string
		var compare func(int64, int64) int
		if c.Mode != 0 {
			compare = func(x, y int64) int {
				call := len(calls)
				calls = append(calls, fmt.Sprintf("%d,%d", x, y))
				if call == c.FailAt {
					panic(sortPanic(call))
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
		second := captureSort(func() {
			switch c.Kind {
			case 0:
				slices.GoLeanInsertionSortCmpFunc(data, c.A, c.B, compare)
			case 1:
				slices.GoLeanSiftDownCmpFunc(data, c.Root, c.B-c.A, c.A, compare)
			case 2:
				slices.GoLeanHeapSortCmpFunc(data, c.A, c.B, compare)
			}
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

func captureSort(call func()) (out string) {
	defer func() {
		switch failure := recover().(type) {
		case nil:
		case sortPanic:
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
	call()
	return "ok"
}
