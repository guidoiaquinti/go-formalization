package bytes_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestGoLeanReaderSetup(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_READER_SETUP_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var rows []struct {
		Input struct {
			Tag, Offset, Previous int
			Cursor                int64
			Raw, Spare            []byte
		}
		Want string
	}
	if err := json.NewDecoder(f).Decode(&rows); err != nil {
		t.Fatal(err)
	}
	for i, row := range rows {
		c := row.Input
		if got := bytes.GoLeanReaderSetup(c.Tag, c.Offset, c.Cursor, c.Previous, c.Raw, c.Spare); got != row.Want {
			t.Fatalf("row %d: got %q want %q", i, got, row.Want)
		}
	}
}
