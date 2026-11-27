package strings_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func observeByteStringReplace(table *[256][]byte, input string) (result string) {
	defer func() {
		if value := recover(); value != nil {
			fault, ok := value.(runtime.Error)
			if table != nil || !ok || fault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			result = "panic"
		}
	}()
	return hex.EncodeToString([]byte(strings.GoLeanByteStringReplacerReplace(table, input)))
}

func TestGoLeanByteStringReplacerReplace(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_BYTE_STRING_REPLACER_REPLACE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		row, want, ok := strings.Cut(scanner.Text(), "\t")
		if !ok {
			t.Fatalf("invalid fixture %q", scanner.Text())
		}
		fields := strings.Fields(row)
		if len(fields) < 2 {
			t.Fatalf("invalid fixture %q", scanner.Text())
		}
		numbers := make([]int, len(fields))
		for i, value := range fields {
			numbers[i], err = strconv.Atoi(value)
			if err != nil {
				t.Fatal(err)
			}
		}
		var table *[256][]byte
		if numbers[0] == 1 {
			table = new([256][]byte)
		}
		entries, raw := numbers[1], numbers[2:]
		for range entries {
			if len(raw) < 2 || raw[0] < 0 || raw[0] > 255 || raw[1] < 0 || raw[1] > len(raw)-2 || table == nil {
				t.Fatalf("invalid mapping fixture %q", row)
			}
			key, width := raw[0], raw[1]
			entry := make([]byte, width)
			for i := range entry {
				entry[i] = byte(raw[i+2])
			}
			table[key] = entry
			raw = raw[width+2:]
		}
		input := make([]byte, len(raw))
		for i, value := range raw {
			input[i] = byte(value)
		}
		if got := observeByteStringReplace(table, string(input)); got != want {
			t.Fatalf("byteStringReplacer.Replace(%s): Go=%s want=%s", row, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.byteStringReplacer.Replace matched %d oracle rows\n", count)
}
