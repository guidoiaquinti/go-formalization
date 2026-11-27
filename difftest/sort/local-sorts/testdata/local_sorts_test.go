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

type goLeanLocalSortData struct {
	values []int
	calls  []string
}

func (*goLeanLocalSortData) Len() int { panic("unexpected Len") }
func (d *goLeanLocalSortData) Less(i, j int) bool {
	d.calls = append(d.calls, fmt.Sprintf("L%d,%d", i, j))
	return d.values[i] < d.values[j]
}
func (d *goLeanLocalSortData) Swap(i, j int) {
	d.calls = append(d.calls, fmt.Sprintf("S%d,%d", i, j))
	d.values[i], d.values[j] = d.values[j], d.values[i]
}
func (d *goLeanLocalSortData) wire() string {
	shown := make([]string, len(d.values))
	for i, v := range d.values {
		shown[i] = strconv.Itoa(v)
	}
	return strings.Join(shown, ",") + ":" + strings.Join(d.calls, ";")
}

// Carry the original root value down the larger-child path, filling each
// vacated position from its selected child, then write the saved value once.
func goLeanSiftOracle(values []int, root, end, first int) {
	if root >= end {
		return
	}
	carry := values[first+root]
	for {
		child := root*2 + 1
		if child >= end {
			break
		}
		if child+1 < end && values[first+child+1] > values[first+child] {
			child++
		}
		if values[first+child] <= carry {
			break
		}
		values[first+root] = values[first+child]
		root = child
	}
	values[first+root] = carry
}

func TestGoLeanLocalSorts(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SORT_LOCAL_SORT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	counts := [3][2]int{}
	for scanner.Scan() {
		var c struct {
			Kind, Variant, A, B, C int
			Values                 []int
			Want                   string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		a, b, first := c.A, c.B, c.C
		data := &goLeanLocalSortData{values: slices.Clone(c.Values)}
		want := slices.Clone(c.Values)
		callbacks := lessSwap{Less: data.Less, Swap: data.Swap}
		switch c.Kind {
		case 0:
			if c.Variant == 0 {
				insertionSort(data, a, b)
			} else {
				insertionSort_func(callbacks, a, b)
			}
			slices.Sort(want[a:b])
		case 1:
			if c.Variant == 0 {
				siftDown(data, a, b, first)
			} else {
				siftDown_func(callbacks, a, b, first)
			}
			goLeanSiftOracle(want, a, b, first)
		case 2:
			if c.Variant == 0 {
				heapSort(data, a, b)
			} else {
				heapSort_func(callbacks, a, b)
			}
			slices.Sort(want[a:b])
		}
		if !slices.Equal(want, data.values) || c.Want != data.wire() {
			t.Fatalf("%+v: got %s, want values %v", c, data.wire(), want)
		}
		counts[c.Kind][c.Variant]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("native rows by helper and variant: %v", counts)
}
