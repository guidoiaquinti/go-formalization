package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const path = "difftest/strings/trim-left-unicode/"
	if err := check.TestInPackage("strings", "TestGoLeanTrimLeftUnicodeReview",
		path+"testdata/bridge_test.go", path+"review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
}
