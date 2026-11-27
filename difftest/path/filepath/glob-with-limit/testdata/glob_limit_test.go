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

func TestGoLeanGlobWithLimit(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_GLOB_LIMIT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	t.Chdir(os.Getenv("GOLEAN_GLOB_LIMIT_DIR"))
	s := bufio.NewScanner(file)
	count := 0
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) != 3 {
			t.Fatalf("invalid fixture %q", s.Text())
		}
		depth, err := strconv.Atoi(fields[0])
		if err != nil {
			t.Fatal(err)
		}
		pattern, err := hex.DecodeString(strings.TrimPrefix(fields[1], "x"))
		if err != nil {
			t.Fatal(err)
		}
		got, err := filepath.GoLeanGlobWithLimit(string(pattern), depth)
		wire := "nil"
		if err != nil {
			if got != nil || err != filepath.ErrBadPattern {
				t.Fatalf("unexpected result %q, %v", got, err)
			}
			wire = "bad"
		} else if got != nil {
			parts := make([]string, len(got))
			for i, p := range got {
				values := make([]string, len(p))
				for j := range len(p) {
					values[j] = strconv.Itoa(int(p[j]))
				}
				parts[i] = strings.Join(values, ",")
			}
			wire = strconv.Itoa(len(got)) + ":" + strings.Join(parts, ";")
		}
		if wire != fields[2] {
			t.Fatalf("globWithLimit(%q,%d)=%q, want %q", pattern, depth, wire, fields[2])
		}
		count++
	}
	if err := s.Err(); err != nil || count == 0 {
		t.Fatalf("read %d cases: %v", count, err)
	}
	t.Logf("private globWithLimit matched %d pinned native cases", count)
}
