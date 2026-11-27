package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	profile int
	path    string
	devices []string
	want    string
	ok      bool
}

func main() {
	cases := []fixture{
		{0, ".", nil, ".", true},
		{0, "./a", nil, "", false},
		{0, "x/..", nil, "", false},
		{0, "x//a", nil, "", false},
		{0, "x/", nil, "", false},
		{0, "/x", nil, "", false},
		{0, "x/\x00a", nil, "", false},
		{0, "x/\xff", nil, "", false},
		{0, "x/\xed\xa0\x80", nil, "", false},
		{0, "x/\uFFFD", nil, "x/\uFFFD", true},
		{0, "x/\uFFFF", nil, "x/\uFFFF", true},
		{0, "x/..hidden", nil, "x/..hidden", true},
		{0, "#x/a:b\\c", nil, "#x/a:b\\c", true},
		{1, ".", nil, ".", true},
		{1, "x/y/\uFFFD", nil, "x\\y\\\uFFFD", true},
		{1, "x//y", nil, "", false},
		{1, "x/y\\z", nil, "", false},
		{1, "x/y:z", nil, "", false},
		{1, "x/y\x00z", nil, "", false},
		{1, "x/CONOUT$", nil, "", false},
		{1, "x/COM³", nil, "", false},
		{1, "x/COM⁴", nil, "x\\COM⁴", true},
		{1, "x/PrN .log", nil, "x\\PrN .log", true},
		{1, "x/PrN .log", []string{"PrN .log"}, "", false},
		{1, "x/PrN .log", []string{"prn .log", "PrN"}, "x\\PrN .log", true},
		{1, "x/ordinary", []string{"ordinary"}, "x\\ordinary", true},
		{2, "#x/a:b\\c", nil, "", false},
		{2, "x/#a:b\\c", nil, "x/#a:b\\c", true},
		{2, "./#x", nil, "", false},
		{2, ".", nil, ".", true},
	}
	var input bytes.Buffer
	for _, f := range cases {
		fmt.Fprintf(&input, "%d %d", f.profile, len(f.devices))
		for _, component := range f.devices {
			fmt.Fprintf(&input, " %d", len(component))
			writeBytes(&input, component)
		}
		writeBytes(&input, f.path)
		fmt.Fprintln(&input)
	}
	native := 0
	for i, got := range check.LeanRows("godiff-path-filepath-localize", len(cases), &input) {
		f := cases[i]
		if want := observation(f.want, f.ok); got != want {
			log.Fatalf("profile %d path %q devices %q: Lean=%s want=%s", f.profile, f.path, f.devices, got, want)
		}
		if (f.profile == 0 && runtime.GOOS != "plan9" && runtime.GOOS != "windows") ||
			(f.profile == 2 && runtime.GOOS == "plan9") {
			actual, failure := filepath.Localize(f.path)
			if actual != f.want || (failure == nil) != f.ok || (failure != nil && failure.Error() != "invalid path") {
				log.Fatalf("native %q: got %q, %v; want %q, success=%t", f.path, actual, failure, f.want, f.ok)
			}
			native++
		}
	}

	// Five plausible implementations differ on these witnesses.
	if !fs.ValidPath(cases[6].path) || cases[6].ok ||
		!filepath.IsLocal(cases[1].path) || cases[1].ok ||
		!cases[14].ok || cases[14].want == cases[14].path ||
		!cases[22].ok || cases[23].ok || !cases[24].ok ||
		!cases[12].ok || cases[26].ok || !cases[27].ok {
		log.Fatal("invalid validation, slash conversion, callback or device-root witnesses")
	}
	fmt.Printf("filepath.Localize: %d native and %d Lean profile fixtures; five incorrect variants distinguished\n", native, len(cases))
}

func writeBytes(out *bytes.Buffer, raw string) {
	for i := range len(raw) {
		fmt.Fprintf(out, " %d", raw[i])
	}
}

func observation(path string, ok bool) string {
	if !ok {
		return "invalid path"
	}
	encoded := make([]string, len(path))
	for i := range len(path) {
		encoded[i] = strconv.Itoa(int(path[i]))
	}
	return "ok:" + strings.Join(encoded, ",")
}
