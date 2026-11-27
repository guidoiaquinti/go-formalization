package dag

import (
	"encoding/json"
	"os"
	"runtime"
	"testing"
)

func TestGoLeanRulesSyntaxError(t *testing.T) {
	data, err := os.ReadFile(os.Getenv("GOLEAN_DAG_SYNTAXERROR_INPUT"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Nil       bool
		Line      int
		Msg, Last []byte
		Want      []byte
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	for i, c := range cases {
		var parser *rulesParser
		if !c.Nil {
			parser = &rulesParser{lineno: c.Line, lastWord: string(c.Last), text: "unread input"}
		}
		got := goLeanSyntaxErrorPanic(parser, string(c.Msg))
		if c.Nil {
			runtimeErr, ok := got.(runtime.Error)
			if !ok || runtimeErr.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				t.Fatalf("case %d: got %T %[2]q, want nil receiver panic", i, got)
			}
		} else {
			panicValue, ok := got.(syntaxError)
			if !ok || string(panicValue) != string(c.Want) {
				t.Fatalf("case %d line=%d msg=%q last=%q: got %T %[5]q, want syntaxError %q",
					i, c.Line, c.Msg, c.Last, got, c.Want)
			}
		}
	}
}

func goLeanSyntaxErrorPanic(parser *rulesParser, msg string) (value any) {
	defer func() { value = recover() }()
	parser.syntaxError(msg)
	return nil
}
