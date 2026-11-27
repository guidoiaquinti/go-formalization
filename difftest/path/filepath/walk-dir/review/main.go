package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	cases := []struct{ input, want string }{
		{"0 0 0 0 0 1 0 0 0", "ok|root:d,root/a:f,root/d:d,root/d/b:f,root/z:f"},
		{"0 0 0 0 0 2 0 0 0", "ok|root:d,root/a:f,root/d:d,root/d/b:f"},
		{"0 1 0 0 0 0 0 0 0", "ok|root:d,root:d:7,root/a:f,root/d:d,root/d/b:f,root/d/c:f,root/z:f"},
		{"0 1 0 0 0 0 0 0 3", "9|root:d,root:d:7"},
		{"0 2 0 0 0 0 0 0 2", "ok|missing:nil:8"},
		{"0 2 0 0 0 0 0 0 3", "9|missing:nil:8"},
	}
	var model bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&model, c.input)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-path-filepath-walk-dir")
	cmd.Stdin = &model
	out, err := cmd.Output()
	must(err)
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		panic("unexpected model row count")
	}
	for i, c := range cases {
		if rows[i] != c.want {
			panic(fmt.Sprintf("row %d got %q want %q", i, rows[i], c.want))
		}
	}
	dir, err := os.MkdirTemp("", "golean-review-walkdir-")
	must(err)
	defer os.RemoveAll(dir)
	must(os.Mkdir(filepath.Join(dir, "d"), 0700))
	for _, p := range []string{"z", "d/c", "d/b", "a"} {
		must(os.WriteFile(filepath.Join(dir, filepath.FromSlash(p)), nil, 0600))
	}
	for _, scenario := range []struct {
		sentinel error
		trace    string
	}{
		{filepath.SkipDir, ".,a,d,d/b,z"}, {filepath.SkipAll, ".,a,d,d/b"},
	} {
		var paths []string
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			must(err)
			rel, err := filepath.Rel(dir, path)
			must(err)
			rel = filepath.ToSlash(rel)
			paths = append(paths, rel)
			if rel == "d/b" {
				return scenario.sentinel
			}
			return nil
		})
		must(err)
		if strings.Join(paths, ",") != scenario.trace {
			panic(fmt.Sprintf("native walk trace %q", paths))
		}
	}
	missing := filepath.Join(dir, "missing")
	wrapped := fmt.Errorf("callback wraps skip: %w", filepath.SkipAll)
	for _, want := range []error{nil, filepath.SkipDir, filepath.SkipAll, wrapped} {
		called := 0
		got := filepath.WalkDir(missing, func(path string, d fs.DirEntry, err error) error {
			called++
			if path != missing || d != nil || !os.IsNotExist(err) {
				panic("missing-root callback differed")
			}
			return want
		})
		expected := want
		if want == filepath.SkipDir || want == filepath.SkipAll {
			expected = nil
		}
		if called != 1 || got != expected {
			panic("missing-root return or count differed")
		}
	}
	fmt.Printf("WalkDir passed %d independent model rows and six native sequences\n", len(cases))
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
