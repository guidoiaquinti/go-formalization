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

func TestGoLeanMatchChunk(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_FILEPATH_MATCH_CHUNK_CASES"))
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
	readBool := func(value string) bool {
		result, err := strconv.ParseBool(value)
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		f := strings.Fields(scanner.Text())
		if len(f) != 5 {
			t.Fatalf("invalid fixture %d", rows)
		}
		pattern, input, want := decode(f[0]), decode(f[1]), decode(f[4])
		ok, bad := readBool(f[2]), readBool(f[3])
		var wantErr error
		if bad {
			wantErr = filepath.ErrBadPattern
		}
		got, matched, err := filepath.GoLeanMatchChunk(pattern, input)
		if got != want || matched != ok || err != wantErr {
			t.Fatalf("matchChunk(%q,%q)=(%q,%t,%v), want (%q,%t,%v)",
				pattern, input, got, matched, err, want, ok, wantErr)
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
