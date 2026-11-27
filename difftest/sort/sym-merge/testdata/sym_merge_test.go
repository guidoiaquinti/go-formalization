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

type goLeanMergeData struct {
	keys, ids []int
	calls     []string
}

func (*goLeanMergeData) Len() int { panic("unexpected Len") }
func (d *goLeanMergeData) Less(i, j int) bool {
	d.calls = append(d.calls, fmt.Sprintf("L%d,%d", i, j))
	return d.keys[d.ids[i]] < d.keys[d.ids[j]]
}
func (d *goLeanMergeData) Swap(i, j int) {
	d.calls = append(d.calls, fmt.Sprintf("S%d,%d", i, j))
	d.ids[i], d.ids[j] = d.ids[j], d.ids[i]
}

func TestGoLeanSymMerge(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SORT_SYM_MERGE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var counts [2]int
	for scanner.Scan() {
		var c struct {
			Variant, A, M, B int
			Keys             []int
			Want             string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		d := &goLeanMergeData{keys: c.Keys, ids: make([]int, len(c.Keys))}
		for i := range d.ids {
			d.ids[i] = i
		}
		if c.Variant == 0 {
			symMerge(d, c.A, c.M, c.B)
		} else {
			symMerge_func(lessSwap{Less: d.Less, Swap: d.Swap}, c.A, c.M, c.B)
		}
		want := make([]int, 0, c.B-c.A)
		l, r := c.A, c.M
		for l < c.M && r < c.B {
			if c.Keys[l] <= c.Keys[r] {
				want = append(want, l)
				l++
			} else {
				want = append(want, r)
				r++
			}
		}
		for ; l < c.M; l++ {
			want = append(want, l)
		}
		for ; r < c.B; r++ {
			want = append(want, r)
		}
		if !slices.Equal(want, d.ids[c.A:c.B]) {
			t.Fatalf("%+v: unstable merge %v want %v", c, d.ids, want)
		}
		fields := make([]string, len(d.ids))
		for i, id := range d.ids {
			if (i < c.A || i >= c.B) && i != id {
				t.Fatal("merged beyond range")
			}
			fields[i] = strconv.Itoa(id)
		}
		got := strings.Join(fields, ",") + ":" + strings.Join(d.calls, ";")
		if got != c.Want {
			t.Fatalf("%+v: native %q", c, got)
		}
		counts[c.Variant]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("symMerge exact native observation counts: %v", counts)
}
