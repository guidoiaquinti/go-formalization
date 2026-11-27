package slices_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"slices"
	"testing"
)

type float32Alias float32
type float64Alias float64
type intAlias int64
type stringAlias string

func TestGoLeanSortPrimitives(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SLICES_SORT_PRIMITIVES_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Length  int
			Word    uint64
			Present bool
			Want    string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing fixtures")
			}
			t.Logf("three pinned sort primitives matched %d rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		next := observeXorshift(c.Word, c.Present)
		f32, f64 := math.Float32frombits(uint32(c.Word)), math.Float64frombits(c.Word)
		nonfloat := slices.GoLeanSortIsNaN(intAlias(c.Word)) || slices.GoLeanSortIsNaN(stringAlias("x"))
		got := fmt.Sprintf("%s/%d/%t/%t/%t", next, slices.GoLeanNextPowerOfTwo(c.Length), slices.GoLeanSortIsNaN(float32Alias(f32)), slices.GoLeanSortIsNaN(float64Alias(f64)), nonfloat)
		if got != c.Want {
			t.Fatalf("row %d %+v: got %q", i, c, got)
		}
	}
}

func observeXorshift(word uint64, present bool) (out string) {
	defer func() {
		if failure := recover(); failure != nil {
			if err, ok := failure.(runtime.Error); !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			out = "panic"
		}
	}()
	value, updated := slices.GoLeanXorshiftNext(word, present)
	return fmt.Sprintf("%d,%d", value, updated)
}
