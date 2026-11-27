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

func wordResult[T string | []byte](data T) (out string) {
	defer func() {
		if p := recover(); p != nil {
			if err, ok := p.(runtime.Error); ok && strings.Contains(err.Error(), "index out of range") {
				out = "bounds"
				return
			}
			panic(p)
		}
	}()
	return strconv.FormatUint(uint64(utf8.GoLeanWord(data)), 10)
}

func TestGoLeanWord(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_UTF8_WORD_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewReader(file)
	for {
		var encoded, want string
		n, err := fmt.Fscan(input, &encoded, &want)
		if err == io.EOF && n == 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		data, err := hex.DecodeString(strings.TrimPrefix(encoded, "x"))
		if err != nil {
			t.Fatal(err)
		}
		if slice, str := wordResult(data), wordResult(string(data)); slice != want || str != want {
			t.Fatalf("word(%x): slice=%s string=%s want=%s", data, slice, str, want)
		}
	}
}
