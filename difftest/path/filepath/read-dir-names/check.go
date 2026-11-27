package main

import (
	"fmt"
	"log"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type observation struct {
	stage, failure, close int
	names                 []string
}

func main() {
	edges := []observation{
		{0, 1, 0, nil}, {1, 0, 0, nil}, {1, 2, 8, []string{"stale"}}, {1, 0, 8, []string{"b", "a", "aa", "a", "\xff", "\x00"}},
	}
	for i := range 256 {
		raw := string([]byte{byte(i)})
		edges = append(edges, observation{1, 0, 0, []string{raw + "b", raw, "a"}})
	}
	check.CompareSamples("godiff-path-filepath-read-dir-names", edges, sample, encode, expected, expected)
	const dir = "difftest/path/filepath/read-dir-names/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanReadDirNames", dir+"bridge_test.go", dir+"read_dir_names_test.go"); err != nil {
		log.Fatal(err)
	}
}

func sample(r *rand.Rand, i int) observation {
	o := observation{stage: 1, close: r.Intn(3)}
	switch i % 5 {
	case 0:
		o.stage, o.failure = 0, 1+r.Intn(9)
	case 1:
		o.failure = 1 + r.Intn(9)
	}
	for range r.Intn(15) {
		raw := make([]byte, r.Intn(16))
		r.Read(raw)
		o.names = append(o.names, string(raw))
	}
	return o
}

func encode(o observation) string {
	row := fmt.Appendf(nil, "%d %d %d %d", o.stage, o.failure, o.close, len(o.names))
	for _, name := range o.names {
		row = fmt.Appendf(row, " %d", len(name))
		for i := range len(name) {
			row = strconv.AppendUint(append(row, ' '), uint64(name[i]), 10)
		}
	}
	return string(row)
}

func expected(o observation) string {
	if o.stage == 0 {
		return fmt.Sprintf("1|nil|%d", o.failure)
	}
	if o.failure != 0 {
		return fmt.Sprintf("3|nil|%d", o.failure)
	}
	names := slices.Clone(o.names)
	slices.Sort(names)
	for i, s := range names {
		bytes := make([]string, len(s))
		for j := range len(s) {
			bytes[j] = strconv.Itoa(int(s[j]))
		}
		names[i] = strings.Join(bytes, ",")
	}
	return fmt.Sprintf("3|%d|%s", len(names), strings.Join(names, ";"))
}
