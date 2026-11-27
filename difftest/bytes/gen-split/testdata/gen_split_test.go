package bytes_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanBytesGenSplit(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_BYTES_GEN_SPLIT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var count int
	for {
		var save, limit, spare int
		var text, separator, want string
		n, err := fmt.Fscan(reader, &save, &limit, &spare, &text, &separator, &want)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		s, err := hex.DecodeString(text[1:])
		if err != nil {
			t.Fatal(err)
		}
		sep, err := hex.DecodeString(separator[1:])
		if err != nil {
			t.Fatal(err)
		}
		src := make([]byte, len(s)+spare)
		copy(src, s)
		for i := len(s); i < len(src); i++ {
			src[i] = 0xa5
		}
		got := goLeanGenSplitResult(src[:len(s)], sep, save, limit)
		if got != want {
			t.Fatalf("input %q, sep %q, save %d, limit %d: Go=%q want=%q", s, sep, save, limit, got, want)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no genSplit comparison rows")
	}
	t.Logf("private genSplit matched %d comparison rows", count)
}

func goLeanGenSplitResult(s, sep []byte, save, limit int) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(runtime.Error); ok && strings.HasPrefix(e.Error(), "runtime error: slice bounds out of range") {
				result = "panic"
			} else {
				panic(p)
			}
		}
	}()
	parts := bytes.GoLeanGenSplit(s, sep, save, limit)
	if parts == nil {
		return "nil"
	}
	fields := make([]string, len(parts))
	for i, part := range parts {
		encoded := make([]string, len(part))
		for j := range len(part) {
			encoded[j] = strconv.Itoa(int(part[j]))
		}
		fields[i] = "f" + strconv.Itoa(cap(part)) + ":" + strings.Join(encoded, ",")
	}
	return "slice" + strings.Join(fields, ";")
}
