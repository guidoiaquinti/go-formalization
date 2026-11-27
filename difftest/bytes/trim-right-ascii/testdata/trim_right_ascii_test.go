package bytes_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanTrimRightASCII(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTES_TRIM_RIGHT_ASCII_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) != 4 {
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
		spare, err := strconv.Atoi(parts[1])
		if err != nil {
			t.Fatal(err)
		}
		raw, err := hex.DecodeString(parts[2])
		if err != nil {
			t.Fatal(err)
		}
		var input []byte
		if spare >= 0 {
			input = make([]byte, len(raw), len(raw)+spare)
			copy(input, raw)
		}
		if got := observeTrimRightASCII(input, members, nilTable); got != parts[3] {
			t.Fatalf("trimRightASCII(%x cap=%d, %q): Go=%s want=%s", input, cap(input), parts[0], got, parts[3])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private bytes.trimRightASCII matched %d oracle rows\n", count)
}

func observeTrimRightASCII(input []byte, members []byte, nilTable bool) (output string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			output = "panic"
		}
	}()
	got := bytes.GoLeanTrimRightASCII(input, members, nilTable)
	if got == nil {
		return "nil"
	}
	row := strconv.AppendInt(nil, int64(cap(got)), 10)
	row = append(row, ':')
	for i, b := range got {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(b), 10)
	}
	return string(row)
}
