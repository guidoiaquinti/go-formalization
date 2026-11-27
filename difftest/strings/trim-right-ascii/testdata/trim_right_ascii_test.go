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

func TestGoLeanTrimRightASCII(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_TRIM_RIGHT_ASCII_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) != 3 {
			t.Fatalf("invalid fixture row %q", scanner.Text())
		}
		nilTable := parts[0] == "-"
		var members []byte
		if !nilTable {
			members, err = hex.DecodeString(parts[0])
			if err != nil {
				t.Fatal(err)
			}
		}
		input, err := hex.DecodeString(parts[1])
		if err != nil {
			t.Fatal(err)
		}
		if got := observeTrimRightASCII(string(input), members, nilTable); got != parts[2] {
			t.Fatalf("trimRightASCII(%x, %q): Go=%s want=%s", input, parts[0], got, parts[2])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.trimRightASCII matched %d oracle rows\n", count)
}

func observeTrimRightASCII(input string, members []byte, nilTable bool) (output string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			output = "panic"
		}
	}()
	return hex.EncodeToString([]byte(strings.GoLeanTrimRightASCII(input, members, nilTable)))
}
