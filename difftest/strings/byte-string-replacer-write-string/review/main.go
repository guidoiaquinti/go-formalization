package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const path = "difftest/strings/byte-string-replacer-write-string/"
	if err := check.TestInPackage("strings", "TestGoLeanByteStringReplacerWriteStringReview",
		path+"testdata/bridge_test.go", path+"review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
}
