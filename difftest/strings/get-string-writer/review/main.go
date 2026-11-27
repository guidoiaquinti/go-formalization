package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const path = "difftest/strings/get-string-writer/review/testdata/"
	if err := check.TestInPackage("strings", "TestGoLeanGetStringWriterReview",
		path+"bridge_test.go", path+"review_test.go"); err != nil {
		log.Fatal(err)
	}
}
