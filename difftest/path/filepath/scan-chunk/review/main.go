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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type portions struct{ chunk, rest string }

type fixture struct {
	pattern       string
	star          bool
	unix, windows portions
}

func same(pattern string, star bool, chunk, rest string) fixture {
	p := portions{chunk, rest}
	return fixture{pattern, star, p, p}
}

func decimal(s string) string {
	values := make([]string, len(s))
	for i := range len(s) {
		values[i] = strconv.Itoa(int(s[i]))
	}
	return strings.Join(values, ",")
}

func wire(star bool, part portions) string {
	return fmt.Sprintf("%t:%s:%s", star, decimal(part.chunk), decimal(part.rest))
}

func main() {
	tests := []fixture{
		same("", false, "", ""),
		same("***", true, "", ""),
		same("***ab**cd", true, "ab", "**cd"),
		same("a*b*c", false, "a", "*b*c"),
		same("a\\", false, "a\\", ""),
		{"a\\*b*c", false, portions{"a\\*b", "*c"}, portions{"a\\", "*b*c"}},
		same("a\\\\*x", false, "a\\\\", "*x"),
		{"[\\]*z", false, portions{"[\\]*z", ""}, portions{"[\\]", "*z"}},
		{"\\[*z", false, portions{"\\[", "*z"}, portions{"\\[*z", ""}},
		same("[[*]a*b", false, "[[*]a", "*b"),
		same("[a*", false, "[a*", ""),
		same("a]*b", false, "a]", "*b"),
		same("***[\x80]\x00*\xff", true, "[\x80]\x00", "*\xff"),
	}
	var modelInput, nativeFixture bytes.Buffer
	var ignoredClass, escapedWindows, nestedClass int
	for _, test := range tests {
		fmt.Fprint(&modelInput, 0)
		for _, value := range []byte(test.pattern) {
			fmt.Fprintf(&modelInput, " %d", value)
		}
		fmt.Fprintln(&modelInput)
		want := test.unix
		if runtime.GOOS == "windows" {
			want = test.windows
		}
		fmt.Fprintf(&nativeFixture, "x%x %t x%x x%x\n", test.pattern, test.star, want.chunk, want.rest)
		body := strings.TrimLeft(test.pattern, "*")
		noClass := portions{body, ""}
		if at := strings.IndexByte(body, '*'); at >= 0 {
			noClass = portions{body[:at], body[at:]}
		}
		if noClass != test.windows {
			ignoredClass++
		}
		if test.unix != test.windows {
			escapedWindows++
		}
		if nested(body) != test.windows {
			nestedClass++
		}
	}
	for i, row := range check.LeanRows("godiff-path-filepath-scan-chunk", len(tests), &modelInput) {
		test := tests[i]
		want := wire(test.star, test.unix) + "|" + wire(test.star, test.windows)
		if row != want {
			log.Fatalf("scanChunk(%q): Lean=%s want=%s", test.pattern, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-scan-chunk-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, nativeFixture.Bytes(), 0600))
	must(os.Setenv("GOLEAN_FILEPATH_SCAN_CHUNK_CASES", path))
	const overlay = "difftest/path/filepath/scan-chunk/testdata/"
	must(check.TestInPackage("path/filepath", "TestGoLeanScanChunk", overlay+"bridge_test.go", overlay+"scan_chunk_test.go"))
	if ignoredClass != 3 || escapedWindows != 3 || nestedClass != 1 {
		log.Fatalf("mutants: ignored class=%d, Windows escapes=%d, nested class=%d",
			ignoredClass, escapedWindows, nestedClass)
	}
	fmt.Printf("scanChunk: %d native fixtures and %d platform model rows; ignored class=%d, Windows escapes=%d, nested class=%d mutant failures\n",
		len(tests), 2*len(tests), ignoredClass, escapedWindows, nestedClass)
}

// Deliberately incorrect Windows interpretation that nests brackets.
func nested(body string) portions {
	var depth int
	for i := range len(body) {
		switch body[i] {
		case '[':
			depth++
		case ']':
			depth = max(0, depth-1)
		case '*':
			if depth == 0 {
				return portions{body[:i], body[i:]}
			}
		}
	}
	return portions{body, ""}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
