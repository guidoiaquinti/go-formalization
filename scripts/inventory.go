// Enumerate named declarations selected by the pinned Go toolchain.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Package struct {
	ImportPath string
	Dir        string
	GoFiles    []string
}

type Declaration struct {
	Package string
	Symbol  string
	File    string
	Line    int
}

// A method belongs to its named base type, including pointer and generic receivers.
func receiver(expr ast.Expr) string {
	switch typ := expr.(type) {
	case *ast.Ident:
		return typ.Name
	case *ast.StarExpr:
		return receiver(typ.X)
	case *ast.IndexExpr:
		return receiver(typ.X)
	case *ast.IndexListExpr:
		return receiver(typ.X)
	case *ast.ParenExpr:
		return receiver(typ.X)
	default:
		panic(fmt.Sprintf("unexpected receiver %T", expr))
	}
}

func scan(patterns []string) error {
	if runtime.Version() != os.Getenv("GOTOOLCHAIN") {
		return fmt.Errorf("inventory built with %s, expected %s", runtime.Version(), os.Getenv("GOTOOLCHAIN"))
	}
	args := append([]string{"list", "-json=ImportPath,Dir,GoFiles"}, patterns...)
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64", "GOAMD64=v1", "CGO_ENABLED=0", "GOEXPERIMENT=", "GOFLAGS=")
	cmd.Stderr = os.Stderr
	list, err := cmd.Output()
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(list))
	var declarations []Declaration
	for {
		var pkg Package
		err := decoder.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		for _, filename := range pkg.GoFiles {
			positions := token.NewFileSet()
			file, err := parser.ParseFile(positions, filepath.Join(pkg.Dir, filename), nil, parser.SkipObjectResolution)
			if err != nil {
				return err
			}
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Name.Name == "_" {
					continue
				}
				line := positions.PositionFor(fn.Pos(), false).Line
				name := fn.Name.Name
				if fn.Recv != nil {
					name = receiver(fn.Recv.List[0].Type) + "." + name
				} else if name == "init" {
					name = fmt.Sprintf("init@%s:%d", filename, line)
				}
				declarations = append(declarations, Declaration{pkg.ImportPath, name, filename, line})
			}
		}
	}
	return json.NewEncoder(os.Stdout).Encode(declarations)
}

func main() {
	if err := scan(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
