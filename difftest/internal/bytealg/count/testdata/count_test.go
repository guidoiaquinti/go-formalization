package bytealg_test

import (
	"encoding/json"
	"internal/bytealg"
	"io"
	"os"
	"testing"
)

func TestGoLeanBytealgCount(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTEALG_COUNT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			Bytes        []byte
			Target       byte
			Offset, Want int
		}
		if err := rows.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("no count cases")
			}
			t.Logf("matched %d inputs through four native count entry points", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		var slice []byte
		var text string
		if row.Offset >= 0 {
			data := make([]byte, row.Offset+len(row.Bytes)+48)
			for j := range data {
				data[j] = row.Target
			}
			copy(data[row.Offset:], row.Bytes)
			end := row.Offset + len(row.Bytes)
			slice = data[row.Offset:end]
			text = string(data)[row.Offset:end]
		}
		want := [4]int{row.Want, row.Want, row.Want, row.Want}
		if got := bytealg.GoLeanBytealgCounts(slice, text, row.Target); got != want {
			t.Fatalf("row %d bytes=%x target=%02x offset=%d: got %v want %v", i, slice, row.Target, row.Offset, got, want)
		}
	}
}
