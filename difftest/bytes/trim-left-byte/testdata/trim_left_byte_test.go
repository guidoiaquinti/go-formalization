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

func TestGoLeanTrimLeftByte(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTES_TRIM_LEFT_BYTE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) != 5 {
			t.Fatalf("invalid fixture row %q", scanner.Text())
		}
		cut, err := strconv.ParseUint(parts[0], 10, 8)
		if err != nil {
			t.Fatal(err)
		}
		capacity, err := strconv.Atoi(parts[1])
		if err != nil {
			t.Fatal(err)
		}
		raw, err := hex.DecodeString(parts[2])
		if err != nil {
			t.Fatal(err)
		}
		var input []byte
		if capacity >= 0 {
			input = make([]byte, len(raw), capacity)
			copy(input, raw)
		}
		wantCap, err := strconv.Atoi(parts[3])
		if err != nil {
			t.Fatal(err)
		}
		wantNil := parts[4] == ""
		got := bytes.GoLeanTrimLeftByte(input, byte(cut))
		if (got == nil) != wantNil || cap(got) != wantCap || hex.EncodeToString(got) != parts[4] {
			t.Fatalf("trimLeftByte(%x cap=%d,%02x) = %x cap=%d nil=%t, want %s cap=%d nil=%t",
				input, capacity, cut, got, cap(got), got == nil, parts[4], wantCap, wantNil)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private bytes.trimLeftByte matched %d oracle rows\n", count)
}
