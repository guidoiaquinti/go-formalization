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

type goLeanPdqData struct {
	values []int
	calls  []string
}

func (*goLeanPdqData) Len() int { panic("unexpected Len") }
func (d *goLeanPdqData) Less(i, j int) bool {
	d.calls = append(d.calls, fmt.Sprintf("L%d,%d", i, j))
	return d.values[i] < d.values[j]
}
func (d *goLeanPdqData) Swap(i, j int) {
	d.calls = append(d.calls, fmt.Sprintf("S%d,%d", i, j))
	d.values[i], d.values[j] = d.values[j], d.values[i]
}

func TestGoLeanPdqsort(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_PDQ_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var counts [2]int
	for scanner.Scan() {
		var c struct {
			Variant, A, B, Limit int
			Values               []int
			Want                 string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		d := &goLeanPdqData{values: slices.Clone(c.Values)}
		if c.Variant == 0 {
			pdqsort(d, c.A, c.B, c.Limit)
		} else {
			pdqsort_func(lessSwap{Less: d.Less, Swap: d.Swap}, c.A, c.B, c.Limit)
		}
		want := slices.Clone(c.Values)
		slices.Sort(want[c.A:c.B])
		if !slices.Equal(want, d.values) {
			t.Fatalf("range [%d:%d] from %v: native %v, sorted %v", c.A, c.B, c.Values, d.values, want)
		}
		fields := make([]string, len(c.Values))
		for i, v := range d.values {
			fields[i] = strconv.Itoa(v)
		}
		got := strings.Join(fields, ",") + ":" + strings.Join(d.calls, ";")
		if got != c.Want {
			t.Fatalf("kind %d sort [%d:%d] limit %d input %v: native %s, Lean %s", c.Variant, c.A, c.B, c.Limit, c.Values, got, c.Want)
		}
		counts[c.Variant]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("pdqsort exact trace counts: %v", counts)
}
