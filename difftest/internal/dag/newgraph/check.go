package main

import (
	"fmt"
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const want = "nil 0 present 0 present 0"
	const fixture = "difftest/internal/dag/newgraph/testdata/newgraph_test.go"
	if row := check.LeanRows("godiff-internal-dag-newgraph", 1, nil)[0]; row != want {
		log.Fatalf("newGraph: Lean=%q, want=%q", row, want)
	}
	if err := check.TestInPackage("internal/dag", "TestGoLeanNewGraph", fixture); err != nil {
		log.Fatal(err)
	}
	fmt.Println("internal/dag.newGraph: Go and Lean agree on the constructor's only input.")
}
