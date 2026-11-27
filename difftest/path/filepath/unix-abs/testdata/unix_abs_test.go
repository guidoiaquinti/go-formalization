package filepath_test

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestGoLeanUnixAbs(t *testing.T) {
	fixtures, err := os.Open(os.Getenv("GOLEAN_UNIX_ABS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer fixtures.Close()
	scanner := bufio.NewScanner(fixtures)
	decode := func(field string) string {
		if len(field) == 0 || field[0] != 'x' {
			t.Fatalf("invalid fixture field %q", field)
		}
		raw, err := hex.DecodeString(field[1:])
		if err != nil {
			t.Fatal(err)
		}
		return string(raw)
	}
	if !scanner.Scan() {
		t.Fatal("missing working directory")
	}
	t.Chdir(decode(scanner.Text()))
	count := 0
	for scanner.Scan() {
		input := decode(scanner.Text())
		if !scanner.Scan() {
			t.Fatal("missing expected path")
		}
		want := decode(scanner.Text())
		if got, err := filepath.GoLeanUnixAbs(input); err != nil || got != want {
			t.Fatalf("unixAbs(%q)=(%q,%v), want %q", input, got, err, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil || count == 0 {
		t.Fatalf("read %d fixtures: %v", count, err)
	}
	t.Logf("private unixAbs matched %d native rows", count)
}
