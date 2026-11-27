package slices_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanReverseRange(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SLICE_REVERSE_RANGE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			Backing            []int
			A, B, Length, Mode int
			Want               string
		}
		if err := decoder.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("missing range cases")
			}
			t.Logf("both pinned private reversals matched %d rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		ordered, cmp := slices.Clone(row.Backing), slices.Clone(row.Backing)
		var compare func(int, int) int
		if row.Mode != 0 {
			compare = func(int, int) int { panic("reverseRange called the comparator") }
		}
		slices.GoLeanReverseRangeOrdered(ordered[:row.Length:len(ordered)], row.A, row.B)
		slices.GoLeanReverseRangeCmpFunc(cmp[:row.Length:len(cmp)], row.A, row.B, compare)
		got := fmt.Sprintf("%s|%s", backingWire(ordered), backingWire(cmp))
		if got != row.Want {
			t.Fatalf("row %d range [%d,%d) got %q want %q", i, row.A, row.B, got, row.Want)
		}
	}
}

func backingWire(cells []int) string {
	if cells == nil {
		return "nil"
	}
	result := make([]string, len(cells))
	for i, value := range cells {
		result[i] = strconv.Itoa(value)
	}
	return "slice:" + strings.Join(result, ",")
}
