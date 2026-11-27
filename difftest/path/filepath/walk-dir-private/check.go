package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	windows := 0
	if filepath.Separator == '\\' {
		windows = 1
	}
	// Exhaust the finite action table once. Each sequence supplies all
	// six callback results, even when an earlier callback stops traversal.
	cases := make([]string, 0, 4096+60000)
	for i := range 4096 {
		cases = append(cases, sequence(i))
	}
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for range 20000 {
			cases = append(cases, sequence(r.Intn(4096)))
		}
	}
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&rows, "%d 0", windows)
		for _, v := range c {
			fmt.Fprintf(&rows, " %c", v)
		}
		fmt.Fprintln(&rows, " 0")
	}
	got := check.LeanRows("godiff-path-filepath-walk-dir-private", len(cases), &rows)
	file, err := os.CreateTemp("", "golean-walk-dir-*")
	must(err)
	defer os.Remove(file.Name())
	var fixture bytes.Buffer
	for i, row := range got {
		fmt.Fprintf(&fixture, "%s %s\n", cases[i], strings.ReplaceAll(row, "\\", "/"))
	}
	_, err = file.Write(fixture.Bytes())
	must(err)
	must(file.Close())
	must(os.Setenv("GOLEAN_WALK_DIR_CASES", file.Name()))
	const data = "difftest/path/filepath/walk-dir-private/testdata/"
	must(check.TestInPackage("path/filepath", "TestGoLeanPrivateWalkDir", data+"bridge_test.go", data+"walk_dir_test.go"))
	var fixed bytes.Buffer
	for retry := range 4 {
		fmt.Fprintf(&fixed, "0 1 0 0 0 0 0 0 %d\n", retry)
	}
	wants := []string{
		"ok|root:d,root:d:7,root/a:f,root/d:d,root/d/b:f,root/d/c:f,root/z:f",
		"ok|root:d,root:d:7", "skipAll|root:d,root:d:7", "9|root:d,root:d:7",
	}
	for i, row := range check.LeanRows("godiff-path-filepath-walk-dir-private", len(wants), &fixed) {
		if row != wants[i] {
			log.Fatalf("partial read retry %d: got %s, want %s", i, row, wants[i])
		}
	}
	fmt.Println("walkDir: 4,096 exhaustive callback maps, 60,000 seeded repeats and four partial-read model fixtures passed.")
}

func sequence(value int) string {
	var actions [6]byte
	for i := range actions {
		actions[i] = byte('0' + value%4)
		value /= 4
	}
	return string(actions[:])
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
