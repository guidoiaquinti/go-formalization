package zip

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestGoLeanZipZone(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ZIP_ZONE_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := bufio.NewReader(f)
	for row := 0; ; row++ {
		var offset time.Duration
		var seconds int
		var cached bool
		n, err := fmt.Fscan(input, &offset, &seconds, &cached)
		if n == 0 && err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		loc := timeZone(offset)
		if got := loc == timeZone(offset); got != cached {
			t.Fatalf("offset %d: cache=%t want=%t", offset, got, cached)
		}
		for _, epoch := range []int64{-2208988800, 0, 1719792000, 4102444800} {
			instant := time.Unix(epoch, 0).In(loc)
			name, got := instant.Zone()
			start, end := instant.ZoneBounds()
			if name != "" || loc.String() != "" || got != seconds || instant.IsDST() || !start.IsZero() || !end.IsZero() {
				t.Fatalf("offset %d epoch %d: zone=%q/%d want empty/%d", offset, epoch, name, got, seconds)
			}
		}
	}
}
