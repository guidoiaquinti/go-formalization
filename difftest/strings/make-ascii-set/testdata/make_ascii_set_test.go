package strings_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanMakeASCIISet(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_MAKE_ASCII_SET_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var count int
	for {
		var input, want string
		n, err := fmt.Fscan(reader, &input, &want)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		payload, err := hex.DecodeString(input[1:])
		if err != nil {
			t.Fatal(err)
		}
		table, valid := strings.GoLeanMakeASCIISet(string(payload))
		got := strconv.AppendBool(nil, valid)
		got = append(got, '|')
		for _, member := range table {
			digit := byte('0')
			if member {
				digit = '1'
			}
			got = append(got, digit)
		}
		if string(got) != want {
			t.Fatalf("input %q: Go=%s want=%s", payload, got, want)
		}
		count++
	}
	if count == 0 {
		t.Fatal("no private makeASCIISet comparison rows")
	}
	t.Logf("private makeASCIISet matched %d comparison rows", count)
}
