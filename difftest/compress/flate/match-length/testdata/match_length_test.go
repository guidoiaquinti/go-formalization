package flate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

func observedMatch(left, right []byte) (result string) {
	defer func() {
		if e := recover(); e != nil {
			failure, ok := e.(runtime.Error)
			if !ok || !strings.Contains(failure.Error(), "slice bounds out of range") {
				panic(e)
			}
			result = "bounds"
		}
	}()
	return fmt.Sprint(matchLen(left, right))
}

func TestGoLeanMatchLength(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_MATCH_LENGTH_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := json.NewDecoder(f)
	for i := 0; ; i++ {
		var c struct {
			Left, Right []byte
			RightLength int
			Want        string
		}
		err := input.Decode(&c)
		if err == io.EOF && i > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		if got := observedMatch(c.Left, c.Right[:c.RightLength:len(c.Right)]); got != c.Want {
			t.Fatalf("row %d: got %q want %q", i, got, c.Want)
		}
	}
}
