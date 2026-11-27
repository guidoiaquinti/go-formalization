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

const executor = "godiff-path-filepath-unix-abs"

func main() {
	cwd, err := os.Getwd()
	must(err)
	profile := map[string]int{"windows": 1, "plan9": 2}[runtime.GOOS]
	fixed := []struct {
		profile         int
		failure         bool
		path, cwd, want string
	}{
		{0, true, "/a/../b", "cwd unavailable", wire("ok:", "/b")},
		{0, true, "a", "cwd unavailable", wire("err:", "cwd unavailable")},
		{0, true, "", "\xff\x00", wire("err:", "\xff\x00")},
		{0, false, "", "/cwd", wire("ok:", "/cwd")},
		{0, false, "../a", "/cwd/inner", wire("ok:", "/cwd/a")},
		{0, false, "#dev", "/cwd", wire("ok:", "/cwd/#dev")},
		{2, true, "#dev/./x", "unused", wire("ok:", "#dev/x")},
		{2, false, "a", "#dev/dir", wire("ok:", "#dev/dir/a")},
		{1, true, `C:/a/../b`, "unused", wire("ok:", `C:\b`)},
		{1, false, `\child`, `C:\cwd`, wire("ok:", `C:\cwd\child`)},
		{1, false, `leaf`, `C:\cwd`, wire("ok:", `C:\cwd\leaf`)},
		{1, true, `\child`, "cwd unavailable", wire("err:", "cwd unavailable")},
	}
	var input bytes.Buffer
	edges := make([]string, 0, len(fixed)+512)
	for _, f := range fixed {
		fmt.Fprintln(&input, encode(f.profile, f.failure, f.path, f.cwd))
		edges = append(edges, f.path)
	}
	for i, row := range check.LeanRows(executor, len(fixed), &input) {
		if row != fixed[i].want {
			log.Fatalf("fixed case %d: Lean %q, want %q", i, row, fixed[i].want)
		}
	}
	fmt.Printf("filepath.unixAbs: %d fixed profile and error-boundary cases passed.\n", len(fixed))
	for i := range 256 {
		b := string([]byte{byte(i)})
		edges = append(edges, "a/"+b, "/a/"+b)
	}
	rawResult := func(s string) string {
		if filepath.IsAbs(s) {
			return filepath.Clean(s)
		}
		return filepath.Join(cwd, s)
	}
	want := func(s string) string { return wire("ok:", rawResult(s)) }
	check.CompareSamples(executor, edges, sample,
		func(s string) string { return encode(profile, false, s, cwd) }, want, want)

	file, err := os.CreateTemp("", "golean-unix-abs-*")
	must(err)
	defer os.Remove(file.Name())
	var fixture bytes.Buffer
	write := func(s string) {
		fmt.Fprintf(&fixture, "x%s\nx%s\n", hex.EncodeToString([]byte(s)), hex.EncodeToString([]byte(rawResult(s))))
	}
	fmt.Fprintf(&fixture, "x%s\n", hex.EncodeToString([]byte(cwd)))
	for _, s := range edges {
		write(s)
	}
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for i := range 20000 {
			write(sample(r, i))
		}
	}
	_, err = file.Write(fixture.Bytes())
	must(err)
	must(file.Close())
	must(os.Setenv("GOLEAN_UNIX_ABS_CASES", file.Name()))
	const data = "difftest/path/filepath/unix-abs/testdata/"
	must(check.TestInPackage("path/filepath", "TestGoLeanUnixAbs", data+"bridge_test.go", data+"unix_abs_test.go"))
}

func sample(r *rand.Rand, i int) string {
	raw := make([]byte, r.Intn(64))
	r.Read(raw)
	for j := 0; j < len(raw); j += 5 {
		raw[j] = "/\\.#:"[r.Intn(5)]
	}
	base := []string{"", "", "/", "../", "#dev/", `C:\`, "a/./b/../"}[i%7]
	return base + string(raw)
}

func encode(profile int, failure bool, path, cwd string) string {
	flag := 0
	if failure {
		flag = 1
	}
	out := fmt.Appendf(nil, "%d %d %d", profile, flag, len(path))
	for _, raw := range []string{path, cwd} {
		for i := range len(raw) {
			out = strconv.AppendUint(append(out, ' '), uint64(raw[i]), 10)
		}
	}
	return string(out)
}

func wire(kind, path string) string {
	fields := make([]string, len(path))
	for i := range len(path) {
		fields[i] = strconv.Itoa(int(path[i]))
	}
	return kind + strings.Join(fields, ",")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
