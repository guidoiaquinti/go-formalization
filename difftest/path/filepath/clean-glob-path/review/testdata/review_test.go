package filepath_test

import (
	"path/filepath"
	"testing"
	"unicode/utf8"
)

func requireTrailingSeparator(path string) string {
	if len(path) < 2 || path[len(path)-1] == filepath.Separator {
		if path != "" && path != string(filepath.Separator) {
			return path[:len(path)-1]
		}
		return filepath.GoLeanCleanGlobPath(path)
	}
	return path
}

func removeLastRune(path string) string {
	if path == "" || path == string(filepath.Separator) {
		return filepath.GoLeanCleanGlobPath(path)
	}
	_, width := utf8.DecodeLastRuneInString(path)
	return path[:len(path)-width]
}

func TestGoLeanCleanGlobPathReview(t *testing.T) {
	native := string(filepath.Separator)
	other := "\\"
	if native == other {
		other = "/"
	}
	fixtures := []struct{ path, want string }{
		{"", "."}, {native, native}, {other, ""}, {"a", ""},
		{"abc", "ab"}, {"a" + native + native, "a" + native},
		{"\xffé", "\xff\xc3"}, {"a\xc3", "a"}, {"abc\x00", "abc"},
	}
	for _, c := range []struct {
		name  string
		clean func(string) string
		want  int
	}{
		{"Go", filepath.GoLeanCleanGlobPath, 0},
		{"Require an actual trailing separator", requireTrailingSeparator, 4},
		{"Drop the last decoded rune", removeLastRune, 1},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.clean(f.path); got != f.want {
				bad++
				if c.want == 0 {
					t.Errorf("cleanGlobPath(%x)=%x want %x", f.path, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
