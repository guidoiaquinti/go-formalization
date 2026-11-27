package bytes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"
)

func TestGoLeanReadFrom(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BUFFER_READ_FROM_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			Input struct {
				Tag, Len, Off, Cap, Flag int
				NilReader                bool
				Steps                    []bytes.GoLeanReadStep
			}
			Want string
		}
		if err := decoder.Decode(&row); errors.Is(err, io.EOF) {
			if i == 0 {
				t.Fatal("no fixtures")
			}
			t.Logf("pinned ReadFrom matched %d rows", i)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		c := row.Input
		if got := bytes.GoLeanReadFrom(c.Tag, c.Len, c.Off, c.Cap, c.Flag, c.NilReader, c.Steps); got != row.Want {
			t.Fatalf("row %d %+v: got %q want %q", i, c, got, row.Want)
		}
	}
}
