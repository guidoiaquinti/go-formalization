package strings_test

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanShouldUseASCIISet(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_SHOULD_USE_ASCII_SET_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		lengthField, want, ok := strings.Cut(scanner.Text(), "\t")
		if !ok {
			t.Fatalf("invalid fixture row %q", scanner.Text())
		}
		length, err := strconv.Atoi(lengthField)
		if err != nil {
			t.Fatal(err)
		}
		if got := strings.GoLeanShouldUseASCIISet(length); strconv.FormatBool(got) != want {
			t.Fatalf("shouldUseASCIISet(%d): Go=%t want=%s", length, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.shouldUseASCIISet matched %d oracle rows\n", count)
}
