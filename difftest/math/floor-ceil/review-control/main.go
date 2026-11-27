package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const dir = "difftest/math/floor-ceil/"
	if err := check.TestInPackage("math", "TestGoLeanGenericFloorCeilControl",
		dir+"testdata/bridge_test.go", dir+"review-control/testdata/control_test.go"); err != nil {
		log.Fatal(err)
	}
}
