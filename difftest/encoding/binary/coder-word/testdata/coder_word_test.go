package binary_test

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type wordFixture struct {
	Width, Order, Length, Offset int
	Word                         uint64
	Backing                      []byte
	Want                         string
}

type tracingOrder struct {
	binary.ByteOrder
	t       *testing.T
	fixture wordFixture
	backing []byte
	method  int
	calls   int
	trace   string
}

func (o *tracingOrder) access(b []byte, width int, input uint64, writing bool) uint64 {
	o.calls++
	c := o.fixture
	if o.calls != 1 || width != c.Width || len(b) != width || cap(b) != len(o.backing)-c.Offset ||
		&b[0] != &o.backing[c.Offset] || writing != (o.method >= 2) {
		o.t.Fatal("incorrect callback shape or alias")
	}
	if writing && input != c.Word&wordMask(width) {
		o.t.Fatalf("incorrect encoder argument: got %x want %x", input, c.Word&wordMask(width))
	}
	o.trace = strconv.Itoa(cap(b))
	b[0] = 165
	b[:cap(b)][cap(b)-1] = 90
	if c.Order == 3 {
		panic("scripted ByteOrder panic")
	}
	return wordMask(width)
}

func (o *tracingOrder) Uint16(b []byte) uint16       { return uint16(o.access(b, 2, 0, false)) }
func (o *tracingOrder) Uint32(b []byte) uint32       { return uint32(o.access(b, 4, 0, false)) }
func (o *tracingOrder) Uint64(b []byte) uint64       { return o.access(b, 8, 0, false) }
func (o *tracingOrder) PutUint16(b []byte, x uint16) { o.access(b, 2, uint64(x), true) }
func (o *tracingOrder) PutUint32(b []byte, x uint32) { o.access(b, 4, uint64(x), true) }
func (o *tracingOrder) PutUint64(b []byte, x uint64) { o.access(b, 8, x, true) }

func TestGoLeanCoderWord(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BINARY_CODER_WORD_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	counts := [9]int{}
	for reader.Scan() {
		var c wordFixture
		if err := json.Unmarshal(reader.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		rows := make([]string, 0, 4)
		for method := range 4 {
			got, expected := nativeWord(t, c, method), oracleWord(c, method)
			if got != expected {
				t.Fatalf("fixture %+v method %d: Go %q, oracle %q", c, method, got, expected)
			}
			rows = append(rows, got)
		}
		if got := strings.Join(rows, "|"); got != c.Want {
			t.Fatalf("fixture %+v: native %q", c, got)
		}
		counts[c.Width]++
	}
	if err := reader.Err(); err != nil {
		t.Fatal(err)
	}
	if counts[2] == 0 || counts[4] == 0 || counts[8] == 0 {
		t.Fatal("missing widths")
	}
	t.Logf("matched Go, Lean and independent oracle on %d/%d/%d rows (16/32/64-bit)", counts[2], counts[4], counts[8])
}

func nativeWord(t *testing.T, c wordFixture, method int) string {
	backing := slices.Clone(c.Backing)
	trace := &tracingOrder{t: t, fixture: c, backing: backing, method: method, trace: "-"}
	var order binary.ByteOrder
	switch c.Order {
	case 0:
		order = binary.LittleEndian
	case 1:
		order = binary.BigEndian
	case 2, 3:
		order = trace
	}
	value, status, offset := binary.GoLeanCoderWord(order, c.Width, method, c.Length, c.Offset, backing, c.Word)
	switch {
	case status == "ok":
	case status == "scripted ByteOrder panic":
		status = "callback"
	case strings.HasPrefix(status, "runtime error: slice bounds out of range"):
		status = "bounds"
	case strings.Contains(status, "nil pointer dereference"):
		status = "nil"
	}
	return wordWire(status, value, offset, trace.trace, backing)
}

func oracleWord(c wordFixture, method int) string {
	cells := slices.Clone(c.Backing)
	if c.Offset < 0 || c.Offset > len(cells)-c.Width {
		return wordWire("bounds", "0", c.Offset, "-", cells)
	}
	if c.Order == 4 {
		return wordWire("nil", "0", c.Offset, "-", cells)
	}
	var value uint64
	trace := "-"
	if c.Order == 2 || c.Order == 3 {
		cells[c.Offset], cells[len(cells)-1] = 165, 90
		trace = strconv.Itoa(len(cells) - c.Offset)
		if c.Order == 3 {
			return wordWire("callback", "0", c.Offset, trace, cells)
		}
		value = wordMask(c.Width)
	} else {
		for i := range c.Width {
			k := i
			if c.Order == 1 {
				k = c.Width - i - 1
			}
			index := c.Offset + i
			if method >= 2 {
				cells[index] = byte(c.Word >> (8 * k))
			} else {
				value += uint64(cells[index]) << (8 * k)
			}
		}
	}
	number := "0"
	if method == 0 {
		number = strconv.FormatUint(value, 10)
	} else if method == 1 {
		shift := 64 - 8*c.Width
		number = strconv.FormatInt(int64(value<<shift)>>shift, 10)
	}
	return wordWire("ok", number, c.Offset+c.Width, trace, cells)
}

func wordMask(width int) uint64 { return ^uint64(0) >> (64 - 8*width) }

func wordWire(status, value string, offset int, trace string, backing []byte) string {
	parts := make([]string, len(backing))
	for i, b := range backing {
		parts[i] = strconv.Itoa(int(b))
	}
	return fmt.Sprintf("%s:%s:%d:%s:%s", status, value, offset, trace, strings.Join(parts, ","))
}
