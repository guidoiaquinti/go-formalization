package strings_test

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanGenericReplacerReplace(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_GENERIC_REPLACER_REPLACE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			OldNew [][]byte
			Text   []byte
			Want   string
		}
		if err := decoder.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("no replacement cases")
			}
			t.Logf("%d native constructor-built replacements agree", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		oldnew := make([]string, len(row.OldNew))
		for j, raw := range row.OldNew {
			oldnew[j] = string(raw)
		}
		value := strings.GoLeanGenericReplacerReplace(oldnew, string(row.Text))
		raw := make([]string, len(value))
		for j := range len(value) {
			raw[j] = strconv.Itoa(int(value[j]))
		}
		got := "bytes:" + strings.Join(raw, ",")
		if got != row.Want {
			t.Fatalf("row %d pairs=%q input=%x: got %q want %q", i, oldnew, row.Text, got, row.Want)
		}
	}
}
