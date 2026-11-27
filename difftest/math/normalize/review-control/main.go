package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const dir = "difftest/math/normalize/"
	if err := check.TestInPackage("math", "TestGoLeanNormalizeControl",
		dir+"testdata/bridge_test.go", dir+"review-control/testdata/control_test.go"); err != nil {
		log.Fatal(err)
	}
}
