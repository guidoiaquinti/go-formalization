package bytes_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanIsSeparator(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_BYTES_IS_SEPARATOR_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		runeField, want, ok := strings.Cut(scanner.Text(), "\t")
		if !ok {
			t.Fatalf("invalid fixture row %q", scanner.Text())
		}
		value, err := strconv.ParseInt(runeField, 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		if got := bytes.GoLeanIsSeparator(rune(value)); strconv.FormatBool(got) != want {
			t.Fatalf("isSeparator(%d): Go=%t want=%s", value, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private bytes.isSeparator matched %d oracle rows\n", count)
}
