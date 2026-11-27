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

func TestGoLeanMakeSingleStringReplacer(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_MAKE_SINGLE_STRING_REPLACER_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) != 3 {
			t.Fatalf("invalid comparison row %q", scanner.Text())
		}
		pattern, err := hex.DecodeString(fields[0])
		if err != nil {
			t.Fatal(err)
		}
		value, err := hex.DecodeString(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		stored, replacement, bad, good := strings.GoLeanMakeSingleStringReplacer(string(pattern), string(value))
		if stored != string(pattern) || replacement != string(value) || good == nil {
			t.Fatalf("makeSingleStringReplacer(%x,%x): stored=%x replacement=%x nil good table=%t",
				pattern, value, stored, replacement, good == nil)
		}
		got := formatBytes(stored) + ":" + formatBytes(replacement) + ":" +
			formatIntegers(bad[:]) + ":" + formatIntegers(good)
		if got != fields[2] {
			t.Fatalf("makeSingleStringReplacer(%x,%x): Go=%s want=%s", pattern, value, got, fields[2])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.makeSingleStringReplacer matched %d oracle rows\n", count)
}

func formatIntegers(values []int) string {
	fields := make([]string, len(values))
	for i, value := range values {
		fields[i] = strconv.Itoa(value)
	}
	return strings.Join(fields, ",")
}

func formatBytes(raw string) string {
	values := make([]int, len(raw))
	for i := range len(raw) {
		values[i] = int(raw[i])
	}
	return formatIntegers(values)
}
