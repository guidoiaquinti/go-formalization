package strings_test

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanReplacerBuildOnce(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_REPLACER_BUILD_ONCE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			Pairs                     [][]byte
			Present, OnceZero, Repeat bool
			Want                      string
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
		got, once := strings.GoLeanReplacerBuildOnce(pairs, row.Present, row.OnceZero, row.Repeat)
		if once != nil {
			state, implementation, _ := strings.Cut(got, ":")
			got = state + ":" + strconv.FormatBool(reflect.ValueOf(once).Elem().IsZero()) + ":" + implementation
		}
		if got != row.Want {
			t.Fatalf("row %d pairs=%q: Go=%s want=%s", i, pairs, got, row.Want)
		}
	}
}
