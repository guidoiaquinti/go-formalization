package flate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestGoLeanMatchCodes(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_MATCH_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := bufio.NewReader(f)
	for row := 0; ; row++ {
		var offset, wantOffset uint32
		var length, wantLength uint8
		n, err := fmt.Fscan(input, &offset, &length, &wantOffset, &wantLength)
		if n == 0 && err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		if offsetCode(offset) != wantOffset || lengthCode(length) != wantLength {
			t.Fatalf("row %d offset %d length %d: Go %d %d want %d %d", row, offset, length, offsetCode(offset), lengthCode(length), wantOffset, wantLength)
		}
	}
}
