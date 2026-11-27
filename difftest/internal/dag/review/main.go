package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const fixture = "difftest/internal/dag/review/testdata/review_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanDAGReview", fixture); err != nil {
		log.Fatal(err)
	}
}
