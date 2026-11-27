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

func TestGoLeanHasMeta(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_FILEPATH_HAS_META_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 || fields[0][0] != 'x' {
			t.Fatalf("invalid fixture %d", rows)
		}
		path, err := hex.DecodeString(fields[0][1:])
		if err != nil {
			t.Fatal(err)
		}
		want, err := strconv.ParseBool(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		if got := filepath.GoLeanHasMeta(string(path)); got != want {
			t.Fatalf("hasMeta(%q) = %t, want %t", path, got, want)
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
