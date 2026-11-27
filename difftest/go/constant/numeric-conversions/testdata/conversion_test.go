package constant

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanNumericConversions(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_CONSTANT_CONVERSION_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := json.NewDecoder(file)
	float := func(f *big.Float) string {
		integer, accuracy := f.Int(nil)
		if integer == nil || accuracy != big.Exact || f.Acc() != big.Exact ||
			f.Mode() != big.ToNearestEven || integer.Sign() == 0 && f.Signbit() {
			t.Fatal("float lost exact integer value or default mode")
		}
		return fmt.Sprintf("%s/%d", integer, f.Prec())
	}
	rational := func(r *big.Rat) string {
		return fmt.Sprintf("%s/%s", r.Num(), r.Denom())
	}
	for row := 0; ; row++ {
		var c struct {
			X    int64
			Big  string
			Tag  int
			Want string
		}
		err := reader.Decode(&c)
		if err == io.EOF && row > 0 {
			t.Logf("%d pinned native numeric conversion rows", row)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		zeroI, zeroR, zeroF := newInt(), newRat(), newFloat()
		if zeroI == newInt() || zeroR == newRat() || zeroF == newFloat() {
			t.Fatal("constructor reused its destination")
		}
		i, r, f := i64toi(int64Val(c.X)), i64tor(int64Val(c.X)), i64tof(int64Val(c.X))
		if i.val == i64toi(int64Val(c.X)).val || r.val == i64tor(int64Val(c.X)).val ||
			f.val == i64tof(int64Val(c.X)).val {
			t.Fatal("conversion reused its destination")
		}
		source, ok := new(big.Int).SetString(c.Big, 10)
		if !ok {
			t.Fatal("bad integer fixture")
		}
		converted := itor(intVal{source})
		if converted.val == itor(intVal{source}).val || source.String() != c.Big {
			t.Fatal("rational conversion reused a destination or changed input")
		}
		source.SetInt64(91)
		references := []Value{nil, (*stringVal)(nil), MakeBool(true), MakeUnknown(), i}
		ref := references[c.Tag]
		z := vtoc(ref)
		if z.re != ref {
			t.Fatal("complex conversion changed real identity")
		}
		imaginary, isSmall := z.im.(int64Val)
		if !isSmall {
			t.Fatal("complex conversion changed the imaginary representation")
		}
		real := strconv.Itoa(c.Tag)
		if z.re == nil {
			real = "nil"
		} else if z.re == (*stringVal)(nil) {
			real = "typedNil"
		}
		got := strings.Join([]string{
			zeroI.String(), rational(zeroR), float(zeroF), i.val.String(), rational(r.val),
			float(f.val), rational(converted.val), fmt.Sprintf("%s:%d", real, imaginary),
		}, "|")
		if got != c.Want {
			t.Fatalf("row %d: Go %q want %q", row, got, c.Want)
		}
	}
}
