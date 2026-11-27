package json

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type compatFixture struct {
	Present                             bool
	Type, Nested                        int
	Value, Struct, Field, Source, Inner []byte
	Want                                string
}

var marker = new(int)
var numericSuffix = regexp.MustCompile(`(^|\.)[0-9]+$`)

type compatInner struct {
	value string
	calls *int
}

func (e *compatInner) Error() string {
	if e == nil {
		panic(marker)
	}
	*e.calls++
	return e.value
}

func TestGoLeanCompatErrors(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_JSON_COMPAT_ERRORS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	decoder := NewDecoder(f)
	for row := 0; ; row++ {
		var c compatFixture
		if err := decoder.Decode(&c); err != nil {
			if err == io.EOF && row > 0 {
				t.Logf("%d pinned native rows per declaration", row)
				return
			}
			t.Fatal(err)
		}
		native, oracle := checkCompatFixture(c)
		if native != c.Want || native != oracle {
			t.Fatalf("row %d: Go %q Lean %q independent %q", row, native, c.Want, oracle)
		}
	}
}

func checkCompatFixture(c compatFixture) (string, string) {
	var s *SyntaxError
	var v *UnsupportedValueError
	var typErr *UnsupportedTypeError
	var invalid *InvalidUnmarshalError
	var marshaler *MarshalerError
	var unmarshal *UnmarshalTypeError
	checkCalls := 0
	typ := []reflect.Type{nil, reflect.TypeFor[int](), reflect.TypeFor[*int](), reflect.TypeFor[[]string]()}[c.Type]
	var inner error
	switch c.Nested {
	case 1:
		inner = &compatInner{string(c.Inner), &checkCalls}
	case 2:
		inner = (*compatInner)(nil)
	}
	if c.Present {
		s = &SyntaxError{msg: string(c.Value), Offset: 17}
		v = &UnsupportedValueError{Value: reflect.ValueOf(117), Str: string(c.Value)}
		typErr = &UnsupportedTypeError{Type: typ}
		invalid = &InvalidUnmarshalError{Type: typ}
		marshaler = &MarshalerError{Type: typ, Err: inner, sourceFunc: string(c.Source)}
		unmarshal = &UnmarshalTypeError{Type: typ, Err: inner, Value: string(c.Value), Struct: string(c.Struct), Field: string(c.Field), Offset: 27}
	}
	wrapped := func(f func() error) string {
		return capture(func() string {
			if f() != inner {
				return "wrong underlying error"
			}
			return strconv.Itoa(c.Nested)
		})
	}
	fields := []string{captureText(s.Error), captureText(v.Error), captureText(typErr.Error), captureText(invalid.Error),
		captureText(marshaler.Error), wrapped(marshaler.Unwrap), captureText(unmarshal.Error), wrapped(unmarshal.Unwrap)}
	wantCalls := 0
	if c.Present && typ != nil && c.Nested == 1 {
		wantCalls = 2
	}
	if checkCalls != wantCalls {
		panic("unexpected nested Error call count")
	}
	return strings.Join(fields, "|"), compatOracle(c)
}

func capture(call func() string) (result string) {
	defer func() {
		switch p := recover(); p {
		case nil:
		case marker:
			result = "nestedpanic"
		default:
			if _, ok := p.(runtime.Error); !ok {
				panic(p)
			}
			result = "nilpanic"
		}
	}()
	return call()
}

func captureText(call func() string) string {
	return capture(func() string { return wire(call()) })
}

func wire(s string) string {
	parts := make([]string, len(s))
	for i, b := range []byte(s) {
		parts[i] = strconv.Itoa(int(b))
	}
	return strings.Join(parts, ",")
}

func compatOracle(c compatFixture) string {
	if !c.Present {
		return strings.TrimSuffix(strings.Repeat("nilpanic|", 8), "|")
	}
	value := string(c.Value)
	underlying := strconv.Itoa(c.Nested)
	out := []string{wire(value), wire("json: unsupported value: " + value), "nilpanic",
		wire("json: Unmarshal(nil)"), "nilpanic", underlying, "nilpanic", underlying}
	if c.Type == 0 {
		return strings.Join(out, "|")
	}
	name := []string{"", "int", "*int", "[]string"}[c.Type]
	out[2] = wire(fmt.Sprintf("json: unsupported type: %s", name))
	if c.Type == 2 {
		out[3] = wire(fmt.Sprintf("json: Unmarshal(nil %s)", name))
	} else {
		out[3] = wire(fmt.Sprintf("json: Unmarshal(non-pointer %s)", name))
	}
	if c.Nested == 2 {
		out[4], out[6] = "nestedpanic", "nestedpanic"
		return strings.Join(out, "|")
	}
	target := "Go value"
	if len(c.Struct)+len(c.Field) > 0 {
		target = fmt.Sprintf("%s.%s", c.Struct, c.Field)
		if !numericSuffix.Match(c.Field) {
			target = "Go struct field " + target
		}
	}
	base := fmt.Sprintf("json: cannot unmarshal %s into %s of type %s", value, target, name)
	if c.Nested == 1 {
		source := string(c.Source)
		if source == "" {
			source = "MarshalJSON"
		}
		out[4] = wire(fmt.Sprintf("json: error calling %s for type %s: %s", source, name, c.Inner))
		base = fmt.Sprintf("%s: %s", base, c.Inner)
	}
	out[6] = wire(base)
	return strings.Join(out, "|")
}
