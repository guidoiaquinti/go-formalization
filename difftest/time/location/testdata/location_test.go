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

func TestGoLeanLocation(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TIME_LOCATION_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := json.NewDecoder(f)
	locations := []*time.Location{nil, time.UTC, time.Local, time.FixedZone("alias", 900), time.FixedZone("alias", 900)}
	id := func(loc *time.Location) int {
		for i, p := range locations {
			if loc == p {
				return i
			}
		}
		t.Fatal("location pointer outside fixture")
		return -1
	}
	wire := func(value time.Time) string {
		wall, ext, loc := time.GoLeanLocatedParts(value)
		return fmt.Sprintf("%d,%d,%d", wall, ext, id(loc))
	}
	for row := 0; ; row++ {
		var c struct {
			Wall, Ext uint64
			Old, Next int
			Want      string
		}
		err := in.Decode(&c)
		if err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		value := time.GoLeanLocatedTime(c.Wall, c.Ext, locations[c.Old])
		inLocation := func() (result string) {
			defer func() {
				if p := recover(); p != nil {
					if p != "time: missing Location in call to Time.In" {
						panic(p)
					}
					result = "nil"
				}
			}()
			return wire(value.In(locations[c.Next]))
		}
		got := strings.Join([]string{wire(time.GoLeanStripMono(value)), wire(time.GoLeanSetLoc(value, locations[c.Next])),
			wire(value.UTC()), wire(value.Local()), inLocation(), fmt.Sprint(id(value.Location()))}, "|")
		if got != c.Want {
			t.Fatalf("row %d: Go %q want %q", row, got, c.Want)
		}
	}
}
