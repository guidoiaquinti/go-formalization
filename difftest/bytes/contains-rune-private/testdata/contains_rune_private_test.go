package bytes_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanContainsRunePrivate(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTES_CONTAINS_RUNE_PRIVATE_CASES"))
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
		raw, err := hex.DecodeString(fields[0])
		if err != nil {
			t.Fatal(err)
		}
		seek, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		if got := bytes.GoLeanContainsRunePrivate(string(raw), rune(seek)); strconv.FormatBool(got) != fields[2] {
			t.Fatalf("containsRune(%x, %d) = %t, want %s", raw, seek, got, fields[2])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private bytes.containsRune matched %d oracle rows\n", count)
}
