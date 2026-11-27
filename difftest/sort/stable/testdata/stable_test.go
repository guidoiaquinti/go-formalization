package sort

import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type goLeanStableData struct {
	keys, ids   []int
	calls       []string
	n, lenCalls int
}

func (d *goLeanStableData) Len() int { d.lenCalls++; return d.n }
func (d *goLeanStableData) Less(i, j int) bool {
	d.calls = append(d.calls, fmt.Sprintf("L%d,%d", i, j))
	return d.keys[d.ids[i]] < d.keys[d.ids[j]]
}
func (d *goLeanStableData) Swap(i, j int) {
	d.calls = append(d.calls, fmt.Sprintf("S%d,%d", i, j))
	d.ids[i], d.ids[j] = d.ids[j], d.ids[i]
}

type goLeanStableIDs []int

func TestGoLeanStable(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SORT_STABLE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var counts [4]int
	for scanner.Scan() {
		var c struct {
			Kind, Tag, N int
			Keys         []int
			Want         string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		d := &goLeanStableData{keys: c.Keys, ids: make([]int, len(c.Keys)), n: c.N}
		for i := range d.ids {
			d.ids[i] = i
		}
		status := func() (result string) {
			defer func() {
				if p := recover(); p != nil {
					text := fmt.Sprint(p)
					switch {
					case strings.HasPrefix(text, "reflect: call of Swapper on zero Value"):
						result = "invalid"
					case strings.HasPrefix(text, "reflect: call of Swapper on "):
						result = "non-slice"
					case strings.Contains(text, "invalid memory address or nil pointer dereference"):
						if c.Kind == 2 {
							result = "nil-interface"
						} else {
							result = "nil-less"
						}
					default:
						panic(p)
					}
				}
			}()
			switch c.Kind {
			case 0:
				stable(d, c.N)
			case 1:
				stable_func(lessSwap{Less: d.Less, Swap: d.Swap}, c.N)
			case 2:
				if c.Tag == 0 {
					Stable(d)
				} else {
					Stable(nil)
				}
			case 3:
				switch c.Tag {
				case 0:
					SliceStable(goLeanStableIDs(d.ids), d.Less)
				case 1:
					SliceStable(d.ids, nil)
				case 2:
					SliceStable(nil, d.Less)
				case 3:
					SliceStable([2]int{2, 1}, d.Less)
				case 4:
					SliceStable(goLeanStableIDs(nil), d.Less)
				}
			}
			return "ok"
		}()
		wantStatus, n, wantLen := "ok", max(c.N, 0), 0
		if c.Kind == 2 {
			if c.Tag == 0 {
				wantLen = 1
			} else {
				wantStatus = "nil-interface"
				n = 0
			}
		}
		if c.Kind == 3 {
			n = len(c.Keys)
			switch c.Tag {
			case 1:
				if n > 1 {
					wantStatus = "nil-less"
				}
				n = 0
			case 2:
				wantStatus = "invalid"
				n = 0
			case 3:
				wantStatus = "non-slice"
				n = 0
			}
		}
		want := make([]int, len(c.Keys))
		for i := range want {
			want[i] = i
		}
		slices.SortFunc(want[:n], func(i, j int) int {
			if byKey := cmp.Compare(c.Keys[i], c.Keys[j]); byKey != 0 {
				return byKey
			}
			return cmp.Compare(i, j)
		})
		if !slices.Equal(want, d.ids) || status != wantStatus || d.lenCalls != wantLen {
			t.Fatalf("%+v: native %s Len=%d IDs=%v, want %s Len=%d IDs=%v", c, status, d.lenCalls, d.ids, wantStatus, wantLen, want)
		}
		fields := make([]string, len(d.ids))
		for i, id := range d.ids {
			fields[i] = strconv.Itoa(id)
		}
		got := fmt.Sprintf("%s:%d:%s:%s", status, d.lenCalls, strings.Join(fields, ","), strings.Join(d.calls, ";"))
		if got != c.Want {
			t.Fatalf("%+v: native %q", c, got)
		}
		counts[c.Kind]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("stable exact native observations by declaration: %v", counts)
}
