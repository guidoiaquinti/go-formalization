package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type fixture struct {
	pattern string
	want    []string
	bad     bool
}

func main() {
	cases := []fixture{
		{pattern: ""},
		{pattern: "plain", want: []string{"plain"}},
		{pattern: "./alpha//a.txt", want: []string{"./alpha//a.txt"}},
		{pattern: "alpha/../plain", want: []string{"alpha/../plain"}},
		{pattern: "*/*.txt", want: []string{filepath.Join("alpha", "a.txt"), filepath.Join("beta", "b.txt")}},
		{pattern: "?lpha/t*/*.txt", want: []string{filepath.Join("alpha", "tree", "deep.txt")}},
		{pattern: "*/*/d??p.[t]xt", want: []string{filepath.Join("alpha", "tree", "deep.txt")}},
		{pattern: "[ab]*/*.go", want: []string{filepath.Join("alpha", "z.go")}},
		{pattern: "missing/*"},
		{pattern: "plain/*"},
		{pattern: "*/tree/[", bad: true},
		{pattern: "alpha/\x00"},
	}
	if filepath.Separator == '/' {
		cases = append(cases,
			fixture{pattern: "missing/../plain"},
			fixture{pattern: "plain/../plain"},
			fixture{pattern: "missing/../*"},
			fixture{pattern: "plain/../*"},
			fixture{pattern: "alpha/../*", want: []string{"alpha", "beta", "plain"}},
		)
	}
	var rows bytes.Buffer
	profile := 0
	if filepath.Separator == '\\' {
		profile = 1
	}
	for _, c := range cases {
		fmt.Fprintf(&rows, "%d 0", profile)
		for _, b := range []byte(c.pattern) {
			fmt.Fprintf(&rows, " %d", b)
		}
		fmt.Fprintln(&rows)
	}
	checkNative(cases)
	cmd := exec.Command("lean/.lake/build/bin/godiff-path-filepath-glob")
	cmd.Stdin = &rows
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	results := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(results) != len(cases) {
		log.Fatalf("unexpected row count %d", len(results))
	}
	for i, c := range cases {
		if want := c.wire(); results[i] != want {
			log.Fatalf("fixture %d: Lean=%q want %q", i, results[i], want)
		}
	}

	rawLoss, noRecursion, nilLoss := 0, 0, 0
	for _, c := range cases {
		if len(c.want) == 1 && c.want[0] == c.pattern && filepath.Clean(c.pattern) != c.pattern {
			rawLoss++
		}
		if len(c.want) != 0 && strings.Count(c.pattern, "/") == 2 && strings.ContainsAny(strings.SplitN(c.pattern, "/", 2)[0], "*?") {
			noRecursion++
		}
		if c.want == nil && !c.bad {
			nilLoss++
		}
	}
	wantNilLoss := 4
	if filepath.Separator == '/' {
		wantNilLoss += 4
	}
	if rawLoss != 2 || noRecursion != 2 || nilLoss != wantNilLoss {
		log.Fatalf("mutation counts: raw=%d, recursion=%d, nilness=%d", rawLoss, noRecursion, nilLoss)
	}
	fmt.Printf("Glob reviewer control: %d native and Lean fixtures passed; cleanup, one-level and nonnil-empty alternatives fail %d/%d/%d.\n",
		len(cases), rawLoss, noRecursion, nilLoss)
}

func checkNative(cases []fixture) {
	root, err := os.MkdirTemp("", "golean-glob-public-control-")
	must(err)
	defer os.RemoveAll(root)
	for _, dir := range []string{"alpha", "beta", "alpha/tree"} {
		must(os.Mkdir(filepath.Join(root, dir), 0700))
	}
	for _, name := range []string{"alpha/a.txt", "alpha/z.go", "alpha/tree/deep.txt", "beta/b.txt", "plain"} {
		must(os.WriteFile(filepath.Join(root, name), nil, 0600))
	}
	cwd, err := os.Getwd()
	must(err)
	must(os.Chdir(root))
	defer func() { must(os.Chdir(cwd)) }()
	for _, c := range cases {
		got, err := filepath.Glob(c.pattern)
		if c.bad {
			if got != nil || err != filepath.ErrBadPattern {
				log.Fatalf("Glob(%q)=%q, %v, want nil and ErrBadPattern", c.pattern, got, err)
			}
		} else if err != nil || !reflect.DeepEqual(got, c.want) {
			log.Fatalf("Glob(%q)=%q, %v, want %q", c.pattern, got, err, c.want)
		}
	}
}

func (c fixture) wire() string {
	if c.bad {
		return "bad"
	}
	if c.want == nil {
		return "nil"
	}
	parts := make([]string, len(c.want))
	for i, p := range c.want {
		fields := make([]string, len(p))
		for j := range len(p) {
			fields[j] = strconv.Itoa(int(p[j]))
		}
		parts[i] = strings.Join(fields, ",")
	}
	return strconv.Itoa(len(parts)) + ":" + strings.Join(parts, ";")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
