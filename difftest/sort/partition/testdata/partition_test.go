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

type goLeanPartitionData struct {
	values []int
	calls  []string
}

func (*goLeanPartitionData) Len() int { panic("unexpected Len") }
func (d *goLeanPartitionData) Less(i, j int) bool {
	d.calls = append(d.calls, fmt.Sprintf("L%d,%d", i, j))
	return d.values[i] < d.values[j]
}
func (d *goLeanPartitionData) Swap(i, j int) {
	d.calls = append(d.calls, fmt.Sprintf("S%d,%d", i, j))
	d.values[i], d.values[j] = d.values[j], d.values[i]
}
func (d *goLeanPartitionData) wire(mid int, status string) string {
	values := make([]string, len(d.values))
	for i, v := range d.values {
		values[i] = strconv.Itoa(v)
	}
	return fmt.Sprintf("%d:%s:%s:%s", mid, status, strings.Join(values, ","), strings.Join(d.calls, ";"))
}

func TestGoLeanPartition(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SORT_PARTITION_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var counts [2][2]int
	for scanner.Scan() {
		var c struct {
			Kind, Variant, A, B, Pivot int
			Values                     []int
			Want                       string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		a, b, p := c.A, c.B, c.Pivot
		data := &goLeanPartitionData{values: slices.Clone(c.Values)}
		callbacks := lessSwap{Less: data.Less, Swap: data.Swap}
		pivot := c.Values[p]
		status := "-"
		var mid int
		if c.Kind == 0 {
			var already bool
			if c.Variant == 0 {
				mid, already = partition(data, a, b, p)
			} else {
				mid, already = partition_func(callbacks, a, b, p)
			}
			status = strconv.FormatBool(already)
			wantMid := a
			for _, v := range c.Values[a:b] {
				if v < pivot {
					wantMid++
				}
			}
			// The first pivot swap is outside the already-partitioned test.
			start := slices.Clone(c.Values)
			start[a], start[p] = start[p], start[a]
			wantAlready := !slices.ContainsFunc(start[a+1:wantMid+1], func(v int) bool { return v >= pivot })
			if mid != wantMid || data.values[mid] != pivot || already != wantAlready ||
				slices.ContainsFunc(data.values[a:mid], func(v int) bool { return v >= pivot }) ||
				slices.ContainsFunc(data.values[mid+1:b], func(v int) bool { return v < pivot }) {
				t.Fatalf("%+v: incorrect partition %v at %d, already=%v", c, data.values, mid, already)
			}
		} else {
			if c.Variant == 0 {
				mid = partitionEqual(data, a, b, p)
			} else {
				mid = partitionEqual_func(callbacks, a, b, p)
			}
			wantMid := a
			for _, v := range c.Values[a:b] {
				if v == pivot {
					wantMid++
				}
			}
			if mid != wantMid ||
				slices.ContainsFunc(data.values[a:mid], func(v int) bool { return v != pivot }) ||
				slices.ContainsFunc(data.values[mid:b], func(v int) bool { return v <= pivot }) {
				t.Fatalf("%+v: incorrect equal partition %v at %d", c, data.values, mid)
			}
		}
		if !slices.Equal(c.Values[:a], data.values[:a]) || !slices.Equal(c.Values[b:], data.values[b:]) ||
			!slices.Equal(sorted(c.Values[a:b]), sorted(data.values[a:b])) ||
			data.wire(mid, status) != c.Want {
			t.Fatalf("%+v: native %s", c, data.wire(mid, status))
		}
		counts[c.Kind][c.Variant]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("partition native rows: %v", counts)
}
func sorted(values []int) []int {
	result := slices.Clone(values)
	slices.Sort(result)
	return result
}
