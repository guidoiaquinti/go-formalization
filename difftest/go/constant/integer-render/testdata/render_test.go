package constant

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"testing"
)

func TestGoLeanIntegerRender(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_CONSTANT_INTEGER_RENDER"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := json.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Word    int64
			Big     string
			Present bool
			Want    string
		}
		err := reader.Decode(&c)
		if err == io.EOF && row > 0 {
			t.Logf("%d pinned integer render rows", row)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		n, ok := new(big.Int).SetString(c.Big, 10)
		if !ok {
			t.Fatal("bad integer fixture")
		}
		large := intVal{n}
		if !c.Present {
			large.val = nil
		}
		x := int64Val(c.Word)
		got := fmt.Sprintf("%s|%s|%s|%s|%t|%t|%t", x.String(), x.ExactString(),
			large.String(), large.ExactString(), is32bit(c.Word), is63bit(c.Word), smallInt(n))
		if got != c.Want || n.String() != c.Big {
			t.Fatalf("row %d: Go %q want %q, operand %q", row, got, c.Want, n)
		}
	}
}
