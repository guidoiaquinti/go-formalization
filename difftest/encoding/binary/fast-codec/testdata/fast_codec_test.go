package binary_test

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type codecFixture struct {
	Kind, Form, Order, Length int
	Words                     []uint64
	Backing                   []byte
	Want                      string
}

type definedByte uint8

var codecAtoms = []any{false, int8(0), uint8(0), int16(0), uint16(0), int32(0),
	uint32(0), int64(0), uint64(0), float32(0), float64(0), definedByte(0)}

func (c codecFixture) materialize() (any, reflect.Value) {
	if c.Form == 3 {
		return [2]byte{}, reflect.Value{}
	}
	typ := reflect.TypeOf(codecAtoms[c.Kind])
	storage := reflect.MakeSlice(reflect.SliceOf(typ), len(c.Words), len(c.Words))
	for i, word := range c.Words {
		elem := storage.Index(i)
		switch elem.Kind() {
		case reflect.Bool:
			elem.SetBool(word != 0)
		case reflect.Float32:
			elem.Set(reflect.ValueOf(math.Float32frombits(uint32(word))))
		case reflect.Float64:
			elem.Set(reflect.ValueOf(math.Float64frombits(word)))
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			elem.SetInt(int64(word))
		default:
			elem.SetUint(word)
		}
	}
	switch c.Form {
	case 0:
		return storage.Index(0).Interface(), storage
	case 1:
		return storage.Index(0).Addr().Interface(), storage
	default:
		return storage.Interface(), storage
	}
}

func observedWords(c codecFixture, storage reflect.Value) []uint64 {
	if !storage.IsValid() {
		return c.Words
	}
	mask := ^uint64(0) >> (64 - 8*storage.Type().Elem().Size())
	words := make([]uint64, storage.Len())
	for i := range words {
		elem := storage.Index(i)
		switch value := elem.Interface().(type) {
		case bool:
			if value {
				words[i] = 1
			}
		case float32:
			words[i] = uint64(math.Float32bits(value))
		case float64:
			words[i] = math.Float64bits(value)
		default:
			switch elem.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				words[i] = uint64(elem.Int()) & mask
			default:
				words[i] = elem.Uint()
			}
		}
	}
	return words
}

func TestGoLeanFastCodec(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BINARY_FAST_CODEC_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := 0
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		var c codecFixture
		if err := json.Unmarshal(reader.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		native, expected := codecNative(c), codecOracle(c)
		if native != c.Want || native != expected {
			t.Fatalf("fixture %+v: native %q, oracle %q", c, native, expected)
		}
		rows++
	}
	if err := reader.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("missing fixtures")
	}
	t.Logf("%d exact private observations per method matched Lean and bit-position oracle", rows)
}

func codecNative(c codecFixture) string {
	var order binary.ByteOrder = binary.LittleEndian
	if c.Order == 1 {
		order = binary.BigEndian
	}
	target, storage := c.materialize()
	input := slices.Clone(c.Backing)
	d := codecStatus(binary.GoLeanFastDecode(input[:c.Length], order, target))
	if !slices.Equal(input, c.Backing) {
		panic("fast decode mutated disjoint input")
	}
	source, _ := c.materialize()
	out := slices.Clone(c.Backing)
	e := codecStatus(binary.GoLeanFastEncode(out[:c.Length], order, source))
	return codecWire(d, observedWords(c, storage), e, out)
}

func codecStatus(status string) string {
	if strings.HasPrefix(status, "runtime error: index out of range") ||
		strings.HasPrefix(status, "runtime error: slice bounds out of range") {
		return "bounds"
	}
	return status
}

func codecOracle(c codecFixture) string {
	backing := slices.Clone(c.Backing)
	values := slices.Clone(c.Words)
	d, e := "ok", "ok"
	if c.Kind == 11 || c.Form == 3 {
		return codecWire("unsupported", values, e, backing)
	}
	width := int(reflect.TypeOf(codecAtoms[c.Kind]).Size())
	count := len(values)
	switch {
	case c.Form == 0:
		d = "unsupported"
		count = 0
	case c.Form == 2 && c.Kind <= 2:
		count = min(count, c.Length)
		if c.Kind != 2 && c.Length > len(values) {
			d = "bounds"
		}
	}
	byteShift := func(i int) int {
		if c.Order == 1 {
			i = width - 1 - i
		}
		return 8 * i
	}
	for i := range count {
		end := (i + 1) * width
		if end > c.Length {
			d = "bounds"
			break
		}
		var value uint64
		for j, b := range c.Backing[end-width : end] {
			value |= uint64(b) << byteShift(j)
		}
		if c.Kind == 0 && value != 0 {
			value = 1
		}
		values[i] = value
	}
	for i, word := range c.Words {
		end := (i + 1) * width
		if end > c.Length {
			if c.Form != 2 || c.Kind != 2 {
				e = "bounds"
			}
			break
		}
		for j := range width {
			backing[end-width+j] = byte(word >> byteShift(j))
		}
	}
	return codecWire(d, values, e, backing)
}

func codecWire(d string, words []uint64, e string, backing []byte) string {
	raw, cells := make([]string, len(words)), make([]string, len(backing))
	for i, w := range words {
		raw[i] = strconv.FormatUint(w, 10)
	}
	for i, w := range backing {
		cells[i] = strconv.Itoa(int(w))
	}
	return fmt.Sprintf("%s:%s|%s:%s", d, strings.Join(raw, ","), e, strings.Join(cells, ","))
}
