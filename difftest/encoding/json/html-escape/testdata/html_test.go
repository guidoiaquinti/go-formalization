package json_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanHTMLEscape(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_JSON_HTML_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Present      bool
			Backing, Src []byte
			Length       int
			Want         string
		}
		if err := dec.Decode(&c); err == io.EOF {
			if row == 0 {
				t.Fatal("missing HTML fixtures")
			}
			t.Logf("two pinned HTML escape entries matched %d rows", row)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		original := slices.Clone(c.Backing)
		scratch := slices.Clone(c.Backing)
		out := json.GoLeanAppendHTMLEscape(scratch[:c.Length:len(scratch)], c.Src)
		reuse := len(scratch) > 0 && cap(out) > 0 && &scratch[0] == &out[:cap(out)][0]
		var target *bytes.Buffer
		if c.Present {
			target = bytes.NewBuffer(original[:c.Length:len(original)])
		}
		public := observe(func() string {
			json.HTMLEscape(target, c.Src)
			if target.Len() == 0 {
				return "slice:"
			}
			return wire(target.Bytes())
		})
		got := fmt.Sprintf("%s/%s/%t|%s", wire(out), wire(scratch), reuse, public)
		if got != c.Want {
			t.Fatalf("row %d %+v: got %q", row, c, got)
		}
	}
}

func observe(run func() string) (out string) {
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(runtime.Error); !ok || e.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(p)
			}
			out = "panic"
		}
	}()
	return run()
}

func wire(data []byte) string {
	if data == nil {
		return "nil"
	}
	parts := make([]string, len(data))
	for i, v := range data {
		parts[i] = strconv.Itoa(int(v))
	}
	return "slice:" + strings.Join(parts, ",")
}
