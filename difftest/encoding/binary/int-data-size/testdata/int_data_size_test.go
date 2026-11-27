package binary_test

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type namedByte uint8
type namedBytes []byte

type intSizeFixture struct {
	Form, Kind, Length int
	Bytes              []byte
	Want               string
}

var primitives = []any{true, int8(-1), byte(0xff), int16(-1), uint16(4), int32(-9), uint32(8),
	int64(-55), uint64(123), float32(1.5), float64(-2), int(3), uint(7), uintptr(0), "", complex64(1),
	complex128(2), namedByte(0x50), struct{ X uint32 }{}, new(*int8)}

func (c intSizeFixture) value() any {
	typ := reflect.TypeOf(primitives[c.Kind])
	switch c.Form {
	case 0:
		return primitives[c.Kind]
	case 1:
		return reflect.New(typ).Interface()
	case 2:
		return reflect.Zero(reflect.PointerTo(typ)).Interface()
	case 3:
		if c.Kind == 2 {
			return c.Bytes[:c.Length:len(c.Bytes)]
		}
		if c.Bytes == nil {
			return reflect.Zero(reflect.SliceOf(typ)).Interface()
		}
		return reflect.MakeSlice(reflect.SliceOf(typ), c.Length, len(c.Bytes)).Interface()
	case 4:
		return reflect.New(reflect.ArrayOf(c.Length, typ)).Elem().Interface()
	case 5:
		return namedBytes(c.Bytes)
	case 6:
		return new([]byte)
	case 7:
		return map[uint8]int16{}
	default:
		return nil
	}
}

func intSizeOracle(value any) (int, []byte) {
	if value == nil {
		return 0, nil
	}
	t := reflect.TypeOf(value)
	if data, ok := value.([]byte); ok {
		return len(data), data
	}
	n := 1
	switch t.Kind() {
	case reflect.Pointer:
		if t.Name() != "" {
			return 0, nil
		}
		t = t.Elem()
	case reflect.Slice:
		if t.Name() != "" {
			return 0, nil
		}
		n = reflect.ValueOf(value).Len()
		t = t.Elem()
	}
	if t.PkgPath() != "" {
		return 0, nil
	}
	switch t.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64, reflect.Float32, reflect.Float64:
		return n * int(t.Size()), nil
	}
	return 0, nil
}

func TestGoLeanIntDataSize(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BINARY_INT_DATA_SIZE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		var c intSizeFixture
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		value := c.value()
		size, buf := binary.GoLeanIntDataSize(value)
		wantSize, wantBuf := intSizeOracle(value)
		got := intSizeWire(size, buf)
		if got != c.Want || got != intSizeWire(wantSize, wantBuf) {
			t.Fatalf("form %d kind %d (%T): got %q, Lean %q, oracle %q", c.Form, c.Kind, value, got, c.Want, intSizeWire(wantSize, wantBuf))
		}
		if cap(buf) > 0 && &buf[:cap(buf)][0] != &wantBuf[:cap(wantBuf)][0] {
			t.Fatalf("%T: returned a copy, not the input slice", value)
		}
		if !slices.Equal(buf, wantBuf) {
			t.Fatal("returned view length or elements differ")
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("missing cases")
	}
	t.Logf("%d native rows matched Lean and reflection oracle", rows)
}

func intSizeWire(size int, b []byte) string {
	var fields []string
	for _, v := range b[:cap(b)] {
		fields = append(fields, strconv.Itoa(int(v)))
	}
	return fmt.Sprintf("%d:%t:%d:%s", size, b == nil, cap(b), strings.Join(fields, ","))
}
