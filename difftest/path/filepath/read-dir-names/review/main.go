package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	rows := []struct {
		topology string
		want     string
	}{
		{"0 7 9 1 1 88", "1|nil|7"},
		{"1 4 9 1 1 88", "3|nil|4"},
		{"1 0 9 0", "3|0|"},
		{"1 0 9 7 2 97 97 1 255 1 98 1 65 0 1 97 1 97", "3|7|;65;97;97;97,97;98;255"},
		{"1 0 0 4 2 0 255 2 0 0 1 0 1 128", "3|4|0;0,0;0,255;128"},
	}
	var model bytes.Buffer
	for _, r := range rows {
		fmt.Fprintln(&model, r.topology)
	}
	for i, got := range check.LeanRows("godiff-path-filepath-read-dir-names", len(rows), &model) {
		if got != rows[i].want {
			log.Fatalf("row %d: got %q want %q", i, got, rows[i].want)
		}
	}
	const dir = "difftest/path/filepath/read-dir-names/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanReadDirNames", dir+"bridge_test.go", dir+"read_dir_names_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("readDirNames passed %d fixed model outcomes and the four pinned host fixtures\n", len(rows))
}
