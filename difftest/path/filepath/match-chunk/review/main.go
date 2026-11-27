package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type result struct {
	rest    string
	ok, bad bool
}

type fixture struct {
	chunk, input  string
	unix, windows result
}

func same(chunk, input string, want result) fixture {
	return fixture{chunk, input, want, want}
}

func yes(rest string) result { return result{rest: rest, ok: true} }

func wire(value result) string {
	bytes := make([]string, len(value.rest))
	for i := range len(value.rest) {
		bytes[i] = strconv.Itoa(int(value.rest[i]))
	}
	return fmt.Sprintf("%t:%t:%s", value.ok, value.bad, strings.Join(bytes, ","))
}

func main() {
	bad, no := result{bad: true}, result{}
	tests := []fixture{
		same("", "\xff\x00", yes("\xff\x00")),
		same("a", "", no),
		same("a[", "b", bad),
		same("a[x-]", "b", bad),
		same("a[]", "", bad),
		{"a\\", "b", bad, no},
		{"?", "/tail", no, yes("tail")},
		{"?", "\\tail", yes("tail"), no},
		same("?", "\xfftail", yes("tail")),
		same("?", "🚀tail", yes("tail")),
		same("[/]", "/tail", yes("tail")),
		same("[^/]", "/tail", no),
		same("[z-a]", "ztail", no),
		same("[^z-a]", "ztail", yes("tail")),
		same("[α-γ]", "βtail", yes("tail")),
		same("[�]", "\xfftail", yes("tail")),
		same("[\xff]", "\xfftail", bad),
		same("\xfe", "\xfftail", no),
		{"\\\xff", "\xfftail", yes("tail"), no},
		{"\\*", "*tail", yes("tail"), no},
		{"\\?", "\\éTAIL", no, yes("TAIL")},
	}
	var requests, native bytes.Buffer
	var skippedSyntax, extraWindowsSeparator, decodedLiteral int
	for _, test := range tests {
		fmt.Fprint(&requests, len(test.chunk))
		for _, raw := range []string{test.chunk, test.input} {
			for _, value := range []byte(raw) {
				fmt.Fprintf(&requests, " %d", value)
			}
		}
		fmt.Fprintln(&requests)
		want := test.unix
		if runtime.GOOS == "windows" {
			want = test.windows
		}
		fmt.Fprintf(&native, "x%x x%x %t %t x%x\n", test.chunk, test.input, want.ok, want.bad, want.rest)
		if len(test.chunk) > 0 && test.chunk[0] == 'a' &&
			(len(test.input) == 0 || test.input[0] != 'a') && test.unix.bad {
			skippedSyntax++
		}
		if test.chunk == "?" && len(test.input) > 0 && test.input[0] == '/' && test.windows.ok {
			extraWindowsSeparator++
		}
		if test.chunk == "\xfe" && !test.unix.ok {
			left, _ := utf8.DecodeRuneInString(test.chunk)
			right, _ := utf8.DecodeRuneInString(test.input)
			if left == right {
				decodedLiteral++
			}
		}
	}
	for i, got := range check.LeanRows("godiff-path-filepath-match-chunk", len(tests), &requests) {
		test := tests[i]
		if want := wire(test.unix) + "|" + wire(test.windows); got != want {
			log.Fatalf("matchChunk(%q,%q): Lean=%s want=%s", test.chunk, test.input, got, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-match-chunk-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, native.Bytes(), 0600))
	must(os.Setenv("GOLEAN_FILEPATH_MATCH_CHUNK_CASES", path))
	const overlay = "difftest/path/filepath/match-chunk/testdata/"
	must(check.TestInPackage("path/filepath", "TestGoLeanMatchChunk",
		overlay+"bridge_test.go", overlay+"match_chunk_test.go"))
	if skippedSyntax != 4 || extraWindowsSeparator != 1 || decodedLiteral != 1 {
		log.Fatalf("controls: skipped syntax=%d, extra Windows separator=%d, decoded literal=%d",
			skippedSyntax, extraWindowsSeparator, decodedLiteral)
	}
	fmt.Printf("matchChunk: %d native and %d model observations; skipped syntax=%d, extra Windows separator=%d, decoded literal=%d mutant failures\n",
		len(tests), 2*len(tests), skippedSyntax, extraWindowsSeparator, decodedLiteral)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
