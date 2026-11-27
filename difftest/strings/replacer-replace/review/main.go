package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type fixture struct {
	mode                                 int
	pairs                                []string
	first, second, wantFirst, wantSecond string
}

func native(c fixture) (out string) {
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(runtime.Error)
			if c.mode != 0 || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			out = "panic"
		}
	}()
	var r *strings.Replacer
	if c.mode == 1 {
		r = new(strings.Replacer)
	}
	if c.mode == 2 {
		r = strings.NewReplacer(c.pairs...)
	}
	return wire(r.Replace(c.first)) + "|" + wire(r.Replace(c.second))
}
func wire(s string) string {
	parts := make([]string, len(s))
	for i := range len(s) {
		parts[i] = strconv.Itoa(int(s[i]))
	}
	return strings.Join(parts, ",")
}
func main() {
	cases := []fixture{
		{0, nil, "", "a", "", ""},
		{1, nil, "\xff\x00", "a", "\xff\x00", "a"},
		{2, nil, "", "\x80", "", "\x80"},
		{2, []string{"aa", "R"}, "aaa", "aa", "Ra", "R"},
		{2, []string{"\xff", "A", "\xff", "B"}, "\xff", "z\xff", "A", "zA"},
		{2, []string{"a", "", "b", "wide", "a", "ignored"}, "abba", "b\xffa", "widewide", "wide\xff"},
		{2, []string{"", "X", "\xff", "Y"}, "\xff", "é", "XYX", "X\xc3X\xa9X"},
		{2, []string{"\xff", "Y", "", "X"}, "\xff", "", "YX", "X"},
		{2, []string{"a", "b", "b", "c"}, "ab", "ba", "bc", "cb"},
		{2, []string{"a", "short", "ab", "long"}, "ab", "a", "shortb", "short"},
	}
	var model bytes.Buffer
	var wants []string
	add := func(s string) {
		for i := range len(s) {
			fmt.Fprintf(&model, " %d", s[i])
		}
	}
	for i, c := range cases {
		fmt.Fprintf(&model, "%d %d", c.mode, len(c.pairs)/2)
		for j := 0; j < len(c.pairs); j += 2 {
			key, value := c.pairs[j], c.pairs[j+1]
			fmt.Fprintf(&model, " %d %d", len(key), len(value))
			add(key)
			add(value)
		}
		fmt.Fprintf(&model, " %d", len(c.first))
		add(c.first)
		add(c.second)
		fmt.Fprintln(&model)
		want := wire(c.wantFirst) + "|" + wire(c.wantSecond)
		if c.mode == 0 {
			want = "panic"
		}
		if got := native(c); got != want {
			panic(fmt.Sprintf("native row %d got %q want %q", i, got, want))
		}
		wants = append(wants, want)
	}
	binary := "lean/.lake/build/bin/godiff-strings-replacer-replace"
	if len(os.Args) == 2 {
		binary = os.Args[1]
	}
	cmd := exec.Command(binary)
	cmd.Stdin = &model
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(wants) {
		panic(fmt.Sprintf("Lean produced %d rows", len(rows)))
	}
	for i, want := range wants {
		if rows[i] != want {
			panic(fmt.Sprintf("Lean row %d got %q want %q", i, rows[i], want))
		}
	}
	if strings.ReplaceAll(cases[6].second, "", "X") == cases[6].wantSecond {
		panic("rune-boundary insertion survived")
	}
	c := cases[8]
	if strings.ReplaceAll(strings.ReplaceAll(c.first, "a", "b"), "b", "c") == c.wantFirst {
		panic("cascading replacement survived")
	}
	fmt.Printf("Replacer.Replace passed %d independent pinned Go and Lean sequences\n", len(cases))
}
