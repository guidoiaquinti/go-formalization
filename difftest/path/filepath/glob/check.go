package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	root, err := os.MkdirTemp("", "golean-glob-limit-*")
	must(err)
	defer os.RemoveAll(root)
	for _, dir := range []string{"alpha", "beta", "alpha/tree"} {
		must(os.Mkdir(filepath.Join(root, dir), 0700))
	}
	for _, name := range []string{"plain", "alpha/a.txt", "alpha/z.go", "beta/b.txt", "alpha/tree/deep.txt"} {
		must(os.WriteFile(filepath.Join(root, name), nil, 0600))
	}
	cases := []string{
		"", "plain", "./plain", "missing", "alpha", ".",
		"[", "*/*.txt", "*/[ab].txt", "a*/?.txt", "a*/*/*.txt",
		"a*/t*/*.txt", "*/*/*.txt", "alpha/tree/deep.txt", "*/tree/deep.txt", "missing/../plain", "alpha/a.txt/../plain", "missing/../*", "plain/../*", "alpha/../*",
	}
	for i := range 256 {
		cases = append(cases, "alpha/"+string([]byte{byte(i)}))
	}
	boundaries := len(cases)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for i := range 20000 {
			suffix := []string{"*", "*.txt", "a*", "?", "[ab].txt", "plain", "missing", "[", "a.txt", "z.go", "b.txt"}[r.Intn(11)]
			base := []string{"", "./", "alpha/", "beta/", "a*/", "*/", "missing/", "a*/t*/", "*/tree/", "*/*/"}[r.Intn(10)]
			if i%3 == 0 {
				raw := make([]byte, 1+r.Intn(5))
				r.Read(raw)
				suffix += string(raw)
			}
			cases = append(cases, base+suffix)
		}
	}
	windows := 0
	if filepath.Separator == '\\' {
		windows = 1
	}
	var rows bytes.Buffer
	var wants []string
	cwd, err := os.Getwd()
	must(err)
	must(os.Chdir(root))
	for _, c := range cases {
		fmt.Fprintf(&rows, "%d 0", windows)
		for i := range len(c) {
			fmt.Fprintf(&rows, " %d", c[i])
		}
		fmt.Fprintln(&rows)
		wants = append(wants, expected(c))
	}
	must(os.Chdir(cwd))
	for i, row := range check.LeanRows("godiff-path-filepath-glob", len(cases), &rows) {
		if row != wants[i] {
			log.Fatalf("Glob(%q): Lean %q, Go %q", cases[i], row, wants[i])
		}
	}
	fmt.Printf("filepath.Glob: %d boundaries and %d seeded observations agree.\n", boundaries, len(cases)-boundaries)
}

func expected(c string) string {
	paths, err := filepath.Glob(c)
	if err != nil {
		if err != filepath.ErrBadPattern || paths != nil {
			log.Fatalf("unexpected Glob result %q, %v", paths, err)
		}
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
