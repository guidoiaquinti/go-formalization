package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const dir = "difftest/slices/overlaps/"
	if err := check.TestInPackage("slices", "TestGoLeanOverlapsReview",
		dir+"testdata/bridge_test.go", dir+"control/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
}
