package strings_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

func observeByteReplacerReplace(mapping *[256]byte, input string) (result string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			fault, ok := recovered.(runtime.Error)
			if mapping != nil || !ok || fault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(recovered)
			}
			result = "panic"
		}
	}()
	return hex.EncodeToString([]byte(strings.GoLeanByteReplacerReplace(mapping, input)))
}

func TestGoLeanByteReplacerReplace(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_BYTE_REPLACER_REPLACE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) != 3 {
			t.Fatalf("invalid fixture row %q", scanner.Text())
		}
		mapping, err := hex.DecodeString(fields[0])
		if err != nil || (len(mapping) != 0 && len(mapping) != 256) {
			t.Fatalf("invalid byte table %q", fields[0])
		}
		raw, err := hex.DecodeString(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		var table *[256]byte
		if len(mapping) != 0 {
			table = (*[256]byte)(mapping)
		}
		if got := observeByteReplacerReplace(table, string(raw)); got != fields[2] {
			t.Fatalf("byteReplacer.Replace(%x): Go=%s want=%s", raw, got, fields[2])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.byteReplacer.Replace matched %d oracle rows\n", count)
}
