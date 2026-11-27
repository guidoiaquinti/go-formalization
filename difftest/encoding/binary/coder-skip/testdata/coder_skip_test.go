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

type skipFixture struct {
	Kind, Size, Offset, Length int
	Backing                    []byte
	Want                       string
}

func (c skipFixture) valueAndSize() (reflect.Value, int) {
	switch c.Kind {
	case 0:
		return reflect.Value{}, -1
	case 1:
		return reflect.ValueOf(uint32(0)), 4
	case 2:
		return reflect.Zero(reflect.ArrayOf(c.Size, reflect.TypeFor[int16]())), 2 * c.Size
	case 3:
		return reflect.ValueOf(make([]byte, c.Size)), c.Size
	case 4:
		return reflect.ValueOf(struct {
			B byte
			W int64
		}{}), 9
	case 5:
		return reflect.ValueOf(struct{}{}), 0
	case 6:
		return reflect.Zero(reflect.ArrayOf(c.Size, reflect.TypeFor[string]())), -1
	default:
		return reflect.ValueOf("unsupported"), -1
	}
}

func TestGoLeanCoderSkip(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BINARY_CODER_SKIP_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		var c skipFixture
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		value, size := c.valueAndSize()
		backing := slices.Clone(c.Backing)
		read, write, status := binary.GoLeanCoderSkip(value, backing, c.Length, c.Offset)
		if strings.HasPrefix(status, "runtime error: slice bounds out of range") {
			status = "bounds"
		}
		got := skipWire(read, write, status, backing)
		wantBacking := slices.Clone(c.Backing)
		next := c.Offset + size
		wantWrite, wantStatus := next, "ok"
		if c.Offset < 0 || c.Offset > next || next > len(wantBacking) {
			wantWrite, wantStatus = c.Offset, "bounds"
		} else {
			for i := c.Offset; i < next; i++ {
				wantBacking[i] = 0
			}
		}
		if want := skipWire(next, wantWrite, wantStatus, wantBacking); got != c.Want || got != want {
			t.Fatalf("fixture %+v: native %q, oracle %q", c, got, want)
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("missing rows")
	}
	t.Logf("%d exact private skip calls per method matched Lean and oracle", rows)
}

func skipWire(read, write int, status string, backing []byte) string {
	parts := make([]string, len(backing))
	for i, b := range backing {
		parts[i] = strconv.Itoa(int(b))
	}
	return fmt.Sprintf("%d|%s:%d:%s", read, status, write, strings.Join(parts, ","))
}
