package strings_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanGenericReplacerLookup(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_GENERIC_REPLACER_LOOKUP_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			OldNew     [][]byte
			Text       []byte
			IgnoreRoot bool
			Want       string
		}
		if err := decoder.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("no lookup cases")
			}
			t.Logf("%d native constructor-built lookups agree", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		oldnew := make([]string, len(row.OldNew))
		for j, raw := range row.OldNew {
			oldnew[j] = string(raw)
		}
		value, n, found := strings.GoLeanGenericReplacerLookup(oldnew, string(row.Text), row.IgnoreRoot)
		tag := "miss"
		if found {
			tag = "hit"
		}
		raw := make([]string, len(value))
		for j := range len(value) {
			raw[j] = strconv.Itoa(int(value[j]))
		}
		got := fmt.Sprintf("%s:%d:%s", tag, n, strings.Join(raw, ","))
		if got != row.Want {
			t.Fatalf("row %d pairs=%q input=%x ignoreRoot=%t: got %q want %q", i, oldnew, row.Text, row.IgnoreRoot, got, row.Want)
		}
	}
}
