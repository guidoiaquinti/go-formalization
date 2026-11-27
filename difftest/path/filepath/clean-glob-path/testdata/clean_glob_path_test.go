package filepath_test

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoLeanCleanGlobPath(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_CLEAN_GLOB_PATH_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 || fields[0][0] != 'x' || fields[1][0] != 'x' {
			t.Fatalf("invalid fixture %d", rows)
		}
		input, err := hex.DecodeString(fields[0][1:])
		if err != nil {
			t.Fatal(err)
		}
		want, err := hex.DecodeString(fields[1][1:])
		if err != nil {
			t.Fatal(err)
		}
		if got := filepath.GoLeanCleanGlobPath(string(input)); got != string(want) {
			t.Fatalf("cleanGlobPath(%q) = %q, want %q", input, got, want)
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
