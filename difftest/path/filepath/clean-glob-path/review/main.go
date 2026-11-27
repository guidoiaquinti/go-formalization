package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const path = "difftest/path/filepath/clean-glob-path/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanCleanGlobPathReview",
		path+"testdata/bridge_test.go", path+"review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
}
