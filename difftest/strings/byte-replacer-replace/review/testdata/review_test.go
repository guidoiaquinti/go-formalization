package strings_test

import (
	"bytes"
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGoLeanReviewByteReplacer(t *testing.T) {
	data, err := os.ReadFile(os.Getenv("GOLEAN_REVIEW_BYTE_REPLACER"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Table       *[256]byte
		Input, Want []byte
		Panic       bool
	}
	if err := json.Unmarshal(data, &cases); err != nil || len(cases) != 8 {
		t.Fatalf("invalid fixture set: %v", err)
	}
	for i, c := range cases {
		got, panicked := reviewReplaceBytes(c.Table, c.Input)
		if panicked != c.Panic || !bytes.Equal([]byte(got), c.Want) {
			t.Fatalf("fixture %d: got %q panic=%t, want %q panic=%t", i, got, panicked, c.Want, c.Panic)
		}
	}
}

func reviewReplaceBytes(mapping *[256]byte, input []byte) (got string, panicked bool) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			panicked = true
		}
	}()
	got = strings.GoLeanByteReplacerReplace(mapping, string(input))
	return got, false
}
