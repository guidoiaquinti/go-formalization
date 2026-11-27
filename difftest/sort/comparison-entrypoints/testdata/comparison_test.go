package sort

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type goLeanComparison struct{ less func(int, int) bool }

func (goLeanComparison) Len() int             { panic("unexpected Len") }
func (goLeanComparison) Swap(int, int)        { panic("unexpected Swap") }
func (r goLeanComparison) Less(i, j int) bool { return r.less(i, j) }

func goLeanLenValue(shape, n int) any {
	switch shape {
	case 0:
		return make([]int, n)
	case 1:
		return reflect.Zero(reflect.ArrayOf(n, reflect.TypeFor[byte]())).Interface()
	case 2:
		return strings.Repeat("\xff", n)
	case 3:
		m := make(map[int]struct{}, n)
		for i := range n {
			m[i] = struct{}{}
		}
		return m
	case 4:
		ch := make(chan struct{}, n+1)
		for range n {
			ch <- struct{}{}
		}
		return ch
	case 5:
		return nil
	case 6:
		return 112
	default:
		return new([]int)
	}
}

func TestGoLeanComparisonEntrypoints(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SORT_COMPARISON_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var counts [2]int
	for scanner.Scan() {
		var c struct {
			Kind, Shape, A, B int
			Replies           []int
			Want              string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		marker := new(int)
		var calls []string
		result := func() (out string) {
			defer func() {
				if p := recover(); p != nil {
					if p == marker {
						out = "callback"
						return
					}
					text := fmt.Sprint(p)
					switch {
					case c.Kind == 0 && c.Shape == 5 && strings.Contains(text, "reflectlite.Value.Len on zero Value"):
						out = "invalid"
					case c.Kind == 0 && c.Shape >= 6 && strings.Contains(text, "reflectlite.Value.Len on "):
						out = "unsupported"
					case c.Kind == 1 && c.Shape != 0:
						if _, ok := p.(runtime.Error); !ok || !strings.Contains(text, "nil pointer dereference") {
							panic(p)
						}
						out = "nil"
					default:
						panic(p)
					}
				}
			}()
			less := func(i, j int) bool {
				calls = append(calls, fmt.Sprintf("%d,%d", i, j))
				if c.Kind == 0 {
					r := c.Replies[j]
					if r == 2 {
						panic(marker)
					}
					return r == 1
				}
				if c.B == 2 {
					panic(marker)
				}
				return i < j
			}
			if c.Kind == 0 {
				return strconv.FormatBool(SliceIsSorted(goLeanLenValue(c.Shape, c.A), less))
			}
			var underlying Interface
			if c.Shape == 0 {
				underlying = goLeanComparison{less}
			}
			r := reverse{underlying}
			return strconv.FormatBool(r.Less(c.A, c.Replies[0]))
		}()
		if got := result + ":" + strings.Join(calls, ";"); got != c.Want {
			t.Fatalf("%+v: native %q", c, got)
		}
		counts[c.Kind]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("native exact result and callbacks per helper: %v", counts)
}
