package bytes_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanExplode(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_BYTES_EXPLODE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var count int
	for {
		var limit, spare int
		var text, want string
		n, err := fmt.Fscan(reader, &limit, &spare, &text, &want)
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
		var src []byte
		if len(s)+spare != 0 {
			src = make([]byte, len(s), len(s)+spare)
			copy(src, s)
		}
		got := goLeanExplodeResult(src, limit)
		if got != want {
			t.Fatalf("input %q, limit %d: Go=%q want=%q", s, limit, got, want)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no explode comparison rows")
	}
	t.Logf("private explode matched %d comparison rows", count)
}

func goLeanExplodeResult(s []byte, limit int) string {
	parts := bytes.GoLeanExplode(s, limit)
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
