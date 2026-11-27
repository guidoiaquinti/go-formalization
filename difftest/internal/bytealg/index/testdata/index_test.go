package bytealg_test

import (
	"bytes"
	"encoding/json"
	"internal/bytealg"
	"io"
	"os"
	"testing"
)

func TestGoLeanBytealgIndex(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTEALG_INDEX_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Input, Needle                   []byte
			InputOffset, NeedleOffset, Want int
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("no index cases")
			}
			t.Logf("matched %d native rows with active MaxLen=%d", i, bytealg.MaxLen)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		if len(c.Needle) < 2 || len(c.Needle) > bytealg.MaxLen || len(c.Input) < len(c.Needle) {
			t.Fatalf("row %d violates native domain with len=%d, needle=%d, MaxLen=%d",
				i, len(c.Input), len(c.Needle), bytealg.MaxLen)
		}
		makeView := func(raw []byte, offset int) ([]byte, string) {
			storage := bytes.Repeat([]byte{0x80}, offset+len(raw)+64)
			copy(storage[offset:], raw)
			end := offset + len(raw)
			copy(storage[end:], c.Needle)
			return storage[offset:end], string(storage)[offset:end]
		}
		haystack, text := makeView(c.Input, c.InputOffset)
		needle, pattern := makeView(c.Needle, c.NeedleOffset)
		a, b := bytealg.GoLeanBytealgIndex(haystack, needle, text, pattern)
		if a != c.Want || b != c.Want {
			t.Fatalf("row %d input=%x needle=%x: Index=%d IndexString=%d want=%d", i, c.Input, c.Needle, a, b, c.Want)
		}
	}
}
