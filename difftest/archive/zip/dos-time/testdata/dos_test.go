package zip

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestGoLeanDOSTime(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ZIP_DOS_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := json.NewDecoder(f)
	for row := 0; ; row++ {
		var c struct {
			Date, Clock   uint16
			Seconds       int64
			Nanos, Offset int
			Want          string
		}
		err := in.Decode(&c)
		if err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		instant := time.Unix(c.Seconds, int64(c.Nanos)).In(time.FixedZone("test", c.Offset))
		direct := msDosTimeToTime(c.Date, c.Clock)
		localDate, localTime := timeToMsDosTime(instant)
		header := FileHeader{ModifiedDate: c.Date, ModifiedTime: c.Clock, Modified: instant.Add(time.Hour), Name: "example", Comment: "kept", Method: Deflate}
		modified := header.ModTime()
		if direct.Location() != time.UTC || modified.Location() != time.UTC || direct.Nanosecond() != 0 || modified.Nanosecond() != 0 {
			t.Fatal("decoded timestamp was not integral UTC")
		}
		header.SetModTime(instant)
		got := fmt.Sprintf("%d|%d:%d|%d|%d:%d:%d:%d", direct.Unix(), localDate, localTime, modified.Unix(),
			header.Modified.Unix(), header.Modified.Nanosecond(), header.ModifiedDate, header.ModifiedTime)
		if got != c.Want || header.Modified.Location() != time.UTC ||
			header.Name != "example" || header.Comment != "kept" || header.Method != Deflate {
			t.Fatalf("row %d: pinned Go %q want %q", row, got, c.Want)
		}
	}
}
