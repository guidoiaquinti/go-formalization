package constant

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func goLeanConstantWire(x Value) string {
	switch x.(type) {
	case int64Val:
		return "small:" + x.ExactString()
	case intVal:
		return "large:" + x.ExactString()
	default:
		panic("unexpected representation")
	}
}
func goLeanConstantBytes(raw []byte) string {
	fields := make([]string, len(raw))
	for i, v := range raw {
		fields[i] = strconv.Itoa(int(v))
	}
	return strings.Join(fields, ",")
}
func goLeanConstantCatch(fn func() string) (out string) {
	defer func() {
		if v := recover(); v != nil {
			if v != "true not an Int" && v != "unknown not an Int" {
				panic(v)
			}
			out = "panic"
		}
	}()
	return fn()
}
func goLeanFitsSigned64(x *big.Int) bool {
	return x.Cmp(big.NewInt(math.MinInt64)) >= 0 && x.Cmp(big.NewInt(math.MaxInt64)) <= 0
}

func goLeanConstantOracle(x *big.Int) string {
	kind := "large"
	if goLeanFitsSigned64(x) {
		kind = "small"
	}
	return kind + ":" + x.String()
}

func TestGoLeanConstantInteger(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_CONSTANT_INTEGER_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	lines := bufio.NewScanner(f)
	count := 0
	for lines.Scan() {
		var c struct {
			Tag     int
			Integer string
			Bytes   []byte
			Want    string
		}
		if err := json.Unmarshal(lines.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		x, ok := new(big.Int).SetString(c.Integer, 10)
		if !ok {
			t.Fatal("invalid integer fixture")
		}
		low := uint64(x.Int64())
		unsigned := new(big.Int).SetUint64(low)
		fromRaw := slices.Clone(c.Bytes)
		copied := MakeFromBytes(fromRaw)
		for j := range fromRaw {
			fromRaw[j] ^= 255
		}
		reversed := slices.Clone(c.Bytes)
		slices.Reverse(reversed)
		magnitude := new(big.Int).SetBytes(reversed)
		canonical := makeInt(x)
		if large, ok := canonical.(intVal); ok && large.val != x {
			t.Fatal("makeInt detached the large integer")
		}
		var selected Value
		switch c.Tag {
		case 0:
			selected = canonical
		case 1:
			selected = copied
		case 2:
			selected = MakeUnknown()
		default:
			selected = MakeBool(true)
		}
		native := []string{goLeanConstantWire(MakeInt64(int64(low))), goLeanConstantWire(MakeUint64(low)),
			goLeanConstantWire(canonical), goLeanConstantWire(copied),
			goLeanConstantCatch(func() string { v, ok := Int64Val(selected); return fmt.Sprintf("%d:%t", v, ok) }),
			goLeanConstantCatch(func() string { v, ok := Uint64Val(selected); return fmt.Sprintf("%d:%t", v, ok) }),
			goLeanConstantCatch(func() string { return strconv.Itoa(BitLen(selected)) }),
			goLeanConstantCatch(func() string { return goLeanConstantBytes(Bytes(selected)) }),
		}
		expected := []string{fmt.Sprintf("small:%d", int64(low)), goLeanConstantOracle(unsigned), goLeanConstantOracle(x), goLeanConstantOracle(magnitude)}
		switch c.Tag {
		case 0, 1:
			n := x
			if c.Tag == 1 {
				n = magnitude
			}
			abs := new(big.Int).Abs(n)
			signedLow := int64(abs.Uint64())
			if n.Sign() < 0 {
				signedLow = -signedLow
			}
			lowUnsigned := abs.Uint64()
			if goLeanFitsSigned64(n) {
				lowUnsigned = uint64(signedLow)
			}
			raw := abs.Bytes()
			slices.Reverse(raw)
			expected = append(expected, fmt.Sprintf("%d:%t", signedLow, goLeanFitsSigned64(n)), fmt.Sprintf("%d:%t", lowUnsigned, n.Sign() >= 0 && n.BitLen() <= 64),
				strconv.Itoa(abs.BitLen()), goLeanConstantBytes(raw))
		case 2:
			expected = append(expected, "0:false", "0:false", "0", "panic")
		default:
			expected = append(expected, "panic", "panic", "panic", "panic")
		}
		if strings.Join(native, "|") != c.Want || !slices.Equal(native, expected) {
			t.Fatalf("%+v: native %v, oracle %v", c, native, expected)
		}
		count++
	}
	if err := lines.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("constant integer: %d native observations per declaration", count)
}
