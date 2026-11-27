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

func TestGoLeanScanChunk(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_FILEPATH_SCAN_CHUNK_CASES"))
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
		pattern, chunk, rest := decode(fields[0]), decode(fields[2]), decode(fields[3])
		star, err := strconv.ParseBool(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		s, c, r := filepath.GoLeanScanChunk(pattern)
		if s != star || c != chunk || r != rest {
			t.Fatalf("scanChunk(%q)=(%t,%q,%q), want (%t,%q,%q)", pattern, s, c, r, star, chunk, rest)
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
