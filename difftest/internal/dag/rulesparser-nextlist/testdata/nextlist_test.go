package dag

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
)

func TestGoLeanNextList(t *testing.T) {
	data, err := os.ReadFile(os.Getenv("GOLEAN_DAG_NEXTLIST_INPUT"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Present    bool
		Line       int
		Last, Text []byte
		Expected   string
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	for i, c := range cases {
		var p *rulesParser
		if c.Present {
			p = &rulesParser{lineno: c.Line, lastWord: string(c.Last), text: string(c.Text)}
		}
		if got := goLeanNextList(p); got != c.Expected {
			t.Fatalf("case %d line %d last %q text %q: Go=%q expected=%q", i, c.Line, c.Last, c.Text, got, c.Expected)
		}
	}
}

func goLeanNextList(p *rulesParser) (out string) {
	defer func() {
		if e := recover(); e != nil {
			switch e := e.(type) {
			case syntaxError:
				out = goLeanListRow("syntax", nil, string(e), p)
			case runtime.Error:
				if p == nil && e.Error() == "runtime error: invalid memory address or nil pointer dereference" {
					out = "nil"
				} else {
					panic(e)
				}
			default:
				panic(e)
			}
		}
	}()
	names, token := p.nextList()
	if len(names) == 0 && names != nil {
		return "unexpected allocated empty list"
	}
	return goLeanListRow("return", names, token, p)
}

func goLeanListRow(kind string, names []string, token string, p *rulesParser) string {
	encode := func(s string) []byte {
		row := strconv.AppendInt(nil, int64(len(s)), 10)
		for i := range len(s) {
			row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
		}
		return row
	}
	row := kind
	if kind == "return" {
		row += " " + strconv.Itoa(len(names))
		for _, name := range names {
			row += " " + string(encode(name))
		}
	}
	return fmt.Sprintf("%s %s %d %s %s", row, encode(token), p.lineno, encode(p.lastWord), encode(p.text))
}
