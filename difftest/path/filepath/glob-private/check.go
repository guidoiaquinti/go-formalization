package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type scenario struct {
	state          int
	dir, pattern   string
	initial, names []string
}

func main() {
	boundaries := []scenario{
		{0, "missing", "[", nil, []string{"a"}},
		{1, "file", "[", []string{"old"}, nil},
		{2, "dir", "*", []string{"old"}, []string{"a"}},
		{3, "dir", "[", nil, nil},
		{3, "dir", "[", []string{}, []string{"a"}},
		{3, "dir", "*", []string{"old"}, []string{"z", "ba", "a", "a"}},
		{3, "dir", "a*[", nil, []string{"b"}},
		{3, "dir", "a*[", nil, []string{"a"}},
	}
	for i := range 256 {
		raw := string([]byte{byte(i)})
		boundaries = append(boundaries, scenario{3, "d", "*", nil, []string{"b", raw + "z", raw}})
	}
	check.CompareSamples("godiff-path-filepath-glob-private", boundaries, sample, encode, expected, expected)
	const dir = "difftest/path/filepath/glob-private/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanPrivateGlob", dir+"bridge_test.go", dir+"glob_test.go"); err != nil {
		log.Fatal(err)
	}
}

func sample(r *rand.Rand, i int) scenario {
	patterns := []string{"*", "?", "a*", "[a-z]*", "a*[", "[", "[!x]", "\\*", "", "a?", "[α-ω]", "*a"}
	dirs := []string{"", ".", "a/../d", "a", `C:\d`, `\\host\share`}
	c := scenario{state: r.Intn(6), dir: dirs[r.Intn(len(dirs))], pattern: patterns[r.Intn(len(patterns))]}
	if c.state > 3 {
		c.state = 3
	}
	if i%2 == 0 {
		c.initial = []string{"preexisting"}
	}
	for range r.Intn(9) {
		name := make([]byte, r.Intn(13))
		r.Read(name)
		if i%3 == 0 {
			name = append([]byte("a"), name...)
		}
		c.names = append(c.names, string(name))
	}
	return c
}

func encode(c scenario) string {
	windows := 0
	if filepath.Separator == '\\' {
		windows = 1
	}
	flag := 0
	if c.initial != nil {
		flag = 1
	}
	row := fmt.Appendf(nil, "%d %d %d %d %d", windows, c.state, flag, len(c.initial), len(c.names))
	stringsToEncode := append([]string{c.dir, c.pattern}, c.initial...)
	stringsToEncode = append(stringsToEncode, c.names...)
	for _, s := range stringsToEncode {
		row = fmt.Appendf(row, " %d", len(s))
		for i := range len(s) {
			row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
		}
	}
	return string(row)
}

func expected(c scenario) string {
	out := slices.Clone(c.initial)
	status := "ok"
	if c.state == 3 {
		names := slices.Clone(c.names)
		slices.Sort(names)
		for _, name := range names {
			matched, err := filepath.Match(c.pattern, name)
			if err != nil {
				status = "bad"
				break
			}
			if matched {
				out = append(out, filepath.Join(c.dir, name))
			}
		}
	}
	if out == nil {
		return status + "|nil"
	}
	for i, s := range out {
		b := make([]string, len(s))
		for j := range len(s) {
			b[j] = strconv.Itoa(int(s[j]))
		}
		out[i] = strings.Join(b, ",")
	}
	return fmt.Sprintf("%s|%d:%s", status, len(out), strings.Join(out, ";"))
}
