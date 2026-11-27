package filepath_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGoLeanPrivateGlob(t *testing.T) {
	dir := t.TempDir()
	empty := filepath.Join(dir, "empty")
	if err := os.Mkdir(empty, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"z.txt", "alpha.txt", "b.txt", "a.go"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	cases := []struct {
		dir, pattern  string
		initial, want []string
		err           error
	}{
		{dir, "*.txt", []string{"existing"}, []string{"existing", filepath.Join(dir, "alpha.txt"), filepath.Join(dir, "b.txt"), filepath.Join(dir, "z.txt")}, nil},
		{dir, "!", nil, nil, nil},
		{dir, "!", []string{}, []string{}, nil},
		{dir, "[", nil, nil, filepath.ErrBadPattern},
		{empty, "[", nil, nil, nil},
		{filepath.Join(dir, "missing"), "[", nil, nil, nil},
		{filepath.Join(dir, "a.go"), "[", []string{"existing"}, []string{"existing"}, nil},
	}
	for _, c := range cases {
		got, err := filepath.GoLeanPrivateGlob(c.dir, c.pattern, c.initial)
		if err != c.err || !reflect.DeepEqual(got, c.want) {
			t.Fatalf("glob(%q,%q,%#v)=(%#v,%v), want (%#v,%v)", c.dir, c.pattern, c.initial, got, err, c.want, c.err)
		}
	}
}
