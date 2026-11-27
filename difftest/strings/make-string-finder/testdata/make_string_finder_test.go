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

func TestGoLeanMakeStringFinder(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_MAKE_STRING_FINDER_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		patternField, want, ok := strings.Cut(scanner.Text(), "\t")
		if !ok {
			t.Fatalf("invalid comparison row %q", scanner.Text())
		}
		pattern, err := hex.DecodeString(patternField)
		if err != nil {
			t.Fatal(err)
		}
		stored, bad, good := strings.GoLeanMakeStringFinder(string(pattern))
		if stored != string(pattern) || good == nil {
			t.Fatalf("makeStringFinder(%x): stored=%x nil good table=%t", pattern, stored, good == nil)
		}
		patternValues := make([]int, len(stored))
		for i := range len(stored) {
			patternValues[i] = int(stored[i])
		}
		got := formatIntegers(patternValues) + ":" + formatIntegers(bad[:]) + ":" + formatIntegers(good)
		if got != want {
			t.Fatalf("makeStringFinder(%x): Go=%s want=%s", pattern, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.makeStringFinder matched %d oracle rows\n", count)
}

func formatIntegers(values []int) string {
	fields := make([]string, len(values))
	for i, value := range values {
		fields[i] = strconv.Itoa(value)
	}
	return strings.Join(fields, ",")
}
