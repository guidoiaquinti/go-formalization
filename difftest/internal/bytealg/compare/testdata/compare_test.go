package bytealg_test

import (
	"encoding/json"
	"internal/bytealg"
	"io"
	"os"
	"testing"
)

func TestGoLeanBytealgCompare(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_BYTEALG_COMPARE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows := json.NewDecoder(f)
	for i := 0; ; i++ {
		var row struct {
			Left, Right                    []byte
			LeftOffset, RightOffset, Order int
			Alias                          bool
		}
		if err := rows.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("missing comparison cases")
			}
			t.Logf("%d inputs matched through four native comparison entry points", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		makeView := func(raw []byte, offset int, poison byte) ([]byte, string) {
			if raw == nil {
				return nil, ""
			}
			backing := make([]byte, offset+len(raw)+48)
			for i := range backing {
				backing[i] = poison
			}
			copy(backing[offset:], raw)
			end := offset + len(raw)
			return backing[offset:end], string(backing)[offset:end]
		}
		left, a := makeView(row.Left, row.LeftOffset, 0x55)
		right, b := makeView(row.Right, row.RightOffset, 0xaa)
		if row.Alias {
			right, b = left, a
		}
		orders, equal := bytealg.GoLeanBytealgCompare(left, right, a, b)
		want := [3]int{row.Order, row.Order, row.Order}
		if orders != want || equal != (row.Order == 0) {
			t.Fatalf("row %d a=%x b=%x: orders=%v Equal=%t want=%d", i, left, right, orders, equal, row.Order)
		}
	}
}
