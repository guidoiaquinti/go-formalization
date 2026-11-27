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

func TestGoLeanTrimLeftUnicode(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTES_TRIM_LEFT_UNICODE_CASES"))
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
		cutset, err := hex.DecodeString(parts[0])
		if err != nil {
			t.Fatal(err)
		}
		raw, err := hex.DecodeString(parts[1])
		if err != nil {
			t.Fatal(err)
		}
		spare, err := strconv.Atoi(parts[2])
		if err != nil {
			t.Fatal(err)
		}
		wantCap, err := strconv.Atoi(parts[4])
		if err != nil {
			t.Fatal(err)
		}
		var input []byte
		if spare >= 0 {
			input = make([]byte, len(raw), len(raw)+spare)
			copy(input, raw)
		}
		wantNil := parts[3] == ""
		got := bytes.GoLeanTrimLeftUnicode(input, string(cutset))
		if (got == nil) != wantNil || cap(got) != wantCap || hex.EncodeToString(got) != parts[3] {
			t.Fatalf("trimLeftUnicode(%x cap=%d, %x): Go=%x cap=%d nil=%t want=%s cap=%d nil=%t",
				input, cap(input), cutset, got, cap(got), got == nil, parts[3], wantCap, wantNil)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private bytes.trimLeftUnicode matched %d oracle rows\n", count)
}
