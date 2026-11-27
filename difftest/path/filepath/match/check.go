package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ pattern, name string }
type fixture struct {
	input
	unix, windows string
}

const (
	match  = "true:false"
	miss   = "false:false"
	syntax = "false:true"
	exe    = "godiff-path-filepath-match"
)

func main() {
	fixed := []fixture{
		{input{"", ""}, match, match},
		{input{"*", "a/b"}, miss, match},
		{input{"*", "a\\b"}, match, miss},
		{input{"*/b", "a/b"}, match, match},
		{input{"**?", "é"}, match, match},
		{input{"?", "/"}, miss, match},
		{input{"?", "\\"}, match, miss},
		{input{"\\*", "*"}, match, miss},
		{input{"\\*", "\\x"}, miss, match},
		{input{"a*[", "zzz"}, miss, miss},
		{input{"*[", "zzz"}, syntax, syntax},
		{input{"a[", "zzz"}, syntax, syntax},
		{input{"[", ""}, syntax, syntax},
		{input{"[^x]", "/"}, match, match},
		{input{"[z-a]", "z"}, miss, miss},
		{input{"[�]", "\xff"}, match, match},
		{input{"a*b", "abbb"}, match, match},
		{input{"*a*b", "xaab"}, match, match},
		{input{"*a/b", "foo/a/b"}, miss, match},
		{input{"ab", "a"}, miss, miss},
	}
	var rows bytes.Buffer
	var gold []string
	for _, c := range fixed {
		fmt.Fprintln(&rows, encode(c.input, false))
		fmt.Fprintln(&rows, encode(c.input, true))
		gold = append(gold, c.unix, c.windows)
		want := c.unix
		if runtime.GOOS == "windows" {
			want = c.windows
		}
		if got := native(c.input); got != want {
			log.Fatalf("Go golden %q/%q: %s, want %s", c.pattern, c.name, got, want)
		}
	}
	for i, row := range check.LeanRows(exe, len(gold), &rows) {
		if row != gold[i] {
			c := fixed[i/2]
			log.Fatalf("profile %d golden %q/%q: Lean %s, want %s", i%2, c.pattern, c.name, row, gold[i])
		}
	}
	fmt.Printf("filepath.Match: %d explicit fixtures in both profiles agree.\n", len(fixed))

	var edges []input
	for _, c := range fixed {
		edges = append(edges, c.input)
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		edges = append(edges, input{raw, raw}, input{"*" + raw + "?", "q" + raw + "\xff"})
	}
	check.CompareSamples(exe, edges, sample,
		func(c input) string { return encode(c, runtime.GOOS == "windows") }, native, native)
}

func native(c input) string {
	got, err := filepath.Match(c.pattern, c.name)
	if err != nil && err != filepath.ErrBadPattern {
		log.Fatalf("unexpected Match error %v", err)
	}
	return strconv.FormatBool(got) + ":" + strconv.FormatBool(err != nil)
}

func encode(c input, windows bool) string {
	profile := 0
	if windows {
		profile = 1
	}
	row := fmt.Appendf(nil, "%d %d", profile, len(c.pattern))
	for _, part := range []string{c.pattern, c.name} {
		for j := range len(part) {
			row = strconv.AppendUint(append(row, ' '), uint64(part[j]), 10)
		}
	}
	return string(row)
}

func sample(r *rand.Rand, i int) input {
	switch i % 4 {
	case 0:
		raw := make([]byte, r.Intn(32))
		r.Read(raw)
		return input{string(raw), string(raw)}
	case 1:
		pat := []string{"*a?b", "[é-ê]*?", "[abc]*z*", "a*b[", "?*/*", "\\*?"}
		names := []string{"aaabb", "éxyz", "aczZ", "zaz", "/a/b", "*z", "\xff", "\\"}
		return input{pat[r.Intn(len(pat))], names[r.Intn(len(names))]}
	default:
		pattern, name := make([]byte, r.Intn(24)), make([]byte, r.Intn(24))
		r.Read(pattern)
		r.Read(name)
		for j := range pattern {
			if j%2 == 0 {
				pattern[j] = "*?[]\\-^/ab"[r.Intn(10)]
			}
		}
		return input{string(pattern), string(name)}
	}
}
