package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
	"github.com/guidoiaquinti/go-formalization/difftest/path/filepath/internal/absnative"
)

const executor = "godiff-path-filepath-abs"

func main() {
	cwd, err := os.Getwd()
	must(err)
	profile := map[string]int{"windows": 1, "plan9": 2}[runtime.GOOS]
	fixed := []struct {
		profile                  int
		failure                  bool
		path, query, value, want string
	}{
		{0, true, "/a/../b", "", "unavailable", wire("ok:", "/b")},
		{0, true, "", "", "unavailable", wire("err:", "unavailable")},
		{0, false, "../a", "", "/cwd/inner", wire("ok:", "/cwd/a")},
		{2, true, "#dev/./x", "", "unused", wire("ok:", "#dev/x")},
		{2, false, "a", "", "#dev/dir", wire("ok:", "#dev/dir/a")},
		{1, false, "", ".", `C:\cwd\.`, wire("ok:", `C:\cwd`)},
		{1, true, "", ".", "invalid", wire("err:", "invalid")},
		{1, false, "a", "a", `C:\cwd\..\a`, wire("ok:", `C:\a`)},
		{1, false, `D:\already`, `D:\already`, `D:\x\..\already`, wire("ok:", `D:\already`)},
		{1, true, "a", "a", "invalid", wire("err:", "invalid")},
	}
	var input bytes.Buffer
	edges := make([]string, 0, len(fixed)+512)
	for _, f := range fixed {
		fmt.Fprintln(&input, encode(f.profile, f.failure, f.path, f.query, f.value))
		edges = append(edges, f.path)
	}
	for i, row := range check.LeanRows(executor, len(fixed), &input) {
		if row != fixed[i].want {
			log.Fatalf("fixed case %d: Lean %q, want %q", i, row, fixed[i].want)
		}
	}
	fmt.Printf("filepath.Abs: %d fixed profile and error-boundary cases passed.\n", len(fixed))
	for i := range 256 {
		b := string([]byte{byte(i)})
		edges = append(edges, "a/"+b, "/a/"+b)
	}
	outcome := func(s string) (query, value string, failure bool, result string) {
		if profile == 1 {
			query = s
			if query == "" {
				query = "."
			}
			resolved, err := absnative.FullPath(query)
			if err != nil {
				return query, err.Error(), true, "err:" + err.Error()
			}
			return query, resolved, false, "ok:" + filepath.Clean(resolved)
		}
		if filepath.IsAbs(s) {
			return "", cwd, false, "ok:" + filepath.Clean(s)
		}
		return "", cwd, false, "ok:" + filepath.Join(cwd, s)
	}
	want := func(s string) string {
		_, _, _, result := outcome(s)
		separator := strings.IndexByte(result, ':') + 1
		return wire(result[:separator], result[separator:])
	}
	native := func(s string) string {
		result, err := filepath.Abs(s)
		if err != nil {
			if result != "" {
				log.Fatalf("Abs(%q) returned %q with error %v", s, result, err)
			}
			return wire("err:", err.Error())
		}
		return wire("ok:", result)
	}
	check.CompareSamples(executor, edges, sample, func(s string) string {
		query, value, failed, _ := outcome(s)
		return encode(profile, failed, s, query, value)
	}, native, want)
}

func sample(r *rand.Rand, i int) string {
	raw := make([]byte, r.Intn(64))
	r.Read(raw)
	for j := 0; j < len(raw); j += 5 {
		raw[j] = "/\\.#:"[r.Intn(5)]
	}
	base := []string{"", "//", "/..", "a/../../", "#dev/.", `C:x\`, "a/./b/../"}[i%7]
	return base + string(raw)
}

func encode(profile int, failure bool, path, query, value string) string {
	flag := 0
	if failure {
		flag = 1
	}
	out := fmt.Appendf(nil, "%d %d %d %d", profile, flag, len(path), len(query))
	for _, raw := range []string{path, query, value} {
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
