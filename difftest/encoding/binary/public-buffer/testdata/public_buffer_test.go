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

func publicCall(call func() (int, error)) (status string, n int) {
	defer func() {
		if failure := recover(); failure != nil {
			status = codecStatus(fmt.Sprint(failure))
		}
	}()
	n, err := call()
	if err == nil {
		return "ok", n
	}
	if strings.Contains(err.Error(), "buffer too small") {
		return "short", n
	}
	if strings.Contains(err.Error(), "invalid type") || strings.Contains(err.Error(), "not fixed-sized") {
		return "type", n
	}
	return err.Error(), n
}

func publicByteWire(cells []byte) string {
	parts := make([]string, len(cells))
	for i, v := range cells {
		parts[i] = strconv.Itoa(int(v))
	}
	return strings.Join(parts, ",")
}

func publicWordWire(words []uint64) string {
	parts := make([]string, len(words))
	for i, v := range words {
		parts[i] = strconv.FormatUint(v, 10)
	}
	return strings.Join(parts, ",")
}

func publicValue(c codecFixture) (any, reflect.Value) {
	if c.Form == 3 {
		return "unsupported", reflect.Value{}
	}
	return c.materialize()
}

func TestGoLeanPublicBuffer(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BINARY_PUBLIC_BUFFER_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	count := 0
	for reader.Scan() {
		var c codecFixture
		if err := json.Unmarshal(reader.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		var order binary.ByteOrder = binary.LittleEndian
		if c.Order == 1 {
			order = binary.BigEndian
		}
		source, storage := publicValue(c)
		backing := slices.Clone(c.Backing)
		d, n := publicCall(func() (int, error) { return binary.Decode(backing[:c.Length], order, source) })
		if !slices.Equal(backing, c.Backing) {
			t.Fatal("Decode modified disjoint input")
		}
		read := fmt.Sprintf("%s:%d:%s", d, n, publicWordWire(observedWords(c, storage)))

		value, _ := publicValue(c)
		e, m := publicCall(func() (int, error) { return binary.Encode(backing[:c.Length], order, value) })
		write := fmt.Sprintf("%s:%d:%s", e, m, publicByteWire(backing))
		backing = slices.Clone(c.Backing)
		in := backing[:c.Length:len(backing)]
		out, err := binary.Append(in, order, value)
		status := "ok"
		if err != nil {
			if !strings.Contains(err.Error(), "not fixed-sized") {
				t.Fatalf("Append returned unexpected error %v", err)
			}
			status = "type"
		}
		newStore := len(out) > 0 && (len(backing) == 0 || &out[0] != &backing[0])
		appendWire := fmt.Sprintf("%s:%t:%t:%s:%s", status, newStore, out == nil,
			publicByteWire(out), publicByteWire(backing))
		got := read + "|" + write + "|" + appendWire
		if want := publicOracle(c); got != c.Want || got != want {
			t.Fatalf("fixture %+v: Go %q, oracle %q", c, got, want)
		}
		count++
	}
	if err := reader.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("missing fixtures")
	}
	t.Logf("%d public Go observations per API matched Lean and independent bit oracle", count)
}

func publicOracle(c codecFixture) string {
	namedByte := c.Kind == 11
	if namedByte {
		c.Kind = 2
	}
	n := len(c.Words) * [...]int{1, 1, 1, 2, 2, 4, 4, 8, 8, 4, 8}[c.Kind]
	raw := publicWordWire(c.Words)
	backing := publicByteWire(c.Backing)
	if c.Form == 3 {
		return fmt.Sprintf("type:0:%s|type:0:%s|type:false:true::%s", raw, backing, backing)
	}
	d, e := fmt.Sprintf("ok:%d:%s", n, raw), fmt.Sprintf("ok:%d:%s", n, backing)
	if c.Length < n {
		d, e = "short:0:"+raw, "short:0:"+backing
	} else if n > 0 {
		decoded, encoded, _ := strings.Cut(codecOracle(c), "|")
		readStatus, values, _ := strings.Cut(decoded, ":")
		switch readStatus {
		case "unsupported":
			d = "type:0:" + values
		case "bounds":
			d = "bounds:0:" + values
		default:
			d = fmt.Sprintf("ok:%d:%s", n, values)
		}
		_, written, _ := strings.Cut(encoded, ":")
		e = fmt.Sprintf("ok:%d:%s", n, written)
	}
	if namedByte && c.Form == 0 {
		d = "type:0:" + raw
	}
	output := make([]byte, c.Length+n)
	copy(output, c.Backing)
	payload := c
	payload.Length, payload.Backing = n, output[c.Length:]
	_, encoded, _ := strings.Cut(codecOracle(payload), "|")
	_, written, _ := strings.Cut(encoded, ":")
	outputWire := publicByteWire(output[:c.Length])
	if n != 0 {
		if outputWire != "" {
			outputWire += ","
		}
		outputWire += written
	}
	allocated := len(output) > len(c.Backing)
	if !allocated {
		tail := publicByteWire(c.Backing[len(output):])
		backing = strings.Join(slices.DeleteFunc([]string{outputWire, tail}, func(s string) bool { return s == "" }), ",")
	}
	return fmt.Sprintf("%s|%s|ok:%t:%t:%s:%s", d, e, allocated, c.Backing == nil && len(output) == 0, outputWire, backing)
}
