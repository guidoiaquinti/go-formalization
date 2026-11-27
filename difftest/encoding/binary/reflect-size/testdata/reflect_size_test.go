package binary_test

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type namedWord uint32

var reflectSizeAtoms = map[int]any{
	1: "", 2: true, 3: int16(0), 4: namedWord(0), 5: int64(0), 6: complex128(0),
	12: int8(0), 13: uint8(0), 14: uint16(0), 15: int32(0), 16: uint32(0), 17: uint64(0),
	18: float32(0), 19: float64(0), 20: complex64(0), 21: int(0), 22: map[int]int{},
}

func parseValue(tokens []int) (reflect.Value, []int) {
	tag, tail := tokens[0], tokens[1:]
	if atom, ok := reflectSizeAtoms[tag]; ok {
		return reflect.ValueOf(atom), tail
	}
	switch tag {
	case 0:
		return reflect.Value{}, tail
	case 7, 8:
		var fields []reflect.StructField
		for tag == 8 {
			field, rest := parseValue(tail)
			fields = append(fields, reflect.StructField{Name: fmt.Sprintf("F%d", len(fields)), Type: field.Type()})
			tag, tail = rest[0], rest[1:]
		}
		return reflect.Zero(reflect.StructOf(fields)), tail
	case 9, 10, 11:
		n := tail[0]
		child, rest := parseValue(tail[1:])
		switch tag {
		case 9:
			return reflect.Zero(reflect.ArrayOf(n, child.Type())), rest
		case 10:
			return reflect.MakeSlice(reflect.SliceOf(child.Type()), n, n), rest
		default:
			if n != 0 {
				return reflect.Zero(reflect.PointerTo(child.Type())), rest
			}
			value := reflect.New(child.Type())
			value.Elem().Set(child)
			return value, rest
		}
	}
	panic("invalid layout tag")
}

func TestGoLeanReflectSize(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BINARY_REFLECT_SIZE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		var c struct {
			Tokens []int
			Want   string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		v, tail := parseValue(c.Tokens)
		if len(tail) != 0 {
			t.Fatalf("trailing layout tokens %v", tail)
		}
		var x any
		var typ reflect.Type
		if v.IsValid() {
			x, typ = v.Interface(), v.Type()
		}
		byType := nativeTypeSize(typ)
		cold, warm := binary.GoLeanReflectValueSize(v)
		public := binary.Size(x)
		row := fmt.Sprintf("%s:%d:%d", byType, cold, public)
		if row != c.Want || warm != cold || cold != valueSizeOracle(v) || public != valueSizeOracle(reflect.Indirect(v)) {
			t.Fatalf("tokens %v: native %q warm %d, Lean %q; oracle %d/%d", c.Tokens, row, warm, c.Want, valueSizeOracle(v), valueSizeOracle(reflect.Indirect(v)))
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("missing cases")
	}
	t.Logf("%d Go/Lean layouts matched sizeof, cold/warm dataSize, and Size", rows)
}

func nativeTypeSize(typ reflect.Type) (label string) {
	defer func() {
		if failure := recover(); failure != nil {
			if typ != nil || !strings.Contains(fmt.Sprint(failure), "nil pointer dereference") {
				panic(failure)
			}
			label = "panic"
		}
	}()
	return strconv.Itoa(binary.GoLeanReflectTypeSize(typ))
}

func valueSizeOracle(v reflect.Value) int {
	if !v.IsValid() {
		return -1
	}
	if v.Kind() != reflect.Slice {
		return typeFootprint(v.Type())
	}
	size := typeFootprint(v.Type().Elem())
	if size < 0 {
		return -1
	}
	return v.Len() * size
}

func typeFootprint(t reflect.Type) int {
	if t.Kind() == reflect.Array {
		width := typeFootprint(t.Elem())
		if width < 0 {
			return -1
		}
		return width * t.Len()
	}
	if t.Kind() == reflect.Struct {
		var width int
		for i := range t.NumField() {
			field := typeFootprint(t.Field(i).Type)
			if field < 0 {
				return -1
			}
			width += field
		}
		return width
	}
	switch t.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16, reflect.Int32,
		reflect.Uint32, reflect.Int64, reflect.Uint64, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return int(t.Size())
	default:
		return -1
	}
}
