package slices_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"slices"
	"strconv"
	"testing"
)

type comparatorPanic int

func TestGoLeanOrder2(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SLICES_ORDER2_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Words                          []uint64
			A, B, Kind, Mode, Reply, Count int
			HasCounter                     bool
			Want                           string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing order2 fixtures")
			}
			t.Logf("both pinned order2 entries matched %d rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}

		orderedCount, cmpCount := c.Count, c.Count
		var orderedPointer, cmpPointer *int
		if c.HasCounter {
			orderedPointer, cmpPointer = &orderedCount, &cmpCount
		}
		floats, signed := make([]float64, len(c.Words)), make([]int64, len(c.Words))
		for j, word := range c.Words {
			floats[j], signed[j] = math.Float64frombits(word), int64(word)
		}
		ordered := capture(func() (int, int) {
			if c.Kind == 1 {
				return slices.GoLeanOrder2Float(floats, c.A, c.B, orderedPointer)
			}
			return slices.GoLeanOrder2Int(signed, c.A, c.B, orderedPointer)
		})
		actualWords := slices.Clone(c.Words)
		calls := ""
		var cmp func(uint64, uint64) int
		if c.Mode != 0 {
			cmp = func(x, y uint64) int {
				if calls != "" {
					t.Fatal("order2CmpFunc invoked comparator more than once")
				}
				calls = fmt.Sprintf("%d,%d", x, y)
				if c.Mode == 2 {
					panic(comparatorPanic(c.Reply))
				}
				return c.Reply
			}
		}
		compare := capture(func() (int, int) {
			return slices.GoLeanOrder2Cmp(actualWords, c.A, c.B, cmpPointer, cmp)
		})
		got := fmt.Sprintf("%s/%s/|%s/%s/%s", ordered, counter(orderedPointer), compare, counter(cmpPointer), calls)
		if got != c.Want {
			t.Fatalf("row %d got %q want %q", i, got, c.Want)
		}
		for j, word := range c.Words {
			if math.Float64bits(floats[j]) != word || uint64(signed[j]) != word || actualWords[j] != word {
				t.Fatalf("row %d changed data at %d", i, j)
			}
		}
	}
}

func capture(call func() (int, int)) (result string) {
	defer func() {
		switch failure := recover().(type) {
		case nil:
		case comparatorPanic:
			result = fmt.Sprintf("panic:%d", failure)
		case runtime.Error:
			if failure.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			result = "nil"
		default:
			panic(failure)
		}
	}()
	a, b := call()
	return fmt.Sprintf("ok:%d,%d", a, b)
}

func counter(value *int) string {
	if value == nil {
		return "nil"
	}
	return strconv.Itoa(*value)
}
