package strings_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanTrimLeftByte(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_TRIM_LEFT_BYTE_CASES"))
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
		cut, err := strconv.ParseUint(parts[0], 10, 8)
		if err != nil {
			t.Fatal(err)
		}
		input, err := hex.DecodeString(parts[1])
		if err != nil {
			t.Fatal(err)
		}
		got := strings.GoLeanTrimLeftByte(string(input), byte(cut))
		if encoded := hex.EncodeToString([]byte(got)); encoded != parts[2] {
			t.Fatalf("trimLeftByte(%x,%02x) = %x, want %s", input, cut, got, parts[2])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.trimLeftByte matched %d oracle rows\n", count)
}
