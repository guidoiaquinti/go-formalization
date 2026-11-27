package time_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func reportDurationFormat(d time.Duration, backing *[32]byte) (result string) {
	defer func() {
		if e := recover(); e != nil {
			failure, ok := e.(runtime.Error)
			if !ok || !strings.Contains(failure.Error(), "nil pointer") {
				panic(e)
			}
			result = "nil"
		}
	}()
	start := time.GoLeanDurationFormat(d, backing)
	words := make([]string, len(backing))
	for i, cell := range backing {
		words[i] = strconv.Itoa(int(cell))
	}
	return fmt.Sprintf("%d:%s", start, strings.Join(words, ","))
}

func TestGoLeanDurationText(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_DURATION_TEXT_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := json.NewDecoder(f)
	for row := 0; ; row++ {
		var c struct {
			Value   int64
			Present bool
			Cells   [32]byte
			Want    string
		}
		err := in.Decode(&c)
		if err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		var p *[32]byte
		if c.Present {
			p = &c.Cells
		}
		if got := reportDurationFormat(time.Duration(c.Value), p); got != c.Want {
			t.Fatalf("row %d: Go %q want %q", row, got, c.Want)
		}
	}
}
