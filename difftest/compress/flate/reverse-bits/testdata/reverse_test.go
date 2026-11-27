package flate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestGoLeanReverseBits(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_REVERSE_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	s := bufio.NewReader(f)
	for row := 0; ; row++ {
		var x, want uint16
		var n byte
		parsed, err := fmt.Fscan(s, &x, &n, &want)
		if parsed == 0 && err == io.EOF && row > 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if got := reverseBits(x, n); got != want {
			t.Fatalf("row %d, word %#04x count %d: got %d, want %d", row, x, n, got, want)
		}
	}
}
