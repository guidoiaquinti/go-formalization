package sort

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
)

type goLeanRangeData struct {
	data  []int
	trace []string
}

func (d *goLeanRangeData) Len() int           { panic("range helper unexpectedly called Len") }
func (d *goLeanRangeData) Less(int, int) bool { panic("range helper unexpectedly called Less") }
func (d *goLeanRangeData) Swap(i, j int) {
	d.trace = append(d.trace, fmt.Sprintf("%d,%d", i, j))
	d.data[i], d.data[j] = d.data[j], d.data[i]
}

func TestGoLeanRangeMoves(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SORT_RANGE_MOVES_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	counts := [3][2]int{}
	for scanner.Scan() {
		var c struct {
			Kind, Variant, A, B, N, Width int
			Trace                         string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		a, b, n := c.A, c.B, c.N
		values := make([]int, c.Width)
		for i := range values {
			values[i] = i
		}
		data := &goLeanRangeData{data: slices.Clone(values)}
		callbacks := lessSwap{Less: data.Less, Swap: data.Swap}
		switch c.Kind {
		case 0:
			if c.Variant == 0 {
				reverseRange(data, a, b)
			} else {
				reverseRange_func(callbacks, a, b)
			}
			for i := a; i < b; i++ {
				values[i] = a + b - 1 - i
			}
		case 1:
			if c.Variant == 0 {
				swapRange(data, a, b, n)
			} else {
				swapRange_func(callbacks, a, b, n)
			}
			// Map each final index back through the transpositions in reverse order.
			for i := range values {
				source := i
				for offset := n - 1; offset >= 0; offset-- {
					switch source {
					case a + offset:
						source = b + offset
					case b + offset:
						source = a + offset
					}
				}
				values[i] = source
			}
		case 2:
			if c.Variant == 0 {
				rotate(data, a, b, n)
			} else {
				rotate_func(callbacks, a, b, n)
			}
			for i := a; i < n; i++ {
				values[i] = a + (i-a+b-a)%(n-a)
			}
		}
		got := "swaps:" + strings.Join(data.trace, ";")
		if got != c.Trace || !slices.Equal(data.data, values) {
			t.Fatalf("%+v: native trace=%s, final=%v, want final=%v", c, got, data.data, values)
		}
		counts[c.Kind][c.Variant]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("native rows by helper and variant: %v", counts)
}
