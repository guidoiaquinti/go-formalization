package time_test

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGoLeanFormatHelpers(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TIME_FMT_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := gob.NewDecoder(bufio.NewReader(f))
	for row := 0; ; row++ {
		var c struct {
			Value             uint64
			Length, Precision int
			Backing           []byte
			Nil               bool
			Want              string
		}
		if err := input.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row > 0 {
				t.Logf("%d observations of both pinned helpers", row)
				return
			}
			t.Fatal(err)
		}
		run := func(format func([]byte) (int, uint64)) string {
			var full []byte
			if !c.Nil {
				full = make([]byte, len(c.Backing))
				copy(full, c.Backing)
			}
			result := func() (out string) {
				defer func() {
					if p := recover(); p != nil {
						if e, ok := p.(runtime.Error); ok && strings.Contains(e.Error(), "index out of range") {
							out = "bounds"
						} else {
							panic(p)
						}
					}
				}()
				start, rest := format(full[:c.Length:len(full)])
				return fmt.Sprintf("%d:%d", start, rest)
			}()
			digits := make([]string, len(full))
			for i, b := range full {
				digits[i] = strconv.Itoa(int(b))
			}
			return result + ":" + strings.Join(digits, ",")
		}
		first := run(func(buf []byte) (int, uint64) { return time.GoLeanFmtInt(buf, c.Value), 0 })
		second := run(func(buf []byte) (int, uint64) { return time.GoLeanFmtFrac(buf, c.Value, c.Precision) })
		if got := first + "|" + second; got != c.Want {
			t.Fatalf("row %d: pinned Go %q expected %q", row, got, c.Want)
		}
	}
}
