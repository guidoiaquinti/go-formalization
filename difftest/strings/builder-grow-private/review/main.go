package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const path = "difftest/strings/builder-grow-private/"
	if err := check.TestInPackage("strings", "TestGoLeanBuilderGrowPrivateReview", path+"testdata/bridge_test.go", "difftest/strings/builder-copy-check/testdata/bridge_test.go", path+"review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
}
