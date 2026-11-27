package suffixarray

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanSuffixarrayLookup(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SUFFIXARRAY_LOOKUP_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Profile, Limit int
			Text, Needle   []byte
			Offsets        []int
			Want           string
		}
		if err := dec.Decode(&c); err == io.EOF {
			if row == 0 {
				t.Fatal("empty lookup fixtures")
			}
			t.Logf("two pinned suffixarray lookup entries matched %d rows", row)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		var index *Index
		if c.Profile != 0 {
			index = &Index{data: c.Text}
			switch c.Profile {
			case 1:
				values := make([]int32, len(c.Offsets))
				for i, v := range c.Offsets {
					values[i] = int32(v)
				}
				index.sa.int32 = values[:len(c.Text)]
			case 2:
				values := make([]int64, len(c.Offsets))
				for i, v := range c.Offsets {
					values[i] = int64(v)
				}
				index.sa.int64 = values[:len(c.Text)]
			}
		}
		all := observeLookup(func() string {
			part := index.lookupAll(c.Needle)
			if index.sa.int32 != nil {
				checkAlias(t, index.sa.int32, part.int32)
			} else {
				checkAlias(t, index.sa.int64, part.int64)
			}
			return renderView(part.int32) + "/" + renderView(part.int64)
		})
		public := observeLookup(func() string {
			result := index.Lookup(c.Needle, c.Limit)
			before := ""
			if index != nil {
				before = renderView(index.sa.int32) + "/" + renderView(index.sa.int64)
			}
			if cap(result) != len(result) {
				t.Fatal("Lookup result has unexpected spare capacity")
			}
			description := renderView(result)
			for i := range result {
				result[i] = -1
			}
			if index != nil && before != renderView(index.sa.int32)+"/"+renderView(index.sa.int64) {
				t.Fatal("Lookup result aliases index")
			}
			return description
		})
		if got := all + "|" + public; got != c.Want {
			t.Fatalf("row %d %+v: native=%q want=%q", row, c, got, c.Want)
		}
	}
}

func checkAlias[T comparable](t *testing.T, original, view []T) {
	t.Helper()
	if cap(view) > 0 {
		offset := cap(original) - cap(view)
		backing := original[:cap(original)]
		if offset < 0 || &backing[offset] != &view[:cap(view)][0] ||
			!slices.Equal(backing[offset:], view[:cap(view)]) {
			t.Fatal("lookupAll failed to retain the expected alias")
		}
	}
}

func observeLookup(call func() string) (result string) {
	defer func() {
		if failure := recover(); failure != nil {
			if err, ok := failure.(runtime.Error); !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			result = "panic"
		}
	}()
	return call()
}

func renderView[T ~int | ~int32 | ~int64](view []T) string {
	if view == nil {
		return "nil"
	}
	backing := view[:cap(view)]
	fields := make([]string, len(backing))
	for i, n := range backing {
		fields[i] = strconv.FormatInt(int64(n), 10)
	}
	return fmt.Sprintf("%d:%s", len(view), strings.Join(fields, ","))
}
