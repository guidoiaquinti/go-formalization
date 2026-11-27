package main

import (
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected Go version %s", runtime.Version())
	}
	const fixture = "difftest/image/color/sq-diff/testdata/sq_diff_control_test.go"
	if err := check.TestInPackage("image/color", "TestGoLeanSqDiffControl", fixture); err != nil {
		log.Fatal(err)
	}
}
