package utf8_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

func observe(dst []byte, r rune) (status string) {
	defer func() {
		if p := recover(); p != nil {
			if err, ok := p.(runtime.Error); ok && strings.Contains(err.Error(), "index out of range") {
				status = "bounds"
				return
			}
			panic(p)
		}
	}()
	return strconv.Itoa(utf8.GoLeanEncodeRuneNonASCII(dst, r))
}

func TestGoLeanEncodeRuneNonASCII(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_UTF8_ENCODE_NON_ASCII_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewReader(file)
	for {
		var encoded, status, want string
		var r rune
		n, err := fmt.Fscan(input, &r, &encoded, &status, &want)
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
		got := observe(dst, r)
		if result := "x" + hex.EncodeToString(dst); got != status || result != want {
			t.Fatalf("encodeRuneNonASCII(%s,%d) = %s %s, want %s %s", encoded, r, got, result, status, want)
		}
	}
}
