package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func result(status string, paths ...string) string {
	if paths == nil {
		return status + "|nil"
	}
	raw := make([]string, len(paths))
	for i, path := range paths {
		fields := make([]string, len(path))
		for j := range len(path) {
			fields[j] = strconv.Itoa(int(path[j]))
		}
		raw[i] = strings.Join(fields, ",")
	}
	return fmt.Sprintf("%s|%d:%s", status, len(paths), strings.Join(raw, ";"))
}

func main() {
	cases := []struct {
		state          int
		dir, pattern   string
		initial, names []string
		want           string
	}{
		{0, "x", "[", nil, nil, result("ok")},
		{1, "x", "[", []string{"keep"}, nil, result("ok", "keep")},
		{2, "x", "[", []string{}, nil, "ok|0:"},
		{3, "x", "[", nil, nil, result("ok")},
		{3, "x", "[", []string{"keep"}, []string{"z"}, result("bad", "keep")},
		{3, "x", "a*[", []string{"keep"}, []string{"b", "a"}, result("bad", "keep")},
		{3, "x", "a*[", nil, []string{"b"}, result("ok")},
		{3, "x/./q/..", "*", []string{"keep"}, []string{"z", "aa", "a", "\xff"}, result("ok", "keep", "x/a", "x/aa", "x/z", "x/\xff")},
	}
	var in bytes.Buffer
	for _, c := range cases {
		present := 0
		if c.initial != nil {
			present = 1
		}
		fmt.Fprintf(&in, "0 %d %d %d %d", c.state, present, len(c.initial), len(c.names))
		fields := append([]string{c.dir, c.pattern}, c.initial...)
		fields = append(fields, c.names...)
		for _, field := range fields {
			fmt.Fprintf(&in, " %d", len(field))
			for i := range len(field) {
				fmt.Fprintf(&in, " %d", field[i])
			}
		}
		fmt.Fprintln(&in)
	}
	for i, got := range check.LeanRows("godiff-path-filepath-glob-private", len(cases), &in) {
		if got != cases[i].want {
			log.Fatalf("row %d: got %q want %q", i, got, cases[i].want)
		}
	}
	const data = "difftest/path/filepath/glob-private/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanPrivateGlob", data+"bridge_test.go", data+"glob_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("private glob passed %d fixed model outcomes and seven pinned native fixtures\n", len(cases))
}
