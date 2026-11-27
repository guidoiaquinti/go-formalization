package filepath_test

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGoLeanPrivateSplitList(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_PRIVATE_SPLIT_LIST_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decode := func(field string) string {
		if !strings.HasPrefix(field, "x") {
			t.Fatalf("invalid fixture field %q", field)
		}
		decoded, err := hex.DecodeString(field[1:])
		if err != nil {
			t.Fatal(err)
		}
		return string(decoded)
	}
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		path := decode(fields[0])
		want := make([]string, len(fields)-1)
		for i, field := range fields[1:] {
			want[i] = decode(field)
		}
		if got := filepath.GoLeanPrivateSplitList(path); !reflect.DeepEqual(got, want) {
			t.Fatalf("splitList(%q)=%#v, want %#v", path, got, want)
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("no fixtures")
	}
	t.Logf("private splitList matched %d pinned native rows", rows)
}
