package bytes_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestGoLeanGrowSlice(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_GROW_SLICE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var fixtures []struct {
		Input struct {
			Raw             []byte
			Capacity, Count int
			NilInput        bool
		}
		Want string
	}
	if err := json.NewDecoder(f).Decode(&fixtures); err != nil {
		t.Fatal(err)
	}
	for i, row := range fixtures {
		c := row.Input
		if got := bytes.GoLeanGrowSlice(c.Raw, c.Capacity, c.Count, c.NilInput); got != row.Want {
			t.Fatalf("row %d: got %q want %q", i, got, row.Want)
		}
	}
}
