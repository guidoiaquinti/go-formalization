package bytealg_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internal/bytealg"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGoLeanBytealgRuntimeMemory(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTEALG_RUNTIME_MEMORY_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Left, Right                           []byte
			Length, Size, LeftOffset, RightOffset int
			Alias                                 bool
			Want                                  string
		}
		if err := rows.Decode(&c); err == io.EOF {
			if i == 0 {
				t.Fatal("missing runtime memory cases")
			}
			t.Logf("both native entries matched %d rows", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		makeView := func(raw []byte, offset int, poison byte) []byte {
			if raw == nil {
				return nil
			}
			backing := bytes.Repeat([]byte{poison}, offset+len(raw)+64)
			copy(backing[offset:], raw)
			return backing[offset : offset+len(raw)]
		}
		left, right := makeView(c.Left, c.LeftOffset, 0xa5), makeView(c.Right, c.RightOffset, 0x55)
		if c.Alias {
			right = left
		}
		result := fmt.Sprintf("%s|%t", allocation(c.Length), bytealg.GoLeanRuntimeMemequal(left, right, c.Size))
		if result != c.Want {
			t.Fatalf("row %d %+v: Go=%s want=%s", i, c, result, c.Want)
		}
	}
}

func allocation(length int) (result string) {
	defer func() {
		if fault := recover(); fault != nil {
			err, ok := fault.(runtime.Error)
			if !ok || !strings.HasPrefix(err.Error(), "runtime error: makeslice: len out of range") {
				panic(fault)
			}
			result = "panic"
		}
	}()
	raw := bytealg.GoLeanMakeNoZero(length)
	result = fmt.Sprintf("%d:%d:%t", len(raw), cap(raw), raw != nil)
	// Initialize the whole allocation before reading a single byte.
	all := raw[:cap(raw)]
	for i := range all {
		all[i] = byte(i ^ length)
	}
	for i, b := range all {
		if b != byte(i^length) {
			panic("initialized allocation changed")
		}
	}
	return
}
