package strings_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGoLeanSingleStringReplacerReplace(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_SINGLE_REPLACER_REPLACE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) != 4 {
			t.Fatalf("invalid fixture row %q", scanner.Text())
		}
		var decoded [3]string
		for i := range decoded {
			raw, err := hex.DecodeString(fields[i])
			if err != nil {
				t.Fatal(err)
			}
			decoded[i] = string(raw)
		}
		got := strings.GoLeanSingleStringReplacerReplace(decoded[0], decoded[1], decoded[2])
		if encoded := hex.EncodeToString([]byte(got)); encoded != fields[3] {
			t.Fatalf("Replace(%x,%x,%x): Go=%x want=%s", decoded[0], decoded[1], decoded[2], got, fields[3])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.singleStringReplacer.Replace matched %d oracle rows\n", count)
}
