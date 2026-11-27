package time_test

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestGoLeanClockMutation(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TIME_CLOCK_MUTATION_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := gob.NewDecoder(bufio.NewReader(f))
	locations := []*time.Location{nil, time.FixedZone("alias", 1000), time.FixedZone("alias", 1000)}
	for row := 0; ; row++ {
		var c struct {
			Wall, Ext, Argument uint64
			Location            int
			Want                string
		}
		if err := in.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row > 0 {
				t.Logf("%d pinned observations for both mutations", row)
				return
			}
			t.Fatal(err)
		}
		first, second := time.GoLeanClockMutations(c.Wall, c.Ext, locations[c.Location], c.Argument)
		format := func(v time.Time) string {
			wall, ext, loc := time.GoLeanClockMutationParts(v)
			if loc != locations[c.Location] {
				t.Fatalf("row %d location pointer changed", row)
			}
			return fmt.Sprintf("%d,%d,%d", wall, ext, c.Location)
		}
		if got := format(first) + "|" + format(second); got != c.Want {
			t.Fatalf("row %d raw mutation: Go %q expect %q", row, got, c.Want)
		}
	}
}
