package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const executor = "godiff-path-filepath-has-prefix"

type paths struct{ path, prefix string }

func main() {
	fixed := []paths{
		{"", ""}, {"", "x"}, {"abc", "ab"}, {"a/b", "a/"}, {"A/b", "a\\"},
		{"a", "ab"}, {"//HOST/SHARE", "//host"}, {"abc", "AB"}, {"Éx", "é"},
		{"Éx", "\xc3"}, {"Kelvin", "k"}, {"ẞ", "ß"}, {"Σ", "ς"}, {"Σ", "σ"},
		{"\xffa", "\xef\xbf\xbd"}, {"\xc3x", "\xef\xbf\xbdx"},
	}
	edges := append([]paths(nil), fixed...)
	for i := range 256 {
		raw := string([]byte{byte(i)})
		edges = append(edges, paths{raw + "a", raw}, paths{raw, "A"})
	}
	nativeWindows := filepath.Separator == '\\'
	native := func(c paths) string { return strconv.FormatBool(filepath.HasPrefix(c.path, c.prefix)) }
	windows := func(c paths) string {
		ok := bytes.HasPrefix([]byte(c.path), []byte(c.prefix)) ||
			bytes.HasPrefix([]byte(strings.ToLower(c.path)), []byte(strings.ToLower(c.prefix)))
		return strconv.FormatBool(ok)
	}
	check.CompareSamples(executor, edges, sample, func(c paths) string { return encode(c, nativeWindows) },
		native, func(c paths) string {
			if nativeWindows {
				return windows(c)
			}
			return strconv.FormatBool(bytes.HasPrefix([]byte(c.path), []byte(c.prefix)))
		})
	if !nativeWindows {
		check.CompareSamples(executor, edges, sample, func(c paths) string { return encode(c, true) },
			windows, windows)
	}
}

func sample(r *rand.Rand, i int) paths {
	pieces := []string{"", "A", "a", "/", "\\", "K", "ẞ", "ς", "Σ", "é", "e\xcc\x81", "\xff"}
	token := func() string {
		raw := make([]byte, r.Intn(16))
		r.Read(raw)
		return pieces[r.Intn(len(pieces))] + string(raw)
	}
	path, prefix := token(), token()
	switch i % 5 {
	case 0:
		prefix = path[:r.Intn(len(path)+1)]
	case 1:
		prefix = strings.ToLower(path)
	case 2:
		prefix += path
	case 3:
		path = strings.ToUpper(prefix) + path
	}
	return paths{path, prefix}
}

func encode(c paths, windows bool) string {
	profile := 0
	if windows {
		profile = 1
	}
	row := fmt.Appendf(nil, "%d %d", profile, len(c.path))
	for _, raw := range []string{c.path, c.prefix} {
		for i := range len(raw) {
			row = strconv.AppendUint(append(row, ' '), uint64(raw[i]), 10)
		}
	}
	return string(row)
}
