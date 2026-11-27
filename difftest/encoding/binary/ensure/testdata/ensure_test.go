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

type goLeanEnsureFixture struct {
	Backing       []byte
	Length, Count int
	Want          string
}

func TestGoLeanEnsure(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_BINARY_ENSURE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for i := 0; ; i++ {
		var c goLeanEnsureFixture
		err := dec.Decode(&c)
		if err == io.EOF {
			if i == 0 {
				t.Fatal("missing cases")
			}
			t.Logf("%d native cases matched Lean and the copy oracle", i)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		got, expected := goLeanEnsureNative(t, c), goLeanEnsureOracle(c)
		if got != expected || got != c.Want {
			t.Fatalf("%+v: native %q, oracle %q", c, got, expected)
		}
	}
}

func goLeanEnsureNative(t *testing.T, c goLeanEnsureFixture) (result string) {
	backing := slices.Clone(c.Backing)
	defer func() {
		if failure := recover(); failure != nil {
			if !slices.Equal(backing, c.Backing) {
				t.Fatal("panic mutated the original storage")
			}
			result = fmt.Sprintf("panic:%v", failure)
		}
	}()
	out, pos := binary.GoLeanEnsure(backing[:c.Length:len(backing)], c.Count)
	if !slices.Equal(backing, c.Backing) {
		t.Fatal("ensure mutated the original storage")
	}
	if len(pos) != c.Count || len(out) != c.Length+c.Count ||
		(len(pos) > 0 && &pos[0] != &out[c.Length]) {
		t.Fatal("ensure returned incorrect view geometry")
	}
	allocated := cap(out) > 0 && (len(backing) == 0 || &out[:cap(out)][0] != &backing[0])
	return goLeanEnsureWire(out == nil, allocated, out, pos)
}

func goLeanEnsureOracle(c goLeanEnsureFixture) string {
	if c.Count < 0 {
		return "panic:cannot be negative"
	}
	width := c.Length + c.Count
	out := make([]byte, width)
	copy(out, c.Backing)
	return goLeanEnsureWire(c.Backing == nil && width == 0, width > len(c.Backing), out, out[c.Length:])
}

func goLeanEnsureWire(isNil, allocated bool, buffer, pos []byte) string {
	render := func(b []byte) string {
		parts := make([]string, len(b))
		for i, x := range b {
			parts[i] = strconv.Itoa(int(x))
		}
		return strings.Join(parts, ",")
	}
	return fmt.Sprintf("%t:%t:%s:%s", isNil, allocated, render(buffer), render(pos))
}
