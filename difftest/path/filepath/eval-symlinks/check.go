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

const binary = "godiff-path-filepath-eval-symlinks"

func main() {
	fixed := []struct {
		profile, walk, norm, stat               int
		path, walked, normalized, message, want string
	}{
		{0, 0, 9, 7, "in", "walked", "ignored", "ignored", wire("ok:", "walked")},
		{0, 4, 0, 0, "in", "\xff", "ignored", "", wire("err:4:", "\xff")},
		{1, 0, 0, 0, "in", "walked", `C:\Clean`, "", wire("ok:", `C:\Clean`)},
		{1, 4, 0, 0, "in", "walk error", "ignored", "", wire("err:4:", "walk error")},
		{1, 0, 5, 0, "in", "walked", "norm error", "", wire("err:5:", "norm error")},
		{2, 0, 0, 3, "", "", "", "prefix not a directory suffix", wire("ok:", ".")},
		{2, 0, 0, 0, "a/../b", "", "", "", wire("ok:", "b")},
		{2, 0, 0, 3, "a/b", "", "", "prefix not a directory suffix", "notDirectory"},
		{2, 0, 0, 3, "a/b", "", "", "NOT A DIRECTORY", wire("err:3:", "NOT A DIRECTORY")},
	}
	var rows bytes.Buffer
	for _, c := range fixed {
		fmt.Fprintln(&rows, encode(c.profile, c.walk, c.norm, c.stat, c.path, c.walked, c.normalized, c.message))
	}
	for i, row := range check.LeanRows(binary, len(fixed), &rows) {
		if row != fixed[i].want {
			log.Fatalf("fixed profile row %d: Lean %s, want %s", i, row, fixed[i].want)
		}
	}
	fmt.Printf("%d fixed profile outcomes matched.\n", len(fixed))
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		return
	}
	root, err := os.MkdirTemp("", "golean-symlinks-*")
	must(err)
	defer os.RemoveAll(root)
	must(os.Mkdir(filepath.Join(root, "dir"), 0700))
	must(os.WriteFile(filepath.Join(root, "dir/file"), nil, 0600))
	must(os.Symlink("dir", filepath.Join(root, "link")))
	must(os.Symlink("dir/file", filepath.Join(root, "filelink")))
	must(os.Symlink("missing", filepath.Join(root, "broken")))
	must(os.Symlink("cycle", filepath.Join(root, "cycle")))
	must(os.Symlink("link", filepath.Join(root, "chain")))
	must(os.Symlink(filepath.Join(root, "dir"), filepath.Join(root, "absolute")))
	cases := []string{"", ".", "dir", "./dir/file", "link", "link/file", "filelink", "filelink/x", "link/../dir", "missing", "broken", "cycle", "chain/file", "absolute/file"}
	for i := range 256 {
		cases = append(cases, "link/"+string([]byte{byte(i)}))
	}
	boundaries := len(cases)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for i := range 20000 {
			base := []string{"", "./", "dir/", "link/", "filelink/", "missing/", "chain/", "absolute/"}[r.Intn(8)]
			name := []string{"", "file", ".", "..", "//.", "../dir/file"}[r.Intn(6)]
			if i%2 == 0 {
				name += string([]byte{byte(r.Intn(256))})
			}
			cases = append(cases, base+name)
		}
	}
	cwd, err := os.Getwd()
	must(err)
	must(os.Chdir(root))
	canonicalRoot, err := os.Getwd()
	must(err)
	for _, c := range []struct{ path, want string }{
		{"filelink", "dir/file"},
		{"chain/file", "dir/file"},
		{"absolute/file", filepath.Join(canonicalRoot, "dir/file")},
	} {
		got, err := filepath.EvalSymlinks(c.path)
		if err != nil || got != c.want {
			log.Fatalf("EvalSymlinks(%q)=(%q,%v), want %q", c.path, got, err, c.want)
		}
	}
	rows.Reset()
	wants := make([]string, 0, len(cases))
	var fixture bytes.Buffer
	for _, path := range cases {
		value, err := filepath.EvalSymlinks(path)
		code := 0
		if err != nil {
			if value != "" {
				log.Fatalf("native result %q with error %v", value, err)
			}
			code, value = 1, err.Error()
			wants = append(wants, wire("err:1:", value))
		} else {
			wants = append(wants, wire("ok:", value))
		}
		fmt.Fprintln(&rows, encode(0, code, 0, 0, path, value, "", ""))
		fmt.Fprintln(&fixture, hex.EncodeToString([]byte(path)))
	}
	must(os.Chdir(cwd))
	for i, row := range check.LeanRows(binary, len(cases), &rows) {
		if row != wants[i] {
			log.Fatalf("%q: Lean %s Go %s", cases[i], row, wants[i])
		}
	}
	file, err := os.CreateTemp("", "golean-eval-symlink-cases-*")
	must(err)
	defer os.Remove(file.Name())
	_, err = file.Write(fixture.Bytes())
	must(err)
	must(file.Close())
	must(os.Setenv("GOLEAN_PUBLIC_EVAL_SYMLINK_CASES", file.Name()))
	must(os.Setenv("GOLEAN_PUBLIC_EVAL_SYMLINK_TREE", root))
	const data = "difftest/path/filepath/eval-symlinks/testdata/"
	must(check.TestInPackage("path/filepath", "TestGoLeanPublicEvalSymlinks", data+"bridge_test.go", data+"eval_symlinks_test.go"))
	fmt.Printf("%d native boundaries and %d seeded paths passed.\n", boundaries, len(cases)-boundaries)
}

func encode(p, w, n, s int, strs ...string) string {
	row := fmt.Appendf(nil, "%d %d %d %d", p, w, n, s)
	for _, str := range strs {
		row = fmt.Appendf(row, " %d", len(str))
		for i := range len(str) {
			row = strconv.AppendUint(append(row, ' '), uint64(str[i]), 10)
		}
	}
	return string(row)
}
func wire(prefix, raw string) string {
	out := make([]string, len(raw))
	for i := range len(raw) {
		out[i] = strconv.Itoa(int(raw[i]))
	}
	return prefix + strings.Join(out, ",")
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
