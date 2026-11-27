package types

import (
	"encoding/json"
	"fmt"
	"go/token"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanPackageFields(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_TYPES_PACKAGE_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	wire := func(raw string) string {
		parts := make([]string, len(raw))
		for i, b := range []byte(raw) {
			parts[i] = strconv.Itoa(int(b))
		}
		return strings.Join(parts, ",")
	}
	read := func(f func() string) (result string) {
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
	for row := 0; ; row++ {
		var c struct {
			Mode                             int
			Complete                         bool
			Path, Name, Version, Replacement []byte
			Want                             string
		}
		err := dec.Decode(&c)
		if err == io.EOF && row > 0 {
			t.Logf("%d pinned rows per package declaration", row)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		p := NewPackage(string(c.Path), string(c.Name))
		scope := p.Scope()
		if scope == nil || scope.Parent() != Universe || scope.Pos() != token.NoPos ||
			scope.End() != token.NoPos || scope.Len() != 0 || p.Imports() != nil ||
			scope.comment != fmt.Sprintf("package %q", c.Path) {
			t.Fatal("NewPackage did not initialize an empty package scope")
		}
		fresh := NewPackage(string(c.Path), string(c.Name))
		if fresh == p || fresh.Scope() == scope {
			t.Fatal("NewPackage did not return distinct package and scope")
		}
		fields := func(pkg *Package) string {
			scopeWire := "wrong"
			switch pkg.Scope() {
			case Universe:
				scopeWire = "0"
			case scope:
				scopeWire = "31"
			case nil:
				scopeWire = "nil"
			}
			return strings.Join([]string{
				read(func() string { return wire(pkg.Path()) }),
				read(func() string { return wire(pkg.Name()) }),
				scopeWire,
				read(func() string { return strconv.FormatBool(pkg.Complete()) }),
				read(func() string { return wire(pkg.GoVersion()) }),
			}, "|")
		}
		created := fields(p)
		if c.Mode == 0 {
			p = nil
		} else {
			p.complete, p.goVersion = c.Complete, string(c.Version)
			if c.Mode == 2 {
				p.scope = nil
			}
		}
		before := fields(p)
		var renameTarget, completeTarget *Package
		if p != nil {
			r, m := *p, *p
			renameTarget, completeTarget = &r, &m
		}
		renamed := read(func() string { renameTarget.SetName(string(c.Replacement)); return fields(renameTarget) })
		completed := read(func() string { completeTarget.MarkComplete(); return fields(completeTarget) })
		result := strings.Join([]string{created, before, renamed, completed}, "#")
		if result != c.Want {
			t.Fatalf("row %d: Go %q, want %q", row, result, c.Want)
		}
	}
}
