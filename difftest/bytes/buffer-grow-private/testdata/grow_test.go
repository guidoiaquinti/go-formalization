package bytes_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestGoLeanBufferGrowPrivate(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_BUFFER_GROW_PRIVATE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var fixtures []struct {
		Input struct {
			Tag, Offset, Flag int
			Stored, Spare     []byte
			Count             int
		}
		Want string
	}
	if err := json.NewDecoder(f).Decode(&fixtures); err != nil {
		t.Fatal(err)
	}

	for i, f := range fixtures {
		c := f.Input
		got := bytes.GoLeanBufferGrowPrivate(c.Tag, c.Offset, c.Flag, c.Stored, c.Spare, c.Count)
		if got != f.Want {
			t.Fatalf("row %d: got %q want %q", i, got, f.Want)
		}
	}
}
