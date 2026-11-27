package bytes_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestGoLeanReaderPosition(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_READER_POSITION_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var rows []struct {
		Input struct {
			Tag, Previous, Whence, Script, Count int
			Cursor, Offset                       int64
			Raw, Spare                           []byte
		}
		Want string
	}
	if err := json.NewDecoder(f).Decode(&rows); err != nil {
		t.Fatal(err)
	}
	for i, row := range rows {
		c := row.Input
		if got := bytes.GoLeanReaderPosition(c.Tag, c.Previous, c.Cursor, c.Offset, c.Whence,
			c.Script, c.Count, c.Raw, c.Spare); got != row.Want {
			t.Fatalf("row %d (%+v): got %q want %q", i, c, got, row.Want)
		}
	}
}
