package types

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanBasicFields(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BASIC_TYPE_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	wire := func(name string) string {
		words := make([]string, len(name))
		for i, b := range []byte(name) {
			words[i] = strconv.Itoa(int(b))
		}
		return strings.Join(words, ",")
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
			Kind, Info int
			Present    bool
			Name       []byte
			Want       string
		}
		err := dec.Decode(&c)
		if err == io.EOF && row > 0 {
			t.Logf("%d pinned native rows", row)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		var b *Basic
		if c.Present {
			b = &Basic{kind: BasicKind(c.Kind), info: BasicInfo(c.Info), name: string(c.Name)}
		}
		underlying, isBasic := b.Underlying().(*Basic)
		if !isBasic || underlying != b {
			t.Fatalf("row %d: wrong dynamic type or pointer from Underlying", row)
		}
		identity := "same"
		if underlying == nil {
			identity = "typedNil"
		}
		got := strings.Join([]string{
			read(func() string { return fmt.Sprint(b.Kind()) }), read(func() string { return fmt.Sprint(b.Info()) }),
			read(func() string { return wire(b.Name()) }), identity, read(func() string { return wire(b.String()) }),
		}, "|")
		if got != c.Want {
			t.Fatalf("row %d: Go %q want %q", row, got, c.Want)
		}
	}
}
