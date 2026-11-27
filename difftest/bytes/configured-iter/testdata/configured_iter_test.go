package bytes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanConfiguredIter(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_CONFIGURED_ITER_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var rows []struct {
		Input struct {
			NonNil              bool
			Raw, Spare, Sep     []byte
			Save, First, Second int
		}
		Want string
	}
	if err := json.NewDecoder(file).Decode(&rows); err != nil {
		t.Fatal(err)
	}
	for i, row := range rows {
		c := row.Input
		storage := func() []byte {
			if !c.NonNil {
				return nil
			}
			full := make([]byte, len(c.Raw)+len(c.Spare))
			copy(full, c.Raw)
			copy(full[len(c.Raw):], c.Spare)
			return full[:len(c.Raw):len(full)]
		}
		lines := nativeRuns(bytes.Lines(storage()), c.First, c.Second)
		splits := nativeRuns(bytes.GoLeanConfiguredSplit(storage(), c.Sep, c.Save), c.First, c.Second)
		if got := lines + "#" + splits; got != row.Want {
			t.Fatalf("row %d (%+v): got %q want %q", i, c, got, row.Want)
		}
	}
}

func nativeRuns(action iter.Seq[[]byte], first, second int) string {
	once := func(budget int) (row string) {
		status := "ok"
		var parts []string
		defer func() {
			if p := recover(); p != nil {
				err, ok := p.(runtime.Error)
				if !ok || !strings.HasPrefix(err.Error(), "runtime error: slice bounds out of range") {
					panic(p)
				}
				status = "panic"
			}
			row = status + "/" + strings.Join(parts, ";")
		}()
		if budget == 0 {
			return
		}
		action(func(part []byte) bool {
			tag := "slice"
			if part == nil {
				tag = "nil"
			}
			fields := make([]string, len(part))
			for i, b := range part {
				fields[i] = strconv.Itoa(int(b))
			}
			parts = append(parts, fmt.Sprintf("%s:%d:%s", tag, cap(part), strings.Join(fields, ",")))
			extended := append(part, 0xff)
			if len(part) > 0 && &extended[0] == &part[0] {
				panic("append alias survived clipping")
			}
			return len(parts) < budget
		})
		return
	}
	return strings.Join([]string{once(first), once(second), once(4096), once(4096)}, "|")
}
