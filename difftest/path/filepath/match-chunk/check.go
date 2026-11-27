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
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type scenario struct{ pattern, input string }
type output struct {
	rest    string
	ok, bad bool
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []scenario{{"", "abc"}, {"?", "/"}, {"?", "\\"}, {"?", "\xffx"}, {"[a-c]", "bx"},
		{"[^a-c]", "dx"}, {"[/]", "/x"}, {"a[", "bx"}, {"\\", "x"}, {"[z-a]", "zx"},
		{"[�]", "\xffx"}, {"\\?*", "?*x"}}
	for b := range 256 {
		returned := string([]byte{byte(b), 'x'})
		cases = append(cases, scenario{string([]byte{byte(b)}), returned},
			scenario{string([]byte{'[', byte(b), ']'}), returned})
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		r := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			chunk := make([]byte, r.Intn(48))
			r.Read(chunk)
			for j := 0; j < len(chunk); j += 4 {
				chunk[j] = "?[]-\\^"[r.Intn(6)]
			}
			input := make([]byte, r.Intn(48))
			r.Read(input)
			if i%3 == 0 {
				copy(input, chunk)
			}
			cases = append(cases, scenario{string(chunk), string(input)})
		}
	}
	var lean, fixture bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&lean, len(c.pattern))
		for _, word := range []string{c.pattern, c.input} {
			for i := range len(word) {
				fmt.Fprint(&lean, " ", word[i])
			}
		}
		fmt.Fprintln(&lean)
		want := expected(c, runtime.GOOS == "windows")
		fmt.Fprintf(&fixture, "x%s x%s %t %t x%s\n", hex.EncodeToString([]byte(c.pattern)),
			hex.EncodeToString([]byte(c.input)), want.ok, want.bad, hex.EncodeToString([]byte(want.rest)))
	}
	for i, row := range check.LeanRows("godiff-path-filepath-match-chunk", len(cases), &lean) {
		c := cases[i]
		want := wire(expected(c, false)) + "|" + wire(expected(c, true))
		if row != want {
			return fmt.Errorf("pattern %q, input %q: Lean %q, want %q", c.pattern, c.input, row, want)
		}
	}

	dir, err := os.MkdirTemp("", "golean-match-chunk-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_FILEPATH_MATCH_CHUNK_CASES", path); err != nil {
		return err
	}
	const overlay = "difftest/path/filepath/match-chunk/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanMatchChunk", overlay+"bridge_test.go",
		overlay+"match_chunk_test.go"); err != nil {
		return err
	}
	fmt.Printf("path/filepath.matchChunk: %d boundaries and %d seeded rows agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(c scenario, windows bool) output {
	pat, text := c.pattern, c.input
	matches := true
	for len(pat) > 0 {
		first := pat[0]
		pat = pat[1:]
		switch first {
		case '[':
			r, width := utf8.DecodeRuneInString(text)
			member, rest, valid := class(pat, r, windows)
			if !valid {
				return output{bad: true}
			}
			pat = rest
			matches = matches && width > 0 && member
			text = text[width:]
		case '?':
			separator := byte('/')
			if windows {
				separator = '\\'
			}
			_, width := utf8.DecodeRuneInString(text)
			matches = matches && width > 0 && text[0] != separator
			text = text[width:]
		default:
			if first == '\\' && !windows {
				if pat == "" {
					return output{bad: true}
				}
				first, pat = pat[0], pat[1:]
			}
			matches = matches && len(text) > 0 && text[0] == first
			if text != "" {
				text = text[1:]
			}
		}
	}
	if !matches {
		return output{}
	}
	return output{rest: text, ok: true}
}

func class(pat string, r rune, windows bool) (bool, string, bool) {
	negated := strings.HasPrefix(pat, "^")
	if negated {
		pat = pat[1:]
	}
	matched := false
	for {
		low, tail, ok := endpoint(pat, windows)
		if !ok {
			return false, "", false
		}
		high := low
		if tail[0] == '-' {
			high, tail, ok = endpoint(tail[1:], windows)
			if !ok {
				return false, "", false
			}
		}
		matched = matched || (low <= r && r <= high)
		if tail[0] == ']' {
			return matched != negated, tail[1:], true
		}
		pat = tail
	}
}

func endpoint(pat string, windows bool) (rune, string, bool) {
	if pat == "" || pat[0] == '-' || pat[0] == ']' {
		return 0, "", false
	}
	if !windows {
		pat = strings.TrimPrefix(pat, "\\")
	}
	decoded, width := utf8.DecodeRuneInString(pat)
	return decoded, pat[width:], width < len(pat) && !(decoded == utf8.RuneError && width == 1)
}

func wire(result output) string {
	fields := make([]string, len(result.rest))
	for i := range len(result.rest) {
		fields[i] = strconv.Itoa(int(result.rest[i]))
	}
	return fmt.Sprintf("%t:%t:%s", result.ok, result.bad, strings.Join(fields, ","))
}
