package sort_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestGoLeanPivotHelpers(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SORT_PIVOT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows := json.NewDecoder(f)
	for i := 0; ; i++ {
		var row struct {
			Kind, A, B, C, Initial int
			Script                 []bool
			Want                   string
		}
		if err := rows.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("missing pivot cases")
			}
			t.Logf("%d fixtures matched in interface and function variants", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		for _, variant := range []bool{false, true} {
			var calls []string
			first, second, swaps := sort.GoLeanPivotHelpers(variant, row.Kind, row.A, row.B, row.C, row.Initial, func(a, b int) bool {
				step := len(calls)
				calls = append(calls, fmt.Sprintf("%d,%d", a, b))
				return step < len(row.Script) && row.Script[step]
			})
			got := fmt.Sprintf("%d:%d:%d:%s", first, second, swaps, strings.Join(calls, ";"))
			if got != row.Want {
				t.Fatalf("row %d variant=%t: native %s want %s", i, variant, got, row.Want)
			}
		}
	}
}
