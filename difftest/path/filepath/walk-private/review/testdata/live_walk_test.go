package filepath_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGoLeanPrivateWalkLiveFailure(t *testing.T) {
	sentinel := errors.New("review callback failure")
	for _, want := range []error{nil, filepath.SkipDir, filepath.SkipAll, sentinel} {
		parent := t.TempDir()
		root := filepath.Join(parent, "gone")
		if err := os.Mkdir(root, 0700); err != nil {
			t.Fatal(err)
		}
		info, err := os.Lstat(root)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(root); err != nil {
			t.Fatal(err)
		}
		callbacks := 0
		got := filepath.GoLeanPrivateWalk(root, info, func(path string, entry fs.FileInfo, err error) error {
			callbacks++
			if path != root || entry != info || !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("root read failure callback: %q %v %v", path, entry, err)
			}
			return want
		})
		if callbacks != 1 || got != want {
			t.Fatalf("root read failure returned %v after %d callbacks, want %v", got, callbacks, want)
		}
	}
	for _, reply := range []error{nil, filepath.SkipDir, filepath.SkipAll, sentinel} {
		root := t.TempDir()
		for _, name := range []string{"a", "z"} {
			if err := os.WriteFile(filepath.Join(root, name), nil, 0600); err != nil {
				t.Fatal(err)
			}
		}
		info, err := os.Lstat(root)
		if err != nil {
			t.Fatal(err)
		}
		var events []string
		got := filepath.GoLeanPrivateWalk(root, info, func(path string, entry fs.FileInfo, err error) error {
			switch path {
			case root:
				if err != nil || entry != info {
					t.Fatalf("ordinary root callback: %v %v", entry, err)
				}
				events = append(events, "root")
				// Names have already been read, but the private walker has not lstat'ed a.
				if err := os.Remove(filepath.Join(root, "a")); err != nil {
					t.Fatal(err)
				}
			case filepath.Join(root, "a"):
				if entry != nil || !errors.Is(err, fs.ErrNotExist) {
					t.Fatalf("child stat failure: %v %v", entry, err)
				}
				events = append(events, "a:error")
				return reply
			case filepath.Join(root, "z"):
				if entry == nil || err != nil {
					t.Fatalf("ordinary last child: %v %v", entry, err)
				}
				events = append(events, "z")
			default:
				t.Fatalf("unexpected callback %q", path)
			}
			return nil
		})
		want := reply
		expected := []string{"root", "a:error"}
		if reply == nil || reply == filepath.SkipDir {
			want = nil
			expected = append(expected, "z")
		}
		if got != want || !reflect.DeepEqual(events, expected) {
			t.Fatalf("child stat failure returned %v with %v, want %v with %v", got, events, want, expected)
		}
	}
}
