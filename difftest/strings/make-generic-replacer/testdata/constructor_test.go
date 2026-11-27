package strings_test

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

func TestGoLeanMakeGenericReplacer(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MAKE_GENERIC_REPLACER_CASES"))
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
			t.Logf("%d native trie snapshots agree", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		pairs := make([]string, len(row.Pairs))
		for j, raw := range row.Pairs {
			pairs[j] = string(raw)
		}
		if got := strings.GoLeanMakeGenericReplacer(pairs); got != row.Want {
			t.Fatalf("row %d pairs=%q: Go=%s want=%s", i, pairs, got, row.Want)
		}
	}
}
