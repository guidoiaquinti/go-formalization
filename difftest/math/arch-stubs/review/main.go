package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	path := filepath.Join(runtime.GOROOT(), "src", "math", "stubs.go")
	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	expr, err := constraint.Parse("//go:build !s390x")
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Contains(source, []byte("//go:build !s390x\n")) ||
		expr.Eval(func(tag string) bool { return tag == "s390x" }) ||
		!expr.Eval(func(string) bool { return false }) {
		log.Fatal("unexpected stub build constraint")
	}
	file, err := parser.ParseFile(token.NewFileSet(), path, source, 0)
	if err != nil {
		log.Fatal(err)
	}
	flags := map[string]bool{}
	var names []string
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok != token.CONST {
				continue
			}
			for _, spec := range d.Specs {
				v := spec.(*ast.ValueSpec)
				flag, ok := v.Values[0].(*ast.Ident)
				if ok && flag.Name == "false" {
					flags[v.Names[0].Name] = true
				}
			}
		case *ast.FuncDecl:
			name := d.Name.Name
			if !strings.HasPrefix(name, "arch") {
				continue
			}
			if !flags["haveArch"+strings.TrimPrefix(name, "arch")] || d.Body == nil || len(d.Body.List) != 1 {
				log.Fatalf("%s: unexpected availability or body length", name)
			}
			statement, ok := d.Body.List[0].(*ast.ExprStmt)
			if !ok {
				log.Fatal(name, ": expected panic expression")
			}
			call, ok := statement.X.(*ast.CallExpr)
			if !ok || len(call.Args) != 1 {
				log.Fatal(name, ": invalid panic call")
			}
			fn, ok := call.Fun.(*ast.Ident)
			arg, literal := call.Args[0].(*ast.BasicLit)
			if !ok || fn.Name != "panic" || !literal || arg.Kind != token.STRING || arg.Value != `"not implemented"` {
				log.Fatal(name, ": invalid panic payload")
			}
			names = append(names, name)
		}
	}
	if len(names) != 25 {
		log.Fatalf("found %d stubs", len(names))
	}
	values := [][2]uint64{{0, 0}, {1, 1 << 63}, {0x7ff0000000000000, ^uint64(0)},
		{0x7ff8000000000082, 0x3ff0000000000000}, {0xffefffffffffffff, 0x7ff0000000000000}}
	var input bytes.Buffer
	for index := range names {
		for _, v := range values {
			fmt.Fprintln(&input, 0, index, v[0], v[1])
		}
		fmt.Fprintln(&input, 1, index, 0, 0)
	}
	for i, row := range check.LeanRows("godiff-math-arch-stubs", len(names)*(len(values)+1), &input) {
		want := "false:panic:not implemented"
		if i%(len(values)+1) == len(values) {
			want = "outside-scope"
		}
		if row != want {
			log.Fatalf("%s row %d: %q, want %q", names[i/(len(values)+1)], i, row, want)
		}
	}
	fmt.Printf("Audited %d pinned declarations and %d Lean profile rows.\n", len(names), len(names)*(len(values)+1))
}
