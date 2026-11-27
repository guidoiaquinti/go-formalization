package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const binary = "godiff-path-filepath-walk-symlinks"

func main() {
	fixed := []struct {
		profile    int
		path, want string
	}{
		{0, "/base/dir/file", wire("/base/dir/file")}, {0, "abs/file", wire("/base/dir/file")},
		{0, "filelink/", "notDirectory"}, {0, "cycle", "tooManyLinks"},
		{1, ".", wire(".")}, {1, "C:/dir/file", wire(`C:\dir\file`)},
		{1, "C:/dir/file/", "notDirectory"}, {1, "C:/missing", "missing"},
	}
	var rows bytes.Buffer
	for _, f := range fixed {
		fmt.Fprintln(&rows, encode(f.profile, f.path))
	}
	for i, row := range check.LeanRows(binary, len(fixed), &rows) {
		if row != fixed[i].want {
			log.Fatalf("fixed path %q: Lean %s, want %s", fixed[i].path, row, fixed[i].want)
		}
	}
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		return
	}
	parent, err := os.MkdirTemp("", "golean-walk-symlinks-*")
	must(err)
	defer os.RemoveAll(parent)
	root := filepath.Join(parent, "base")
	must(os.MkdirAll(filepath.Join(root, "dir"), 0700))
	must(os.WriteFile(filepath.Join(root, "dir/file"), nil, 0600))
	for _, link := range [][2]string{{"filelink", "dir/file"}, {"link", "dir"}, {"chain", "link"}, {"broken", "missing"}, {"cycle", "cycle"}, {"up", "dir/../dir"}} {
		must(os.Symlink(link[1], filepath.Join(root, link[0])))
	}
	cases := []string{"", ".", "..", "dir", "dir/file", "link", "link/file", "filelink", "filelink/x", "filelink/", "chain/file", "broken", "cycle", "up/file", "dir/../link/file"}
	for i := range 256 {
		cases = append(cases, "dir/"+string([]byte{byte(i)}))
	}
	boundaries := len(cases)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for i := range 20000 {
			base := []string{"", "./", "dir/", "link/", "chain/", "broken/", "filelink/", "up/"}[r.Intn(8)]
			leaf := []string{"", "file", ".", "../dir/file", "missing", "file/"}[r.Intn(6)]
			if i%2 == 0 {
				leaf += string([]byte{byte(r.Intn(256))})
			}
			cases = append(cases, base+leaf)
		}
	}
	rows.Reset()
	for _, s := range cases {
		fmt.Fprintln(&rows, encode(0, s))
	}
	results := check.LeanRows(binary, len(cases), &rows)
	var fixture bytes.Buffer
	for i, p := range cases {
		fmt.Fprintf(&fixture, "x%s %s\n", hex.EncodeToString([]byte(p)), results[i])
	}
	file, err := os.CreateTemp("", "golean-walk-symlinks-fixture-*")
	must(err)
	defer os.Remove(file.Name())
	_, err = file.Write(fixture.Bytes())
	must(err)
	must(file.Close())
	must(os.Setenv("GOLEAN_WALK_SYMLINK_CASES", file.Name()))
	must(os.Setenv("GOLEAN_WALK_SYMLINK_TREE", root))
	const data = "difftest/path/filepath/walk-symlinks/testdata/"
	must(check.TestInPackage("path/filepath", "TestGoLeanWalkSymlinks", data+"bridge_test.go", data+"walk_symlinks_test.go"))
	fmt.Printf("%d native boundaries, %d seeded paths and %d alternate observations passed.\n", boundaries, len(cases)-boundaries, len(fixed))
}

func encode(profile int, path string) string {
	row := strconv.AppendInt(nil, int64(profile), 10)
	for i := range len(path) {
		row = strconv.AppendUint(append(row, ' '), uint64(path[i]), 10)
	}
	return string(row)
}
func wire(path string) string {
	fields := make([]string, len(path))
	for i := range len(path) {
		fields[i] = strconv.Itoa(int(path[i]))
	}
	return "ok:" + strings.Join(fields, ",")
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
