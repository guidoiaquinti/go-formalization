package bytes_test

import (
	"bufio"
	"bytes"
	"os"
	"runtime"
	"testing"
)

func TestGoLeanASCIISetContains(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTES_ASCII_SET_CONTAINS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var rows int
	for scanner.Scan() {
		want := scanner.Text()
		if len(want) != 256 {
			t.Fatalf("table %d: invalid fixture length %d", rows, len(want))
		}
		var table *[256]bool
		if want[0] != 'p' {
			table = new([256]bool)
			for b := range table {
				table[b] = want[b] == '1'
			}
		}
		for b := range 256 {
			if got := asciiSetContainsResult(table, byte(b)); got != want[b] {
				t.Fatalf("table %d, byte %d: Go=%c want=%c", rows, b, got, want[b])
			}
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("no asciiSet.contains comparison rows")
	}
	t.Logf("asciiSet.contains matched %d lookups", 256*rows)
}

func asciiSetContainsResult(table *[256]bool, b byte) (result byte) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			result = 'p'
		}
	}()
	if bytes.GoLeanASCIISetContains(table, b) {
		return '1'
	}
	return '0'
}
