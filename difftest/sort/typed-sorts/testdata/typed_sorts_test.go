package sort

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func goLeanSameTypedCounts[T comparable](input, output []T) bool {
	if len(input) != len(output) {
		return false
	}
	counts := make(map[T]int)
	for _, v := range input {
		counts[v]++
	}
	for _, v := range output {
		if counts[v] == 0 {
			return false
		}
		counts[v]--
	}
	return true
}

func TestGoLeanTypedSorts(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TYPED_SORT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var counts [6]int
	for scanner.Scan() {
		var c struct {
			Kind   int
			Ints   []int
			Floats []uint64
			Words  [][]byte
			Want   string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		var got string
		switch c.Kind {
		case 0:
			arr := slices.Clone(c.Ints)
			Ints(arr)
			if !slices.IsSorted(arr) || !goLeanSameTypedCounts(c.Ints, arr) {
				t.Fatalf("invalid int output %v", arr)
			}
			fields := make([]string, len(arr))
			for i, v := range arr {
				fields[i] = strconv.Itoa(v)
			}
			got = "values:" + strings.Join(fields, ",")
		case 1, 2:
			var arr []float64
			if c.Floats != nil {
				arr = make([]float64, len(c.Floats))
			}
			for i, v := range c.Floats {
				arr[i] = math.Float64frombits(v)
			}
			if c.Kind == 1 {
				Float64s(arr)
			} else {
				Float64Slice(arr).Sort()
			}
			bits := make([]uint64, len(arr))
			fields := make([]string, len(arr))
			for i, v := range arr {
				if i > 0 && ((!math.IsNaN(arr[i-1]) && math.IsNaN(v)) || v < arr[i-1]) {
					t.Fatalf("nonmonotone float output %v", arr)
				}
				bits[i] = math.Float64bits(v)
				fields[i] = strconv.FormatUint(bits[i], 10)
			}
			if !goLeanSameTypedCounts(c.Floats, bits) {
				t.Fatalf("changed float bits %v want %v", bits, c.Floats)
			}
			got = "values:" + strings.Join(fields, ",")
		case 3, 4:
			var arr []string
			if c.Words != nil {
				arr = make([]string, len(c.Words))
			}
			for i, v := range c.Words {
				arr[i] = string(v)
			}
			before := slices.Clone(arr)
			if c.Kind == 3 {
				Strings(arr)
			} else {
				StringSlice(arr).Sort()
			}
			if !slices.IsSorted(arr) || !goLeanSameTypedCounts(before, arr) {
				t.Fatalf("invalid string output %q", arr)
			}
			fields := make([]string, len(arr))
			for i, v := range arr {
				raw := make([]string, len(v))
				for j := range len(v) {
					raw[j] = strconv.Itoa(int(v[j]))
				}
				fields[i] = fmt.Sprintf("%d:%s", len(v), strings.Join(raw, ","))
			}
			got = "values:" + strings.Join(fields, ";")
		case 5:
			f := math.Float64frombits(c.Floats[0])
			actual := isNaN(f)
			mag := c.Floats[0] & ((1 << 63) - 1)
			if actual != (mag > 0x7ff0000000000000) {
				t.Fatalf("invalid isNaN for %x", c.Floats[0])
			}
			got = strconv.FormatBool(actual)
		}
		if got != c.Want {
			t.Fatalf("%+v: native %q", c, got)
		}
		counts[c.Kind]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("typed sort native observations: %v", counts)
}
