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

func TestGoLeanSuffixarrayMapLms(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_SUFFIXARRAY_MAP_LMS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Length, Count int
			Backing       []int64
			Want          string
		}
		if err := dec.Decode(&c); err == io.EOF {
			if row == 0 {
				t.Fatal("missing map fixtures")
			}
			t.Logf("both pinned LMS map declarations matched %d rows", row)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		wide := slices.Clone(c.Backing)
		narrow := make([]int32, len(c.Backing))
		for i, v := range c.Backing {
			narrow[i] = int32(v)
		}
		status32 := invokeMap(func() { map_32(narrow[:c.Length], c.Count) })
		status64 := invokeMap(func() { map_64(wide[:c.Length], c.Count) })
		if got := status32 + ":" + mapWire(narrow) + "|" + status64 + ":" + mapWire(wide); got != c.Want {
			t.Fatalf("row %d %+v got %q", row, c, got)
		}
	}
}

func invokeMap(call func()) (status string) {
	defer func() {
		if failure := recover(); failure != nil {
			if err, ok := failure.(runtime.Error); !ok || !strings.HasPrefix(err.Error(), "runtime error: index out of range") {
				panic(failure)
			}
			status = "bounds"
		}
	}()
	call()
	return "ok"
}

func mapWire[T ~int32 | ~int64](cells []T) string {
	fields := make([]string, len(cells))
	for i, n := range cells {
		fields[i] = strconv.FormatInt(int64(n), 10)
	}
	return fmt.Sprint(strings.Join(fields, ","))
}
