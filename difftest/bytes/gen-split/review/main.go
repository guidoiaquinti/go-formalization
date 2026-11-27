package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixtures := []struct {
		save, limit, spare int
		input, sep, want   string
	}{
		{3, -1, 2, "a:b", ":", "slicef4:97,58,98,165;f3:98"},
		{5, -1, 2, "a:b", ":", "panic"},
		{-1, -1, 2, "a:b", ":", "slicef0:;f3:98"},
		{-2, -1, 2, "a:b", ":", "panic"},
		{5, 1, 2, "a:b", ":", "slicef5:97,58,98"},
		{5, 0, 2, "a:b", ":", "nil"},
		{-100, -1, 2, "\xff", "", "slicef3:255"},
	}
	var rows strings.Builder
	for _, f := range fixtures {
		fmt.Fprintf(&rows, "%d %d %d x%x x%x %s\n", f.save, f.limit, f.spare, f.input, f.sep, f.want)
	}
	dir, err := os.MkdirTemp("", "golean-bytes-gensplit-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, []byte(rows.String()), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_BYTES_GEN_SPLIT_CASES", path); err != nil {
		log.Fatal(err)
	}
	const overlay = "difftest/bytes/gen-split/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanBytesGenSplit", overlay+"bridge_test.go", overlay+"gen_split_test.go"); err != nil {
		log.Fatal(err)
	}
	// The first fixture rejects a mutant that caps endpoints at length.
	f := fixtures[0]
	input := []byte(f.input)
	match := bytes.Index(input, []byte(f.sep))
	mutated := input[:min(match+f.save, len(input))]
	if bytes.Equal(mutated, []byte{97, 58, 98, 165}) {
		log.Fatal("clamp-to-length mutant agreed")
	}
	fmt.Printf("private Go genSplit passed %d fixtures\n", len(fixtures))
}
