package utf8_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"
	"unicode/utf8"
)

func TestGoLeanDecodeRuneSlow(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_UTF8_DECODE_SLICE_SLOW_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := bufio.NewReader(f)
	for {
		var encoded string
		var wantRune rune
		var wantWidth int
		n, err := fmt.Fscan(in, &encoded, &wantRune, &wantWidth)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		var data []byte
		if encoded != "nil" {
			data, err = hex.DecodeString(encoded[1:])
			if err != nil {
				t.Fatal(err)
			}
		}
		r, w := utf8.GoLeanDecodeRuneSlow(data)
		if r != wantRune || w != wantWidth {
			t.Fatalf("decodeRuneSlow(%x) = (%U, %d), want (%U, %d)", data, r, w, wantRune, wantWidth)
		}
	}
}
