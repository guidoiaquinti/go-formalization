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

type output struct {
	rune rune
	rest string
	bad  bool
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []string{"", "-", "]", "\\", "\\-x", "\\]x", "a", "abc", "éx", "🔧x", "\xef\xbf\xbdx", "\xe2\x82x"}
	for b := range 256 {
		cases = append(cases, string([]byte{byte(b)}), string([]byte{byte(b), 'x'}),
			string([]byte{'\\', byte(b), 0xff}))
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		r := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			raw := make([]byte, r.Intn(64))
			r.Read(raw)
			if i%4 == 0 && len(raw) > 0 {
				raw[0] = '\\'
			}
			cases = append(cases, string(raw))
		}
	}
	var lean, fixture bytes.Buffer
	for _, chunk := range cases {
		fmt.Fprint(&lean, 0)
		for i := range len(chunk) {
			fmt.Fprint(&lean, " ", chunk[i])
		}
		fmt.Fprintln(&lean)
		want := expected(chunk, runtime.GOOS == "windows")
		fmt.Fprintf(&fixture, "x%s %d %t x%s\n", hex.EncodeToString([]byte(chunk)),
			want.rune, want.bad, hex.EncodeToString([]byte(want.rest)))
	}
	for i, row := range check.LeanRows("godiff-path-filepath-get-esc", len(cases), &lean) {
		chunk := cases[i]
		want := wire(expected(chunk, false)) + "|" + wire(expected(chunk, true))
		if row != want {
			return fmt.Errorf("chunk %q: Lean %q, want %q", chunk, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-get-esc-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_FILEPATH_GET_ESC_CASES", path); err != nil {
		return err
	}
	const overlay = "difftest/path/filepath/get-esc/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanGetEsc", overlay+"bridge_test.go",
		overlay+"get_esc_test.go"); err != nil {
		return err
	}
	fmt.Printf("path/filepath.getEsc: %d boundaries and %d seeded rows agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(chunk string, windows bool) output {
	if chunk == "" || strings.HasPrefix(chunk, "-") || strings.HasPrefix(chunk, "]") {
		return output{bad: true}
	}
	if !windows {
		chunk = strings.TrimPrefix(chunk, "\\")
	}
	first, end := rune(0), len(chunk)
	for index, decoded := range chunk {
		if index > 0 {
			end = index
			break
		}
		first = decoded
	}
	return output{first, chunk[end:], end == len(chunk) || first == utf8.RuneError && end == 1}
}

func wire(result output) string {
	fields := make([]string, len(result.rest))
	for i := range len(result.rest) {
		fields[i] = strconv.Itoa(int(result.rest[i]))
	}
	return fmt.Sprintf("%d:%t:%s", result.rune, result.bad, strings.Join(fields, ","))
}
