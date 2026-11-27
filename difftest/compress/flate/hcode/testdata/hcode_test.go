package flate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestGoLeanHCode(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_HCODE_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := bufio.NewReader(f)
	for i := 0; ; i++ {
		var raw, packed, updated uint32
		var code uint16
		var length, wantLength uint8
		var wantCode uint64
		var wantZero bool
		n, err := fmt.Fscan(in, &raw, &code, &length, &packed, &wantLength, &wantCode, &wantZero, &updated)
		if n == 0 && err == io.EOF && i > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		original := hcode(raw)
		state := original
		state.set(code, length)
		if uint32(newhcode(code, length)) != packed || original.len() != wantLength ||
			original.code64() != wantCode || original.zero() != wantZero || uint32(state) != updated {
			t.Fatalf("row %d: word %#x code %#x length %d: got %d %d %d %t %d",
				i, raw, code, length, newhcode(code, length), original.len(), original.code64(), original.zero(), state)
		}
	}
}
