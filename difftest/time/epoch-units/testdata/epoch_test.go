package time_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGoLeanEpochUnits(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TIME_EPOCH_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := json.NewDecoder(f)
	wire := func(value time.Time) string {
		wall, ext, loc := time.GoLeanLocatedParts(value)
		if loc != time.Local {
			t.Fatal("constructor raw location is not time.Local")
		}
		return fmt.Sprintf("%d,%d", wall, ext)
	}
	for row := 0; ; row++ {
		var c struct {
			Wall, Ext        uint64
			Sec, Nsec, Units int64
			Want             string
		}
		err := in.Decode(&c)
		if err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		value := time.GoLeanLocatedTime(c.Wall, c.Ext, time.UTC)
		observed := strings.Join([]string{fmt.Sprint(value.UnixMilli()), fmt.Sprint(value.UnixMicro()), fmt.Sprint(value.UnixNano()),
			wire(time.GoLeanUnixTime(c.Sec, int32(c.Nsec))), wire(time.Unix(c.Sec, c.Nsec)), wire(time.UnixMilli(c.Units)),
			wire(time.UnixMicro(c.Units))}, "|")
		if observed != c.Want {
			t.Fatalf("row %d: Go %q want %q", row, observed, c.Want)
		}
	}
}
