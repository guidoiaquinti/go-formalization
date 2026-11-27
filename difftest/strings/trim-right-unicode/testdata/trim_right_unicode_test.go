package strings_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGoLeanTrimRightUnicode(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_TRIM_RIGHT_UNICODE_CASES"))
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
		cutset, err := hex.DecodeString(parts[0])
		if err != nil {
			t.Fatal(err)
		}
		input, err := hex.DecodeString(parts[1])
		if err != nil {
			t.Fatal(err)
		}
		got := strings.GoLeanTrimRightUnicode(string(input), string(cutset))
		if encoded := hex.EncodeToString([]byte(got)); encoded != parts[2] {
			t.Fatalf("trimRightUnicode(%x,%x): Go=%x want=%s", input, cutset, got, parts[2])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.trimRightUnicode matched %d oracle rows\n", count)
}
