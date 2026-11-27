package binary_test

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type coderByteFixture struct {
	Nil            bool
	Length, Offset int
	Argument       byte
	Backing        []byte
	Want           string
}

func TestGoLeanCoderByte(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_BINARY_CODER_BYTE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for i := 0; ; i++ {
		var c coderByteFixture
		err := dec.Decode(&c)
		if err == io.EOF {
			if i == 0 {
				t.Fatal("missing cases")
			}
			t.Logf("%d rows matched across all six exact private methods", i)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		var rows []string
		for method := range 6 {
			backing := slices.Clone(c.Backing)
			result, offset := binary.GoLeanCoderByte(method, c.Length, c.Offset, backing, c.Argument, c.Nil)
			got := coderByteWire(result, offset, backing)
			switch {
			case strings.Contains(result, "nil pointer dereference"):
				got = "nil"
			case strings.HasPrefix(result, "runtime error: index out of range"):
				got = coderByteWire("bounds", offset, backing)
			}
			if want := coderByteOracle(c, method); got != want {
				t.Fatalf("fixture %+v method %d: native %q, oracle %q", c, method, got, want)
			}
			rows = append(rows, got)
		}
		if got := strings.Join(rows, "|"); got != c.Want {
			t.Fatalf("fixture %+v: native %q", c, got)
		}
	}
}

func coderByteOracle(c coderByteFixture, method int) string {
	if c.Nil {
		return "nil"
	}
	backing := slices.Clone(c.Backing)
	if c.Offset < 0 || c.Offset >= c.Length {
		return coderByteWire("bounds", c.Offset, backing)
	}
	result := "ok"
	switch method {
	case 0:
		result = strconv.FormatBool(backing[c.Offset] != 0)
	case 1:
		result = strconv.Itoa(int(backing[c.Offset]))
	case 2:
		b := int(backing[c.Offset])
		if b >= 128 {
			b -= 256
		}
		result = strconv.Itoa(b)
	case 3:
		backing[c.Offset] = 0
		if c.Argument != 0 {
			backing[c.Offset] = 1
		}
	default:
		backing[c.Offset] = c.Argument
	}
	return coderByteWire(result, c.Offset+1, backing)
}

func coderByteWire(result string, offset int, backing []byte) string {
	parts := make([]string, len(backing))
	for i, b := range backing {
		parts[i] = strconv.Itoa(int(b))
	}
	return fmt.Sprintf("%s:%d:%s", result, offset, strings.Join(parts, ","))
}
