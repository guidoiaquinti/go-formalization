package filepath_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGoLeanReadDirNames(t *testing.T) {
	base := t.TempDir()
	empty := filepath.Join(base, "empty")
	if err := os.Mkdir(empty, 0700); err != nil {
		t.Fatal(err)
	}
	full := filepath.Join(base, "full")
	if err := os.Mkdir(full, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"z", "aa", "é", "a", "B", ".hidden"} {
		if err := os.WriteFile(filepath.Join(full, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, c := range []struct {
		name string
		want []string
	}{
		{empty, []string{}},
		{full, []string{".hidden", "B", "a", "aa", "z", "é"}},
	} {
		got, err := filepath.GoLeanReadDirNames(c.name)
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Fatalf("%s: %#v, %v, want %#v", c.name, got, err, c.want)
		}
	}
	for _, name := range []string{filepath.Join(base, "missing"), filepath.Join(full, "z")} {
		got, err := filepath.GoLeanReadDirNames(name)
		if got != nil || err == nil {
			t.Fatalf("%s: %#v, %v", name, got, err)
		}
	}
}
