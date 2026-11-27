package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const root = "difftest/encoding/binary/ensure/"
	cases := []struct{ in, want string }{
		{"1 0 0", "true:false::"},
		{"0 0 0", "false:false::"},
		{"0 2 0 11 22 33", "false:false:11,22:"},
		{"0 2 1 11 22 33 44", "false:false:11,22,33:33"},
		{"0 2 2 11 22 33 44", "false:false:11,22,33,44:33,44"},
		{"0 2 3 11 22 33 44", "false:true:11,22,33,44,0:33,44,0"},
		{"0 0 3 88 99", "false:true:88,99,0:88,99,0"},
		{"1 0 2", "false:true:0,0:0,0"},
		{"0 2 -1 11 22 33", "panic:cannot be negative"},
	}
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&wire, c.in)
	}
	for i, got := range check.LeanRows("godiff-binary-ensure", len(cases), &wire) {
		if got != cases[i].want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, cases[i].want)
		}
	}
	if err := check.TestInPackage("encoding/binary", "TestGoLeanEnsureReview", root+"testdata/bridge_test.go", root+"review/testdata/ensure_review_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ensure: %d independent Lean rows plus native capacity and write-through controls\n", len(cases))
}
