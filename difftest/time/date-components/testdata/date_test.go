package time_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestGoLeanDateComponents(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TIME_DATE_COMPONENTS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for i := 0; ; i++ {
		var row struct {
			Year, Month, Day, Hi, Lo, Base int
			Want                           string
		}
		if err := dec.Decode(&row); err != nil {
			if err == io.EOF && i != 0 {
				t.Logf("%d pinned rows per declaration", i)
				return
			}
			t.Fatal(err)
		}
		date, leap, fields := time.GoLeanDateComponents(row.Year, row.Month, row.Day, row.Hi, row.Lo, row.Base)
		got := fmt.Sprintf("%d|%t|%d|%d|%d,%d", date, leap, fields[0], fields[1], fields[2], fields[3])
		if got != row.Want {
			t.Fatalf("row %d %+v: Go %q, want %q", i, row, got, row.Want)
		}
	}
}
