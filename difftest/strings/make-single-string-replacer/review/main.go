package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const path = "difftest/strings/make-single-string-replacer/review/testdata/"
	if err := check.TestInPackage("strings", "TestGoLeanSingleMakerReview",
		path+"bridge_test.go", path+"review_test.go"); err != nil {
		log.Fatal(err)
	}
}
