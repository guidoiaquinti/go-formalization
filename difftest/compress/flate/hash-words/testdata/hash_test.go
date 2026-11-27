package flate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestGoLeanHashWords(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_HASH_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := bufio.NewReader(f)
	for row := 0; ; row++ {
		var word uint64
		var bits, count uint8
		var want4, wantN uint32
		n, err := fmt.Fscan(input, &word, &bits, &count, &want4, &wantN)
		if n == 0 && err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		if got4, gotN := hash4u(uint32(word), bits), hashLen(word, bits, count); got4 != want4 || gotN != wantN {
			t.Fatalf("row %d (%#016x,%d,%d): Go %d %d want %d %d", row, word, bits, count, got4, gotN, want4, wantN)
		}
	}
}
