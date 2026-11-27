package bytes_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestGoLeanReaderCopy(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_READER_COPY_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var rows []struct {
		Input struct {
			Tag, Previous        int
			Cursor, Offset       int64
			Sharing, Size        int
			Raw, Spare, Separate []byte
		}
		Want string
	}
	if err := json.NewDecoder(f).Decode(&rows); err != nil {
		t.Fatal(err)
	}
	for i, row := range rows {
		c := row.Input
		if got := bytes.GoLeanReaderCopy(c.Tag, c.Previous, c.Cursor, c.Offset, c.Sharing, c.Size,
			c.Raw, c.Spare, c.Separate); got != row.Want {
			t.Fatalf("row %d (%+v): got %q want %q", i, c, got, row.Want)
		}
	}
}
