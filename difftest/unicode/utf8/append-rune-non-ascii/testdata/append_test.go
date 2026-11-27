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

func TestGoLeanAppendRuneNonASCII(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_UTF8_APPEND_NON_ASCII_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := bufio.NewReader(f)
	for {
		var encoded, want string
		var r rune
		n, err := fmt.Fscan(in, &r, &encoded, &want)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		dst, err := hex.DecodeString(encoded[1:])
		if err != nil {
			t.Fatal(err)
		}
		got := utf8.GoLeanAppendRuneNonASCII(dst, r)
		if hex.EncodeToString(got) != want {
			t.Fatalf("appendRuneNonASCII(%s,%d) = %x, want %s", encoded, r, got, want)
		}
	}
}
