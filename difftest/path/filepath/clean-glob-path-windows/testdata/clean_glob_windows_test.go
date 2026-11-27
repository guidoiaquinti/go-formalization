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

func TestGoLeanCleanGlobWindows(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_FILEPATH_CLEAN_GLOB_WINDOWS_CASES"))
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
		path, cleaned := decode(fields[0]), decode(fields[2])
		want, err := strconv.Atoi(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		prefix, output := filepath.GoLeanCleanGlobWindows(path)
		if prefix != want || output != cleaned {
			t.Fatalf("cleanGlobPathWindows(%q)=(%d,%q), want (%d,%q)", path, prefix, output, want, cleaned)
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
