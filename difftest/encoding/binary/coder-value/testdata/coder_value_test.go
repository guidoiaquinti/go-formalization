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

type valueFixture struct {
	Kind, Shape, Count, Order, Length, Offset int
	Words                                     []uint64
	Backing                                   []byte
	Want                                      string
}

var valueAtoms = []any{false, int8(0), uint8(0), int16(0), uint16(0), int32(0), uint32(0),
	int64(0), uint64(0), float32(0), float64(0), complex64(0), complex128(0), ""}

func (c valueFixture) materialize() (reflect.Value, []reflect.Value) {
	typ := reflect.TypeOf(valueAtoms[c.Kind])
	var root reflect.Value
	var elements []reflect.Value
	switch c.Shape {
	case 0:
		root = reflect.New(typ).Elem()
		elements = append(elements, root)
	case 1, 2:
		if c.Shape == 1 {
			root = reflect.New(reflect.ArrayOf(c.Count, typ)).Elem()
		} else {
			root = reflect.MakeSlice(reflect.SliceOf(typ), c.Count, c.Count)
		}
		for i := range c.Count {
			elements = append(elements, root.Index(i))
		}
	default:
		inner := reflect.ArrayOf(c.Count, typ)
		root = reflect.New(reflect.StructOf([]reflect.StructField{
			{Name: "First", Type: typ},
			{Name: "_", PkgPath: "encoding/binary_test", Type: reflect.TypeFor[uint16]()},
			{Name: "Rest", Type: inner},
		})).Elem()
		elements = append(elements, root.Field(0))
		for i := range c.Count {
			elements = append(elements, root.Field(2).Index(i))
		}
	}
	words := c.Words
	for _, v := range elements {
		switch v.Kind() {
		case reflect.Bool:
			v.SetBool(words[0] != 0)
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(words[0]))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(words[0])
		case reflect.Float32:
			v.SetFloat(float64(math.Float32frombits(uint32(words[0]))))
		case reflect.Float64:
			v.SetFloat(math.Float64frombits(words[0]))
		case reflect.Complex64:
			v.SetComplex(complex(float64(math.Float32frombits(uint32(words[0]))), float64(math.Float32frombits(uint32(words[1])))))
			words = words[1:]
		case reflect.Complex128:
			v.SetComplex(complex(math.Float64frombits(words[0]), math.Float64frombits(words[1])))
			words = words[1:]
		default:
			continue
		}
		words = words[1:]
	}
	return root, elements
}

func valueWords(elements []reflect.Value) []uint64 {
	var words []uint64
	for _, v := range elements {
		switch v.Kind() {
		case reflect.Bool:
			word := uint64(0)
			if v.Bool() {
				word = 1
			}
			words = append(words, word)
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			words = append(words, uint64(v.Int())&(^uint64(0)>>(64-8*v.Type().Size())))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			words = append(words, v.Uint())
		case reflect.Float32:
			words = append(words, uint64(math.Float32bits(float32(v.Float()))))
		case reflect.Float64:
			words = append(words, math.Float64bits(v.Float()))
		case reflect.Complex64:
			x := v.Complex()
			words = append(words, uint64(math.Float32bits(float32(real(x)))), uint64(math.Float32bits(float32(imag(x)))))
		case reflect.Complex128:
			x := v.Complex()
			words = append(words, math.Float64bits(real(x)), math.Float64bits(imag(x)))
		}
	}
	return words
}

func TestGoLeanCoderValue(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BINARY_CODER_VALUE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		var c valueFixture
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		var order binary.ByteOrder = binary.LittleEndian
		if c.Order != 0 {
			order = binary.BigEndian
		}
		backing := slices.Clone(c.Backing)
		value, readElements := c.materialize()
		d, offset := binary.GoLeanCoderValue(false, order, backing[:c.Length:len(backing)], c.Offset, value)
		if !slices.Equal(backing, c.Backing) {
			t.Fatal("decoder mutated input backing")
		}
		source, sourceElements := c.materialize()
		e, next := binary.GoLeanCoderValue(true, order, backing[:c.Length:len(backing)], c.Offset, source)
		if !slices.Equal(valueWords(sourceElements), c.Words) {
			t.Fatal("encoder modified the source value")
		}
		got := valueWire(d, offset, valueWords(readElements), e, next, backing)
		if want := valueOracle(c); got != want || got != c.Want {
			t.Fatalf("fixture %+v: Go %q, oracle %q", c, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("missing fixtures")
	}
	t.Logf("%d native observations per method matched Lean and independent field oracle", count)
}

func valueOracle(c valueFixture) string {
	width := [...]int{1, 1, 1, 2, 2, 4, 4, 8, 8, 4, 8, 4, 8, 0}[c.Kind]
	parts := 1
	if c.Kind == 11 || c.Kind == 12 {
		parts = 2
	}
	read, write := c.Offset, c.Offset
	d, e := true, true
	values, backing := slices.Clone(c.Words), slices.Clone(c.Backing)
	limit := len(backing)
	if width == 1 {
		limit = c.Length
	}
	for i := 0; i <= len(values); i += parts {
		if c.Shape == 3 && (i == parts || (i == 0 && width == 0)) {
			if d {
				read += 2
			}
			if e && write+2 <= len(backing) {
				clear(backing[write : write+2])
				write += 2
			} else {
				e = false
			}
		}
		if i == len(values) {
			break
		}
		var decoded [2]uint64
		for part := range parts {
			d = d && read+width <= limit
			e = e && write+width <= limit
			for b := range width {
				shift := b
				if c.Order != 0 {
					shift = width - b - 1
				}
				if d {
					decoded[part] += uint64(c.Backing[read+b]) << (8 * shift)
				}
				if e {
					backing[write+b] = byte(c.Words[i+part] >> (8 * shift))
				}
			}
			if d {
				read += width
			}
			if e {
				write += width
			}
		}
		if d {
			if c.Kind == 0 && decoded[0] != 0 {
				decoded[0] = 1
			}
			copy(values[i:], decoded[:parts])
		}
	}
	return valueWire(d, read, values, e, write, backing)
}

func valueWire(d bool, offset int, words []uint64, e bool, next int, backing []byte) string {
	a, b := make([]string, len(words)), make([]string, len(backing))
	for i, word := range words {
		a[i] = strconv.FormatUint(word, 10)
	}
	for i, cell := range backing {
		b[i] = strconv.Itoa(int(cell))
	}
	return fmt.Sprintf("%t:%d:%s|%t:%d:%s", d, offset, strings.Join(a, ","), e, next, strings.Join(b, ","))
}
