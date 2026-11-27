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

func TestGoLeanLongestCommonSuffix(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_LONGEST_COMMON_SUFFIX_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var count int
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) != 3 {
			t.Fatalf("invalid comparison row %q", scanner.Text())
		}
		a, err := hex.DecodeString(fields[0])
		if err != nil {
			t.Fatal(err)
		}
		b, err := hex.DecodeString(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		got := strings.GoLeanLongestCommonSuffix(string(a), string(b))
		if strconv.Itoa(got) != fields[2] {
			t.Fatalf("longestCommonSuffix(%x,%x): Go=%d want=%s", a, b, got, fields[2])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.longestCommonSuffix matched %d oracle rows\n", count)
}
