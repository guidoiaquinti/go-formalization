package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := check.TestInPackage("container/heap", "TestGoLeanHeapPushReview",
		"difftest/container/heap/push/review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
}
