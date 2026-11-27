package sort

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type goLeanPartialData struct {
	values []int
	calls  []string
	swaps  int
}

func (*goLeanPartialData) Len() int { panic("unexpected Len") }
func (d *goLeanPartialData) Less(i, j int) bool {
	d.calls = append(d.calls, fmt.Sprintf("L%d,%d", i, j))
	return d.values[i] < d.values[j]
}
func (d *goLeanPartialData) Swap(i, j int) {
	d.calls = append(d.calls, fmt.Sprintf("S%d,%d", i, j))
	d.values[i], d.values[j] = d.values[j], d.values[i]
	d.swaps++
}

func TestGoLeanPartialInsertion(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SORT_PARTIAL_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var count [2]int
	for scanner.Scan() {
		var c struct {
			Variant, A, B int
			Values        []int
			CrossesStart  bool
			Want          string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		d := &goLeanPartialData{values: slices.Clone(c.Values)}
		var done bool
		if c.Variant == 0 {
			done = partialInsertionSort(d, c.A, c.B)
		} else {
			done = partialInsertionSort_func(lessSwap{Less: d.Less, Swap: d.Swap}, c.A, c.B)
		}
		if c.B-c.A < 50 {
			wantDone := c.A < c.B && slices.IsSorted(c.Values[c.A:c.B])
			if done != wantDone || d.swaps != 0 {
				t.Fatalf("short range mutated or wrong flag: %+v", c)
			}
		}
		if c.A == 0 && done && !slices.IsSorted(d.values[c.A:c.B]) {
			t.Fatalf("successful result not sorted: %+v", c)
		}
		if c.CrossesStart && (!done || !slices.Equal(d.values[:5], []int{-1, 100, 101, 0, 1}) || slices.IsSorted(d.values[c.A:c.B])) {
			t.Fatal("left-shift control did not reproduce the prefix mutation and unsorted success")
		}
		before, after := slices.Clone(c.Values), slices.Clone(d.values)
		slices.Sort(before)
		slices.Sort(after)
		if !slices.Equal(before, after) || !slices.Equal(c.Values[c.B:], d.values[c.B:]) {
			t.Fatalf("partial insertion changed the multiset or trailing exterior: %+v", c)
		}
		fields := make([]string, len(d.values))
		for i, v := range d.values {
			fields[i] = strconv.Itoa(v)
		}
		got := fmt.Sprintf("%t:%s:%s", done, strings.Join(fields, ","), strings.Join(d.calls, ";"))
		if got != c.Want {
			t.Fatalf("%+v: native %s", c, got)
		}
		count[c.Variant]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("partial insertion native observations: %v", count)
}
