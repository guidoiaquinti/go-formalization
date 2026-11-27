package filepath_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGoLeanWalkDirReview(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{"a", "a/nested"} {
		if err := os.Mkdir(filepath.Join(root, dir), 0700); err != nil {
			t.Fatal(err)
		}
	}
	for _, file := range []string{"z", "b", "a/1", "a/2", "a/nested/file"} {
		if err := os.WriteFile(filepath.Join(root, file), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	stat, err := os.Lstat(root)
	if err != nil {
		t.Fatal(err)
	}
	wrapped := fmt.Errorf("wrapped: %w", filepath.SkipDir)
	for _, c := range []struct {
		response string
		decision error
		want     []string
		result   error
	}{
		{"a", filepath.SkipDir, []string{".", "a", "b", "z"}, nil},
		{"a/1", filepath.SkipDir, []string{".", "a", "a/1", "b", "z"}, nil},
		{"a/1", filepath.SkipAll, []string{".", "a", "a/1"}, filepath.SkipAll},
		{"a/1", wrapped, []string{".", "a", "a/1"}, wrapped},
	} {
		var events []string
		got := filepath.GoLeanPrivateWalkDir(root, fs.FileInfoToDirEntry(stat), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				t.Fatal(err)
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				t.Fatal(err)
			}
			rel = filepath.ToSlash(rel)
			events = append(events, rel)
			if rel == c.response {
				return c.decision
			}
			return nil
		})
		if got != c.result || !reflect.DeepEqual(events, c.want) {
			t.Fatalf("decision %v from %s: %v %v", c.decision, c.response, events, got)
		}
	}

	// Move a directory into a second temp parent before ReadDir. Restoring
	// it in the error callback must not trigger a rescan or visit its child.
	moved := filepath.Join(t.TempDir(), "moved")
	path := filepath.Join(root, "a")
	var visits []string
	got := filepath.GoLeanPrivateWalkDir(path, entryAt(t, path), func(p string, d fs.DirEntry, err error) error {
		if p != path || !d.IsDir() {
			t.Fatalf("unexpected visit %q directory=%t", p, d.IsDir())
		}
		if err == nil {
			visits = append(visits, "first")
			if err := os.Rename(path, moved); err != nil {
				t.Fatal(err)
			}
		} else {
			if !errors.Is(err, fs.ErrNotExist) {
				t.Fatal(err)
			}
			visits = append(visits, "read error")
			if err := os.Rename(moved, path); err != nil {
				t.Fatal(err)
			}
		}
		return nil
	})
	if got != nil || !reflect.DeepEqual(visits, []string{"first", "read error"}) {
		t.Fatalf("read error sequence %v result=%v", visits, got)
	}
}

func entryAt(t *testing.T, path string) fs.DirEntry {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	return fs.FileInfoToDirEntry(info)
}
