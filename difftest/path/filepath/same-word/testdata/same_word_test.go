package filepath_test

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanSameWord(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_FILEPATH_SAME_WORD_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decode := func(value string) string {
		if !strings.HasPrefix(value, "x") {
			t.Fatalf("invalid field %q", value)
		}
		raw, err := hex.DecodeString(value[1:])
		if err != nil {
			t.Fatal(err)
		}
		return string(raw)
	}
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 3 {
			t.Fatalf("invalid fixture %d", rows)
		}
		left, right := decode(fields[0]), decode(fields[1])
		want, err := strconv.ParseBool(fields[2])
		if err != nil {
			t.Fatal(err)
		}
		if got := filepath.GoLeanSameWord(left, right); got != want {
			t.Fatalf("sameWord(%q,%q)=%t, want %t", left, right, got, want)
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("no fixtures")
	}
}
