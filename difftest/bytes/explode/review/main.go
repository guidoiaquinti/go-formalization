package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixtures := []struct {
		input string
		limit int
		spare int
		want  string
	}{
		{"", 0, 0, "slice"},
		{"A", 0, 3, "slicef4:65"},
		{"é", 0, 3, "slicef2:195,169"},
		{"é", 1, 3, "slicef5:195,169"},
		{"\xffA", -1, 3, "slicef1:255;f4:65"},
		{"\xffA", 1, 3, "slicef5:255,65"},
	}
	var rows bytes.Buffer
	var zeroDrops int
	for _, f := range fixtures {
		fmt.Fprintf(&rows, "%d %d x%x %s\n", f.limit, f.spare, f.input, f.want)
		if f.limit == 0 && bytes.SplitN([]byte(f.input), nil, 0) == nil {
			zeroDrops++
		}
	}
	if zeroDrops != 3 {
		log.Fatalf("zero-limit nil mutant failed %d fixtures, want 3", zeroDrops)
	}
	dir, err := os.MkdirTemp("", "golean-bytes-explode-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_BYTES_EXPLODE_CASES", path); err != nil {
		log.Fatal(err)
	}
	const overlay = "difftest/bytes/explode/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanExplode", overlay+"bridge_test.go", overlay+"explode_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("private Go explode passed %d fixtures; zero-limit nil mutant failed %d\n", len(fixtures), zeroDrops)
}
