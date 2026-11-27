package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const path = "difftest/strings/single-string-replacer-write-string/"
	if err := check.TestInPackage("strings", "TestGoLeanSingleStringReplacerWriteStringReview",
		path+"testdata/bridge_test.go", path+"review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
}
