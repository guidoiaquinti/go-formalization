package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Guard the Go declaration transcribed by GoLean.Math.Bits.mul64Source.
func TestPinnedMul64Source(t *testing.T) {
	pin, err := os.ReadFile("../../.go-version")
	if err != nil {
		t.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		t.Fatalf("source guard needs %s, got %s", want, runtime.Version())
	}
	positions := token.NewFileSet()
	path := filepath.Join(runtime.GOROOT(), "src", "math", "bits", "bits.go")
	file, err := parser.ParseFile(positions, path, nil, parser.SkipObjectResolution)
	if err != nil {
		t.Fatal(err)
	}

	const transcribed = `func Mul64(x, y uint64) (hi, lo uint64) {
	const mask32 = 1<<32 - 1
	x0 := x & mask32
	x1 := x >> 32
	y0 := y & mask32
	y1 := y >> 32
	w0 := x0 * y0
	t := x1*y0 + w0>>32
	w1 := t & mask32
	w2 := t >> 32
	w1 += x0 * y1
	hi = x1*y1 + w2 + w1>>32
	lo = x * y
	return
}`
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name.Name != "Mul64" {
			continue
		}
		var formatted bytes.Buffer
		if err := format.Node(&formatted, positions, fn); err != nil {
			t.Fatal(err)
		}
		if got := formatted.String(); got != transcribed {
			t.Fatalf("pinned Mul64 declaration changed; recheck the Lean transcription\ngot:\n%s\nwant:\n%s", got, transcribed)
		}
		return
	}
	t.Fatal("pinned math/bits source has no Mul64 declaration")
}
