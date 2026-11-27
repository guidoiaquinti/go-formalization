package types

import (
	"encoding/json"
	"fmt"
	"go/constant"
	"go/token"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type objectFixture struct {
	Mode, Pos, Pkg, Type, Value, Import, Sig, Scope int
	Name                                            []byte
	Want                                            string
}

func TestGoLeanObjectConstructors(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_TYPES_OBJECT_CONSTRUCT"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	pkg, imported := NewPackage("package", "same"), NewPackage("imported", "same")
	scope := NewScope(nil, 0, 0, "body")
	boolValue := constant.MakeBool(true)
	nilValue := reflect.Zero(reflect.TypeOf(constant.MakeString("test"))).Interface().(constant.Value)
	for row := 0; ; row++ {
		var c objectFixture
		if err := dec.Decode(&c); err != nil {
			if err == io.EOF && row > 0 {
				t.Logf("%d pinned native cases per declaration", row)
				return
			}
			t.Fatal(err)
		}
		pos, name := token.Pos(c.Pos), string(c.Name)
		owner := []*Package{nil, pkg}[c.Pkg]
		dependency := []*Package{nil, imported}[c.Import]
		value := []constant.Value{nil, nilValue, boolValue}[c.Value]
		typ := []Type{nil, (*Array)(nil), Typ[Int]}[c.Type]
		sig := &Signature{scope: []*Scope{nil, scope}[c.Scope]}
		if c.Sig == 0 {
			sig = nil
		}
		cn := NewConst(pos, owner, name, typ, value)
		pn := NewPkgName(pos, owner, name, dependency)
		tn := NewTypeName(pos, owner, name, typ)
		label := NewLabel(pos, owner, name)
		fn := NewFunc(pos, owner, name, sig)
		if cn == NewConst(pos, owner, name, typ, value) || pn == NewPkgName(pos, owner, name, dependency) ||
			tn == NewTypeName(pos, owner, name, typ) || label == NewLabel(pos, owner, name) || fn == NewFunc(pos, owner, name, sig) {
			t.Fatal("constructor reused identity")
		}
		var funcType Type
		if sig != nil {
			funcType = sig
		}
		for _, pair := range []struct {
			got  *object
			want Type
		}{
			{&cn.object, typ}, {&pn.object, Typ[Invalid]}, {&tn.object, typ}, {&label.object, Typ[Invalid]}, {&fn.object, funcType},
		} {
			want := object{nil, pos, owner, name, pair.want, 0, token.NoPos}
			if *pair.got != want {
				t.Fatalf("row %d: constructor field identity mismatch", row)
			}
		}
		if cn.val != value || pn.imported != dependency || fn.origin != nil || fn.hasPtrRecv_ || fn.nointerface {
			t.Fatalf("row %d: extra state mismatch", row)
		}
		typeName := []string{"nil", "typedNil1", "11"}[c.Type]
		valueName := []string{"nil", "typedNil2", "12"}[c.Value]
		importName := []string{"nil", "22"}[c.Import]
		funcName := "nil"
		if sig != nil {
			funcName = "10"
		}
		observation := []string{typeName + ":" + valueName, "9:" + importName, typeName, fmt.Sprintf("9:%t", label.used), funcName}
		hold := new(Func)
		if c.Mode == 0 {
			cn, pn, fn = nil, nil, nil
		} else {
			if c.Mode == 2 {
				fn.origin = hold
			}
			switch c.Sig {
			case 2:
				fn.typ = (*Signature)(nil)
			case 3:
				fn.typ = Typ[Int]
			}
		}
		observation = append(observation,
			readConstructorResult(func() string {
				if cn.Val() != value {
					return "wrong"
				}
				return valueName
			}),
			readConstructorResult(func() string {
				if pn.Imported() != dependency {
					return "wrong"
				}
				return importName
			}),
			readConstructorResult(func() string {
				got := fn.Signature()
				switch {
				case got == nil:
					return "typedNil10"
				case got == sig:
					return "10"
				case fn.typ == nil && got != fn.Signature() && reflect.DeepEqual(got, new(Signature)):
					return "fresh"
				default:
					return "wrong"
				}
			}),
			readConstructorResult(func() string {
				switch fn.Scope() {
				case nil:
					return "nil"
				case scope:
					return "50"
				default:
					return "wrong"
				}
			}),
			readConstructorResult(func() string {
				switch fn.Origin() {
				case fn:
					return "35"
				case hold:
					return "36"
				default:
					return "wrong"
				}
			}),
			readConstructorResult(func() string {
				switch fn.Pkg() {
				case nil:
					return "nil"
				case pkg:
					return "21"
				default:
					return "wrong"
				}
			}))
		if got := strings.Join(observation, "|"); got != c.Want || got != constructorOracle(c) {
			t.Fatalf("row %d: pinned %q Lean %q independently %q", row, got, c.Want, constructorOracle(c))
		}
	}
}

func readConstructorResult(f func() string) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); !ok {
				panic(p)
			}
			result = "panic"
		}
	}()
	return f()
}

func constructorOracle(c objectFixture) string {
	typeName := []string{"nil", "typedNil1", "11"}[c.Type]
	value := []string{"nil", "typedNil2", "12"}[c.Value]
	imported := []string{"nil", "22"}[c.Import]
	constructorType := "10"
	if c.Sig == 0 {
		constructorType = "nil"
	}
	out := []string{typeName + ":" + value, "9:" + imported, typeName, "9:false", constructorType}
	if c.Mode == 0 {
		return strings.Join(out, "|") + strings.Repeat("|panic", 6)
	}
	scope := "panic"
	if c.Sig == 1 {
		scope = []string{"nil", "50"}[c.Scope]
	}
	return strings.Join(append(out, value, imported, []string{"fresh", "10", "typedNil10", "panic"}[c.Sig],
		scope, []string{"35", "36"}[c.Mode-1], []string{"nil", "21"}[c.Pkg]), "|")
}
