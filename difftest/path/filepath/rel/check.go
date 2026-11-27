package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const executor = "godiff-path-filepath-rel"

type paths struct{ base, target string }

func main() {
	windows := filepath.Separator == '\\'
	fixed := []struct {
		paths
		result string
		fails  bool
	}{
		{paths{"", ""}, ".", false},
		{paths{"C:/a", "c:\\A"}, ".", false},
		{paths{"C:/a", "d:\\A"}, "", true},
		{paths{"C:base", "C:target"}, "..\\target", false},
		{paths{"C:base", "C:/target"}, "", true},
		{paths{"\\\\host\\share", "//HOST/share/file"}, "file", false},
		{paths{"\\\\host\\share\\a", "//HOST/share"}, "", true},
		{paths{"\\\\host\\share", "//host/other/a"}, "", true},
		{paths{"\\a\\b", "/A/c"}, "..\\c", false},
		{paths{"a/..", "b"}, "b", false},
		{paths{"a", "."}, "..", false},
		{paths{"../a", "a"}, "", true},
		{paths{"../a", "../b"}, "..\\b", false},
		{paths{"\xff", "\xef\xbf\xbd"}, ".", false},
		{paths{"x/K", "X/k/z"}, "z", false},
		{paths{"x/a", "x/a/\x00"}, "\x00", false},
	}
	var rows bytes.Buffer
	edges := make([]paths, 0, len(fixed)+512)
	for _, c := range fixed {
		fmt.Fprintln(&rows, encode(c.paths, true))
		edges = append(edges, c.paths)
	}
	for i, row := range check.LeanRows(executor, len(fixed), &rows) {
		c := fixed[i]
		want := wire("ok:", c.result)
		if c.fails {
			want = wire("err:", "Rel: can't make "+c.target+" relative to "+c.base)
		}
		if row != want || (windows && native(c.paths) != want) {
			log.Fatalf("Windows Rel(%q,%q): Lean %s, want %s", c.base, c.target, row, want)
		}
	}
	fmt.Printf("filepath.Rel: %d fixed Windows profile cases agree.\n", len(fixed))
	stalled := paths{`\\host\share`, `\\host\share\`}
	if got := check.LeanRows(executor, 1, strings.NewReader(encode(stalled, true)+"\n"))[0]; got != "diverges" {
		log.Fatalf("Windows UNC cursor stall: Lean %s, want diverges", got)
	}
	edges = append(edges, stalled)
	for i := range 256 {
		raw := string([]byte{byte(i)})
		edges = append(edges, paths{"a/" + raw, "a/x/" + raw}, paths{raw, "/" + raw})
	}
	check.CompareSamples(executor, edges, sample, func(c paths) string { return encode(c, windows) },
		native, native)
}

func native(c paths) string {
	// The pinned Windows cursor loop cannot return after consuming identical
	// bodies reached through UNC-root normalization. The reviewer control
	// exercises this invocation in a bounded child process on Windows.
	if filepath.Separator == '\\' && stallsOnWindows(c) {
		return "diverges"
	}
	result, err := filepath.Rel(c.base, c.target)
	if err != nil {
		if result != "" {
			log.Fatalf("Rel(%q,%q) returned %q with error %q", c.base, c.target, result, err)
		}
		return wire("err:", err.Error())
	}
	return wire("ok:", result)
}

func stallsOnWindows(c paths) bool {
	base, target := filepath.Clean(c.base), filepath.Clean(c.target)
	if strings.EqualFold(base, target) {
		return false
	}
	baseVolume, targetVolume := filepath.VolumeName(c.base), filepath.VolumeName(c.target)
	return len(baseVolume) > 2 && strings.EqualFold(baseVolume, targetVolume) &&
		len(base) == len(baseVolume) && target[len(targetVolume):] == "\\"
}

func sample(r *rand.Rand, i int) paths {
	fragments := []string{"", ".", "..", "a", "A", "?", "??", "\xff", "\x00", "c:", "K"}
	part := func() string {
		var out strings.Builder
		out.WriteString([]string{"", "", "/", "\\", "c:", "d:/", "//host/share"}[r.Intn(7)])
		for range r.Intn(8) {
			if r.Intn(4) == 0 {
				raw := make([]byte, r.Intn(12))
				r.Read(raw)
				out.Write(raw)
			} else {
				out.WriteString(fragments[r.Intn(len(fragments))])
			}
			out.WriteByte("/\\"[r.Intn(2)])
		}
		return out.String()
	}
	left, right := part(), part()
	switch i % 4 {
	case 0:
		right = left + right
	case 1:
		right = left
	case 2:
		prefix := part()
		left, right = prefix+left, prefix+right
	}
	return paths{left, right}
}

func encode(c paths, windows bool) string {
	profile := 0
	if windows {
		profile = 1
	}
	row := fmt.Appendf(nil, "%d %d", profile, len(c.base))
	for _, input := range []string{c.base, c.target} {
		for i := range len(input) {
			row = strconv.AppendUint(append(row, ' '), uint64(input[i]), 10)
		}
	}
	return string(row)
}

func wire(prefix, raw string) string {
	fields := make([]string, len(raw))
	for i := range len(raw) {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return prefix + strings.Join(fields, ",")
}
