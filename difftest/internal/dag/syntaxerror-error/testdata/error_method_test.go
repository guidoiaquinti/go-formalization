package dag

import (
	"encoding/json"
	"os"
	"testing"
)

func TestGoLeanSyntaxErrorMethod(t *testing.T) {
	data, err := os.ReadFile(os.Getenv("GOLEAN_DAG_ERROR_METHOD_INPUT"))
	if err != nil {
		t.Fatal(err)
	}
	var cases [][]byte
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	for i, input := range cases {
		want := string(input)
		if got := syntaxError(want).Error(); got != want {
			t.Fatalf("case %d input=%q: Go=%q want=%q", i, input, got, want)
		}
	}
}
