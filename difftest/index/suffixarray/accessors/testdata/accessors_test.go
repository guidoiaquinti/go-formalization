package suffixarray

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanSuffixarrayAccessors(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SUFFIXARRAY_ACCESSORS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	r := json.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Tag, Len32, Len64, ByteLen, I, J int
			Value                            int64
			Words32, Words64                 []int64
			Bytes                            []byte
			Want                             string
		}
		if err := r.Decode(&c); err == io.EOF {
			if row == 0 {
				t.Fatal("empty suffixarray accessors fixtures")
			}
			t.Logf("six pinned suffixarray accessors matched %d rows", row)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		var a *ints
		var x *Index
		if c.Tag&8 == 0 {
			a = new(ints)
			if c.Words32 != nil {
				data32 := make([]int32, len(c.Words32))
				for i, v := range c.Words32 {
					data32[i] = int32(v)
				}
				a.int32 = data32[:c.Len32]
			}
			if c.Words64 != nil {
				a.int64 = c.Words64[:c.Len64:len(c.Words64)]
			}
			x = &Index{data: c.Bytes[:c.ByteLen:len(c.Bytes)], sa: *a}
		}
		length := observe(func() string { return strconv.Itoa(a.len()) })
		get := observe(func() string { return strconv.FormatInt(a.get(c.I), 10) })
		sub := observe(func() string {
			view := a.slice(c.I, c.J)
			if a.int32 != nil {
				assertAlias(t, view.int32, a.int32, c.I)
			} else {
				assertAlias(t, view.int64, a.int64, c.I)
			}
			return wire(view.int32) + "/" + wire(view.int64)
		})
		data := observe(func() string {
			out := x.Bytes()
			assertAlias(t, out, x.data, 0)
			return wire(out)
		})
		at := observe(func() string {
			out := x.at(c.I)
			assertAlias(t, out, x.data, int(a.get(c.I)))
			return wire(out)
		})
		before := ""
		if a != nil {
			before = wire(a.int32) + "/" + wire(a.int64)
		}
		store := observe(func() string {
			a.set(c.I, c.Value)
			return wire(a.int32) + "/" + wire(a.int64)
		})
		if store == "bounds" && before != wire(a.int32)+"/"+wire(a.int64) {
			t.Fatal("bounds panic changed backing state")
		}
		got := strings.Join([]string{length, get, store, sub, data, at}, "|")
		if got != c.Want {
			t.Fatalf("row %d %+v: native=%q want=%q", row, c, got, c.Want)
		}
	}
}

func assertAlias[T any](t *testing.T, view, original []T, offset int) {
	t.Helper()
	if cap(view) > 0 && &view[:cap(view)][0] != &original[:cap(original)][offset] {
		t.Fatal("returned view does not alias backing at the expected offset")
	}
}

func observe(call func() string) (report string) {
	defer func() {
		if failure := recover(); failure != nil {
			if err, ok := failure.(runtime.Error); ok {
				switch msg := err.Error(); {
				case msg == "runtime error: invalid memory address or nil pointer dereference":
					report = "nilReceiver"
					return
				case strings.HasPrefix(msg, "runtime error: index out of range"), strings.HasPrefix(msg, "runtime error: slice bounds out of range"):
					report = "bounds"
					return
				}
			}
			panic(failure)
		}
	}()
	return "ok:" + call()
}

func wire[T ~byte | ~int32 | ~int64](view []T) string {
	if view == nil {
		return "nil"
	}
	backing := view[:cap(view)]
	words := make([]string, len(backing))
	for i, v := range backing {
		words[i] = strconv.FormatInt(int64(v), 10)
	}
	return fmt.Sprintf("%d:%s", len(view), strings.Join(words, ","))
}
