package filepath_test

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestGoLeanPrivateAbs(t *testing.T) {
	fixtures, err := os.Open(os.Getenv("GOLEAN_PRIVATE_ABS_CASES"))
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
		got, err := filepath.GoLeanPrivateAbs(input)
		if err != nil {
			if got != "" {
				t.Fatalf("abs(%q) returned %q with error %v", input, got, err)
			}
			got = "err:" + err.Error()
		} else {
			got = "ok:" + got
		}
		if got != want {
			t.Fatalf("abs(%q)=%q, want %q", input, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil || count == 0 {
		t.Fatalf("read %d fixtures: %v", count, err)
	}
	t.Logf("private abs matched %d native rows", count)
}
