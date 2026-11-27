package filepath_test

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoLeanPrivateJoin(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_PRIVATE_JOIN_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	read := func(field string) string {
		raw, err := hex.DecodeString(strings.TrimPrefix(field, "x"))
		if err != nil || !strings.HasPrefix(field, "x") {
			t.Fatalf("invalid fixture field %q", field)
		}
		return string(raw)
	}
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			t.Fatal("missing expected path")
		}
		want := read(fields[0])
		var elements []string
		for _, field := range fields[1:] {
			elements = append(elements, read(field))
		}
		if got := filepath.GoLeanPrivateJoin(elements); got != want {
			t.Fatalf("join(%q)=%q, want %q", elements, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil || count == 0 {
		t.Fatalf("read %d fixtures: %v", count, err)
	}
	t.Logf("private join matched %d pinned Go rows", count)
}
