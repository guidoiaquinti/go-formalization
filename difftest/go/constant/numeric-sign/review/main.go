package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
	"log"
	"os"
)

func main() {
	fixed := []struct {
		X, Y int64
		Mode int
		Want string
	}{
		{0, 0, 0, "3|3|4|4|5|0|0|0|0|0|0|1|1"},
		{5, -7, 0, "3|3|4|4|5|1|1|1|1|-1|-1|1|-1"},
		{-5, 7, 0, "3|3|4|4|5|-1|-1|-1|-1|-1|-1|1|1"},
		{0, -2, 1, "3|3|4|4|5|0|0|0|0|-1|-1|panic:<nil> not numeric|panic:<nil> not numeric"},
		{0, 2, 2, "3|3|4|4|5|0|0|0|0|1|1|panic:true not numeric|panic:true not numeric"},
		{-9223372036854775808, 0, 0, "3|3|4|4|5|-1|-1|-1|-1|-1|-1|1|1"},
	}
	f, err := os.CreateTemp("", "golean-constant-sign-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range fixed {
		fmt.Fprintln(&input, c.X, c.Y, c.Mode)
		must(enc.Encode(c))
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-go-constant-numeric-sign", len(fixed), &input) {
		if got != fixed[i].Want {
			log.Fatalf("fixture %d: Lean %q expected %q", i, got, fixed[i].Want)
		}
	}
	must(os.Setenv("GOLEAN_CONSTANT_SIGN_ROWS", f.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanNumericSign", "difftest/go/constant/numeric-sign/testdata/sign_test.go"))
	must(check.TestInPackage("go/constant", "TestGoLeanSignReview", "difftest/go/constant/numeric-sign/review/testdata/review_test.go"))
	fmt.Printf("%d literal numeric-sign rows plus native signed-zero and panic controls\n", len(fixed))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
