package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type scenario struct {
	depth   int
	pattern string
}

func main() {
	root, err := os.MkdirTemp("", "golean-glob-limit-*")
	must(err)
	defer os.RemoveAll(root)
	for _, dir := range []string{"alpha", "beta"} {
		must(os.Mkdir(filepath.Join(root, dir), 0700))
	}
	for _, name := range []string{"plain", "alpha/a.txt", "alpha/z.go", "beta/b.txt"} {
		must(os.WriteFile(filepath.Join(root, name), nil, 0600))
	}
	cases := []scenario{
		{0, ""}, {0, "plain"}, {0, "./plain"}, {0, "missing"}, {0, "alpha"}, {0, "."},
		{0, "["}, {0, "*/*.txt"}, {0, "*/[ab].txt"}, {0, "a*/?.txt"}, {0, "a*/*/*.txt"},
		{9999, "plain"}, {9999, "a*"}, {9999, "*/*.txt"}, {10000, "plain"},
	}
	for i := range 256 {
		cases = append(cases, scenario{0, "alpha/" + string([]byte{byte(i)})})
	}
	boundaries := len(cases)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for i := range 20000 {
			suffix := []string{"*", "*.txt", "a*", "?", "[ab].txt", "plain", "missing", "[", "a.txt", "z.go", "b.txt"}[r.Intn(11)]
			base := []string{"", "./", "alpha/", "beta/", "a*/", "*/", "missing/", "a*/z*/"}[r.Intn(8)]
			if i%3 == 0 {
				suffix += string([]byte{byte(r.Intn(256))})
			}
			cases = append(cases, scenario{0, base + suffix})
		}
	}
	windows := 0
	if filepath.Separator == '\\' {
		windows = 1
	}
	var rows, fixture bytes.Buffer
	cwd, err := os.Getwd()
	must(err)
	must(os.Chdir(root))
	for _, c := range cases {
		fmt.Fprintf(&rows, "%d %d", windows, c.depth)
		for i := range len(c.pattern) {
			fmt.Fprintf(&rows, " %d", c.pattern[i])
		}
		fmt.Fprintln(&rows)
		fmt.Fprintf(&fixture, "%d x%s %s\n", c.depth, hex.EncodeToString([]byte(c.pattern)), expected(c))
	}
	must(os.Chdir(cwd))
	wantRows := strings.Split(strings.TrimSuffix(fixture.String(), "\n"), "\n")
	for i, row := range check.LeanRows("godiff-path-filepath-glob-with-limit", len(cases), &rows) {
		want := strings.Fields(wantRows[i])[2]
		if row != want {
			log.Fatalf("globWithLimit(%q,%d): Lean %q, Go oracle %q", cases[i].pattern, cases[i].depth, row, want)
		}
	}
	// Keep the fixture outside the tree searched by the native tests.
	out, err := os.CreateTemp("", "golean-glob-fixtures-*")
	must(err)
	path := out.Name()
	defer os.Remove(path)
	_, err = out.Write(fixture.Bytes())
	must(err)
	must(out.Close())
	must(os.Setenv("GOLEAN_GLOB_LIMIT_CASES", path))
	must(os.Setenv("GOLEAN_GLOB_LIMIT_DIR", root))
	const data = "difftest/path/filepath/glob-with-limit/testdata/"
	must(check.TestInPackage("path/filepath", "TestGoLeanGlobWithLimit", data+"bridge_test.go", data+"glob_limit_test.go"))
	fmt.Printf("%d boundaries and %d seeded globWithLimit observations agree.\n", boundaries, len(cases)-boundaries)
}

func expected(c scenario) string {
	if c.depth == 10000 || (c.depth == 9999 && c.pattern == "*/*.txt") {
		return "bad"
	}
	paths, err := filepath.Glob(c.pattern)
	if err != nil {
		return "bad"
	}
	if paths == nil {
		return "nil"
	}
	for i, p := range paths {
		fields := make([]string, len(p))
		for j := range len(p) {
			fields[j] = strconv.Itoa(int(p[j]))
		}
		paths[i] = strings.Join(fields, ",")
	}
	return strconv.Itoa(len(paths)) + ":" + strings.Join(paths, ";")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
