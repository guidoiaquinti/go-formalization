package filepath_test

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoLeanPublicWalk(t *testing.T) {
	table := []string{"root", "root/a", "root/d", "root/d/b", "root/d/c", "root/z"}
	file, err := os.Open(os.Getenv("GOLEAN_PUBLIC_WALK_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	t.Chdir(t.TempDir())
	if err := os.MkdirAll("root/d", 0700); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"root/a", "root/d/b", "root/d/c", "root/z"} {
		if err := os.WriteFile(p, nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	other := errors.New("review stop")
	sc := bufio.NewScanner(file)
	count := 0
	for sc.Scan() {
		choices, want, ok := strings.Cut(sc.Text(), " ")
		if !ok || len(choices) != len(table) {
			t.Fatal("invalid fixture")
		}
		var events []string
		err := filepath.Walk("root", func(p string, d os.FileInfo, err error) error {
			if err != nil {
				t.Fatal(err)
			}
			p = filepath.ToSlash(p)
			kind := ":f"
			if d.IsDir() {
				kind = ":d"
			}
			events = append(events, p+kind)
			for i, name := range table {
				if p != name {
					continue
				}
				switch choices[i] {
				case '1':
					return filepath.SkipDir
				case '2':
					return filepath.SkipAll
				case '3':
					return other
				}
				return nil
			}
			t.Fatalf("unexpected callback %s", p)
			return nil
		})
		status := "ok"
		switch err {
		case nil:
		case filepath.SkipDir:
			status = "skipDir"
		case filepath.SkipAll:
			status = "skipAll"
		case other:
			status = "9"
		default:
			t.Fatal(err)
		}
		got := fmt.Sprintf("%s|%s", status, strings.Join(events, ","))
		if got != want {
			t.Fatalf("callbacks %s: Go %s Lean %s", choices, got, want)
		}
		count++
	}
	if err := sc.Err(); err != nil || count == 0 {
		t.Fatalf("read %d fixtures: %v", count, err)
	}
	for _, response := range []error{nil, filepath.SkipDir, filepath.SkipAll, other} {
		calls := 0
		err := filepath.Walk("missing", func(path string, entry os.FileInfo, failure error) error {
			calls++
			if path != "missing" || entry != nil || !os.IsNotExist(failure) {
				t.Fatalf("unexpected missing-root callback %q, %v, %v", path, entry, failure)
			}
			return response
		})
		var want error
		if response == other {
			want = other
		}
		if calls != 1 || err != want {
			t.Fatalf("missing root: %d callbacks and %v, want one and %v", calls, err, want)
		}
	}

	t.Logf("public Walk matched %d pinned native cases", count)
}
