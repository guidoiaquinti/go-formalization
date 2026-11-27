package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const path = "difftest/encoding/binary/reflect-size/"
	cases := []struct{ input, want string }{
		{"0", "panic:-1:-1"},
		{"8 2 8 5 8 3 7", "11:11:11"},
		{"9 0 21", "-1:-1:-1"},
		{"10 0 21", "-1:-1:-1"},
		{"9 0 3", "0:0:0"},
		{"11 0 11 0 5", "-1:-1:-1"},
		{"11 0 10 2 3", "-1:-1:4"},
		{"11 1 4", "-1:-1:-1"},
		{"10 3 8 2 8 5 8 3 7", "-1:33:33"},
		{"8 11 0 3 7", "-1:-1:-1"},
	}
	var wire strings.Builder
	for _, c := range cases {
		fmt.Fprintln(&wire, c.input)
	}
	for i, got := range check.LeanRows("godiff-binary-reflect-size", len(cases), strings.NewReader(wire.String())) {
		if got != cases[i].want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, cases[i].want)
		}
	}
	if err := check.TestInPackage("encoding/binary", "TestGoLeanReflectSizeReview", path+"testdata/bridge_test.go", path+"review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("reflect sizing: %d independent Lean rows and pinned native controls\n", len(cases))
}
