package slices_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type medianPanic int

func TestGoLeanMedian(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SLICES_MEDIAN_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Words                                []int64
			A, B, C, Center, Mode, FailAt, Count int
			Replies                              [3]int
			HasCounter                           bool
			Want                                 string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing median fixtures")
			}
			t.Logf("four pinned median entries matched %d rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		observed := make([]string, 4)
		for kind := range observed {
			data := slices.Clone(c.Words)
			counter := c.Count
			var ptr *int
			if c.HasCounter {
				ptr = &counter
			}
			var calls []string
			var compare func(int64, int64) int
			if c.Mode != 0 {
				compare = func(x, y int64) int {
					call := len(calls)
					calls = append(calls, fmt.Sprintf("%d,%d", x, y))
					if call == c.FailAt {
						panic(medianPanic(call))
					}
					if c.Mode == 1 {
						return cmp.Compare(x, y)
					}
					return c.Replies[call]
				}
			}
			result := captureMedian(func() int {
				switch kind {
				case 0:
					return slices.GoLeanMedianOrdered(data, c.A, c.B, c.C, ptr)
				case 1:
					return slices.GoLeanMedianCmpFunc(data, c.A, c.B, c.C, ptr, compare)
				case 2:
					return slices.GoLeanMedianAdjacentOrdered(data, c.Center, ptr)
				default:
					return slices.GoLeanMedianAdjacentCmpFunc(data, c.Center, ptr, compare)
				}
			})
			count := "nil"
			if ptr != nil {
				count = strconv.Itoa(counter)
			}
			observed[kind] = result + "/" + count + "/" + strings.Join(calls, ";")
			if !slices.Equal(data, c.Words) {
				t.Fatalf("row %d entry %d changed its input", i, kind)
			}
		}
		if got := strings.Join(observed, "|"); got != c.Want {
			t.Fatalf("row %d got %q want %q", i, got, c.Want)
		}
	}
}

func captureMedian(call func() int) (out string) {
	defer func() {
		switch failure := recover().(type) {
		case nil:
		case medianPanic:
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
	return fmt.Sprintf("ok:%d", call())
}
