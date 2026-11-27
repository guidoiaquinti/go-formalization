package flate

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

func observeHashBounds(f func()) (panicked bool) {
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(runtime.Error)
			if !ok || !(strings.Contains(err.Error(), "index out of range") ||
				strings.Contains(err.Error(), "slice bounds out of range")) {
				panic(failure)
			}
			panicked = true
		}
	}()
	f()
	return false
}

func TestGoLeanHashWindow(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_FLATE_HASH_WINDOW_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Size, DstSize int
			Source        []byte
			Dst           []uint32
			Want          string
		}
		err := dec.Decode(&c)
		if err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		src := c.Source[:c.Size:len(c.Source)]
		word := "bounds"
		observeHashBounds(func() { word = fmt.Sprint(hash4(src)) })
		failed := observeHashBounds(func() { bulkHash4(src, c.Dst[:c.DstSize:len(c.Dst)]) })
		fields := make([]string, len(c.Dst))
		for i, x := range c.Dst {
			fields[i] = strconv.FormatUint(uint64(x), 10)
		}
		got := fmt.Sprintf("%s|%t|%s", word, failed, strings.Join(fields, ","))
		if got != c.Want {
			t.Fatalf("row %d: pinned Go %q want %q", row, got, c.Want)
		}
	}
}
