package types

import (
	"encoding/json"
	"fmt"
	"go/token"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGoLeanVariableState(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_TYPES_VAR_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	observe := func(f func() string) (result string) {
		defer func() {
			if p := recover(); p != nil {
				if err, ok := p.(runtime.Error); !ok ||
					err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
					panic(p)
				}
				result = "panic"
			}
		}()
		return f()
	}
	pkg := NewPackage("p", "p")
	typeValues := []Type{nil, (*Array)(nil), Typ[Int], NewSlice(Typ[Byte])}
	for row := 0; ; row++ {
		var c struct {
			Mode, Kind, Replacement, Pos, Pkg, Typ int
			Embedded                               bool
			Name                                   []byte
			Want                                   string
		}
		err := dec.Decode(&c)
		if err == io.EOF && row > 0 {
			t.Logf("%d pinned variable rows", row)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		var selected *Package
		if c.Pkg != 0 {
			selected = pkg
		}
		pos, name, typ := token.Pos(c.Pos), string(c.Name), typeValues[c.Typ]
		kind := VarKind(c.Kind)
		constructors := []*Var{newVar(kind, pos, selected, name, typ), NewVar(pos, selected, name, typ),
			NewParam(pos, selected, name, typ), NewField(pos, selected, name, typ, c.Embedded)}
		again := []*Var{newVar(kind, pos, selected, name, typ), NewVar(pos, selected, name, typ),
			NewParam(pos, selected, name, typ), NewField(pos, selected, name, typ, c.Embedded)}
		var fields []string
		for i, v := range constructors {
			if v == again[i] || v.pos != pos || v.pkg != selected || v.name != name || v.typ != typ ||
				v.parent != nil || v.order_ != 0 || v.scopePos_ != token.NoPos {
				t.Fatalf("row %d constructor %d has wrong object metadata or identity", row, i)
			}
			fields = append(fields, fmt.Sprintf("%d:%t:%t", v.kind, v.embedded, v.origin == nil))
		}
		var v *Var
		origin := new(Var)
		if c.Mode != 0 {
			v = constructors[0]
			v.embedded = c.Embedded
			if c.Mode == 2 {
				v.origin = origin
			}
		}
		identity := func(got *Var) string {
			switch got {
			case origin:
				return "32"
			case v:
				return "31"
			default:
				return "wrong"
			}
		}
		fields = append(fields, observe(func() string { return fmt.Sprint(uint8(v.Kind())) }),
			observe(func() string { return fmt.Sprint(v.Anonymous()) }),
			observe(func() string { return fmt.Sprint(v.Embedded()) }),
			observe(func() string { return fmt.Sprint(v.IsField()) }),
			observe(func() string { return identity(v.Origin()) }),
			observe(func() string {
				before := object{}
				if v != nil {
					before = v.object
				}
				v.SetKind(VarKind(c.Replacement))
				if v.object != before || v.embedded != c.Embedded ||
					(v.origin != nil) != (c.Mode == 2) || v.origin != nil && v.origin != origin {
					t.Fatal("SetKind changed unrelated variable state")
				}
				return fmt.Sprintf("%d:%t:%s", v.kind, v.embedded, identity(v.Origin()))
			}), kind.String())
		if got := strings.Join(fields, "|"); got != c.Want {
			t.Fatalf("row %d Go %q, want %q", row, got, c.Want)
		}
	}
}
