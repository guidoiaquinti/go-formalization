package time_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestGoLeanClockFields(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TIME_FIELDS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for i := 0; ; i++ {
		var row struct {
			Seconds, Days, Century uint64
			Year, Day              uint32
			Want                   string
		}
		if err := dec.Decode(&row); err != nil {
			if err == io.EOF && i > 0 {
				t.Logf("%d pinned rows", i)
				return
			}
			t.Fatal(err)
		}
		fields := time.GoLeanClockFields(row.Seconds, row.Days, row.Century, row.Year, row.Day)
		got := fmt.Sprintf("%d|%d,%d,%d|%d|%d|%d|%d,%d", fields[0], fields[1], fields[2], fields[3],
			fields[4], fields[5], fields[6], fields[7], fields[8])
		if got != row.Want {
			t.Fatalf("row %d %+v: Go %q, want %q", i, row, got, row.Want)
		}
	}
}
