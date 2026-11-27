package constant

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"testing"
)

func TestGoLeanNumericSign(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_CONSTANT_SIGN_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	observe := func(value Value) (result string) {
		defer func() {
			if p := recover(); p != nil {
				if text, ok := p.(string); ok {
					result = "panic:" + text
				} else {
					panic(p)
				}
			}
		}()
		return fmt.Sprint(Sign(value))
	}
	for row := 0; ; row++ {
		var c struct {
			X, Y int64
			Mode int
			Want string
		}
		err := dec.Decode(&c)
		if err == io.EOF && row > 0 {
			t.Logf("%d pinned numeric kind and sign rows", row)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		x := int64Val(c.X)
		b := intVal{new(big.Int).Lsh(big.NewInt(c.X), 67)}
		q := ratVal{new(big.Rat).SetFrac64(c.X, 3)}
		f := floatVal{new(big.Float).SetPrec(512).SetInt64(c.X)}
		y := ratVal{new(big.Rat).SetFrac64(c.Y, 11)}
		z := complexVal{x, y}
		negated := floatVal{new(big.Float).Neg(f.val)}
		var alternate Value
		switch c.Mode {
		case 0:
			alternate = MakeUnknown()
		case 2:
			alternate = MakeBool(true)
		}
		fields := []string{fmt.Sprint(int(x.Kind())), fmt.Sprint(int(b.Kind())), fmt.Sprint(int(q.Kind())),
			fmt.Sprint(int(f.Kind())), fmt.Sprint(int(z.Kind())), observe(x), observe(b), observe(q), observe(f),
			observe(z), observe(complexVal{z, negated}), observe(alternate), observe(complexVal{alternate, y})}
		if got := strings.Join(fields, "|"); got != c.Want {
			t.Fatalf("row %d: Go %q want %q", row, got, c.Want)
		}
	}
}
