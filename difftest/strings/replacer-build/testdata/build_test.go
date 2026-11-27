package strings_test

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGoLeanReplacerBuild(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_REPLACER_BUILD_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			Pairs [][]byte
			Want  string
		}
		if err := rows.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("missing constructor fixtures")
			}
			t.Logf("%d native build snapshots agree", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		pairs := make([]string, len(row.Pairs))
		for j, raw := range row.Pairs {
			pairs[j] = string(raw)
		}
		got, actual, delegate := strings.GoLeanReplacerBuild(pairs)
		if !reflect.DeepEqual(actual, delegate) {
			t.Fatalf("row %d: delegated state differs from the direct constructor", i)
		}
		if got != row.Want {
			t.Fatalf("row %d pairs=%q: Go=%s want=%s", i, pairs, got, row.Want)
		}
	}
}
