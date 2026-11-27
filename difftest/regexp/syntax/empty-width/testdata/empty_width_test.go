package syntax

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"
)

func TestGoLeanEmptyWidth(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_REGEXP_SYNTAX_EMPTY_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Before, After rune
			Present       bool
			Op            uint8
			Arg           uint32
			Want          string
		}
		if err := dec.Decode(&c); err == io.EOF {
			if row == 0 {
				t.Fatal("no empty width fixtures")
			}
			t.Logf("four pinned empty width declarations matched %d rows", row)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		var inst *Inst
		if c.Present {
			inst = &Inst{Op: InstOp(c.Op), Arg: c.Arg}
		}
		matched := catchEmpty(func() string { return fmt.Sprint(inst.MatchEmptyWidth(c.Before, c.After)) })
		op := catchEmpty(func() string { return fmt.Sprint(uint8(inst.op())) })
		got := fmt.Sprintf("%t|%t|%d|%s|%s", IsWordChar(c.Before), IsWordChar(c.After),
			uint8(EmptyOpContext(c.Before, c.After)), matched, op)
		if got != c.Want {
			t.Fatalf("row %d %+v: actual %q", row, c, got)
		}
	}
}

func catchEmpty(f func() string) (report string) {
	defer func() {
		if p := recover(); p != nil {
			if p == "unknown empty width arg" {
				report = "unknown"
				return
			}
			if err, ok := p.(runtime.Error); ok && err.Error() == "runtime error: invalid memory address or nil pointer dereference" {
				report = "nil"
				return
			}
			panic(p)
		}
	}()
	return f()
}
