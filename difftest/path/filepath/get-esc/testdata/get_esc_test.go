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

func TestGoLeanGetEsc(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_FILEPATH_GET_ESC_CASES"))
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
		if len(fields) != 4 {
			t.Fatalf("invalid fixture %d", rows)
		}
		chunk, rest := decode(fields[0]), decode(fields[3])
		expectedRune, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		bad, err := strconv.ParseBool(fields[2])
		if err != nil {
			t.Fatal(err)
		}
		var expectedErr error
		if bad {
			expectedErr = filepath.ErrBadPattern
		}
		got, tail, err := filepath.GoLeanGetEsc(chunk)
		if got != rune(expectedRune) || tail != rest || err != expectedErr {
			t.Fatalf("getEsc(%q)=(%U,%q,%v), want (%U,%q,%v)", chunk, got, tail, err, expectedRune, rest, expectedErr)
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
