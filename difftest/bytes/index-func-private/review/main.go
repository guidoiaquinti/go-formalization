package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixtures := []struct {
		truth       bool
		kind        int
		input, want string
	}{
		{false, 0, "", "-1|"},
		{false, 0, "a", "panic"},
		{false, 3, "\xffé!", "1|65533,233"},
		{false, 3, "\xff\xef\xbf\xbd", "-1|65533,65533"},
		{true, 3, "é\xff", "2|233,65533"},
		{true, 4, "é🚀", "2|233,128640"},
	}
	var rows strings.Builder
	var ignoreTruth, ordinal int
	for _, f := range fixtures {
		fmt.Fprintf(&rows, "%t %d x%x %s\n", f.truth, f.kind, f.input, f.want)
		if !f.truth && f.kind == 3 {
			got := bytes.IndexFunc([]byte(f.input), func(r rune) bool { return r == utf8.RuneError })
			if !strings.HasPrefix(f.want, fmt.Sprintf("%d|", got)) {
				ignoreTruth++
			}
		}
		if f.truth && f.kind >= 3 {
			n := 0
			for _, r := range f.input {
				if f.kind == 3 && r == utf8.RuneError || f.kind == 4 && r == '🚀' {
					break
				}
				n++
			}
			if !strings.HasPrefix(f.want, fmt.Sprintf("%d|", n)) {
				ordinal++
			}
		}
	}
	if ignoreTruth != 2 || ordinal != 2 {
		log.Fatalf("ignore truth failed %d, rune ordinal failed %d", ignoreTruth, ordinal)
	}
	dir, err := os.MkdirTemp("", "golean-bytes-index-func-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, []byte(rows.String()), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_BYTES_INDEX_FUNC_CASES", path); err != nil {
		log.Fatal(err)
	}
	const overlay = "difftest/bytes/index-func-private/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanIndexFunc", overlay+"bridge_test.go", overlay+"index_func_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("private Go passed %d fixtures; both mutants failed two\n", len(fixtures))
}
