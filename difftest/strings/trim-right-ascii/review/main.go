package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const path = "difftest/strings/trim-right-ascii/"
	if err := check.TestInPackage("strings", "TestGoLeanTrimRightASCIIReview",
		path+"testdata/bridge_test.go", path+"review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
}
