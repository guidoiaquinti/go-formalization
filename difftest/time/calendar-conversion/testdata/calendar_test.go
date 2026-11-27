package time_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestGoLeanCalendarConversion(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TIME_CALENDAR"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for i := 0; ; i++ {
		var row struct {
			Days, Century                  uint64
			Year, JanFeb, Month, Day, Leap int
			Want                           string
		}
		if err := dec.Decode(&row); err != nil {
			if err == io.EOF && i > 0 {
				t.Logf("%d pinned rows per declaration", i)
				return
			}
			t.Fatal(err)
		}
		c, v := time.GoLeanCalendarConversion(row.Days, row.Century, row.Year, row.JanFeb, row.Month, row.Day, row.Leap)
		got := fmt.Sprintf("%d,%d,%d|%d|%d|%d|%d,%d,%d|%d,%d", c, v[0], v[1], v[2], v[3], v[4],
			v[5], v[6], v[7], v[8], v[9])
		if got != row.Want {
			t.Fatalf("row %d %+v: Go %q, want %q", i, row, got, row.Want)
		}
	}
}
