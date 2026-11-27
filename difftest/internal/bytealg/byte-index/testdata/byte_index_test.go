package bytealg_test

import (
	"encoding/json"
	"internal/bytealg"
	"io"
	"os"
	"testing"
)

func TestGoLeanBytealgByteIndex(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTEALG_BYTE_INDEX_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			Bytes       []byte
			Target      byte
			Offset      int
			First, Last int
		}
		if err := rows.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("no byte-index cases")
			}
			t.Logf("matched %d inputs through four native byte-index entry points", i)
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
		want := [4]int{row.First, row.First, row.Last, row.Last}
		if got := bytealg.GoLeanBytealgByteIndices(slice, text, row.Target); got != want {
			t.Fatalf("row %d bytes=%x target=%02x offset=%d: got %v want %v", i, slice, row.Target, row.Offset, got, want)
		}
	}
}
