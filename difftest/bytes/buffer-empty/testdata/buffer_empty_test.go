package bytes_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestGoLeanBufferEmpty(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_BUFFER_EMPTY_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var rows []struct {
		Input struct {
			NilReceiver, NilStorage bool
			Length, Spare, Offset   int
		}
		Want string
	}
	if err := json.NewDecoder(f).Decode(&rows); err != nil {
		t.Fatal(err)
	}
	for i, row := range rows {
		c := row.Input
		if got := bytes.GoLeanBufferEmpty(c.NilReceiver, c.NilStorage, c.Length, c.Spare, c.Offset); got != row.Want {
			t.Fatalf("row %d %+v: got %s want %s", i, c, got, row.Want)
		}
	}
}
