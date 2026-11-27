package bytealg_test

import (
	"bytes"
	"encoding/json"
	"internal/abi"
	"io"
	"os"
	"reflect"
	"runtime"
	"testing"
	"unsafe"
)

func TestGoLeanMemequalVarlen(t *testing.T) {
	// Use the compiler-generated closure descriptors. Calling the ABI alias
	// as a bare function would omit the closure register's size word.
	types := map[int]*abi.Type{}
	for _, exemplar := range []any{[3]byte{}, [5]byte{}, [7]byte{}, [9]byte{}, [15]byte{},
		[17]byte{}, [31]byte{}, [33]byte{}, [63]byte{}, [65]byte{}, [127]byte{}, [129]byte{}, [257]byte{}} {
		desc := abi.TypeOf(exemplar)
		code := runtime.FuncForPC(reflect.ValueOf(desc.Equal).Pointer())
		if code == nil || code.Name() != "runtime.memequal_varlen" {
			t.Fatalf("array width %d does not dispatch to runtime.memequal_varlen: %v", desc.Size(), code)
		}
		types[int(desc.Size())] = desc
	}
	file, err := os.Open(os.Getenv("GOLEAN_BYTEALG_MEMEQUAL_VARLEN_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Left, Right                   []byte
			Size, LeftOffset, RightOffset int
			Alias, Want                   bool
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing closure fixtures")
			}
			t.Logf("%d rows matched runtime.memequal_varlen via verified compiler closures", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		desc := types[c.Size]
		if desc == nil || len(c.Left) < c.Size || len(c.Right) < c.Size {
			t.Fatalf("row %d invalid closure size %d", i, c.Size)
		}
		makeView := func(raw []byte, offset int, poison byte) []byte {
			backing := bytes.Repeat([]byte{poison}, offset+len(raw)+64)
			copy(backing[offset:], raw)
			return backing[offset : offset+len(raw)]
		}
		left, right := makeView(c.Left, c.LeftOffset, 0xa5), makeView(c.Right, c.RightOffset, 0x5a)
		if c.Alias {
			right = left
		}
		got := desc.Equal(unsafe.Pointer(unsafe.SliceData(left)), unsafe.Pointer(unsafe.SliceData(right)))
		if got != c.Want {
			t.Fatalf("row %d width=%d a=%x b=%x: equal=%t want=%t", i, c.Size, left, right, got, c.Want)
		}
	}
}
