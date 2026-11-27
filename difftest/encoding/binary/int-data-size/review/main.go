package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const dir = "difftest/encoding/binary/int-data-size/"
	cases := []struct{ wire, want string }{
		{"2 5 0 0", "4:true:0:"},
		{"2 10 0 0", "8:true:0:"},
		{"4 5 0 0", "0:true:0:"},
		{"5 5 2 0", "0:true:0:"},
		{"3 2 2 0 11 22 33", "2:false:3:11,22,33"},
		{"3 2 0 0 44", "0:false:1:44"},
		{"3 2 0 1", "0:true:0:"},
		{"3 5 2 0", "8:true:0:"},
	}
	var wire strings.Builder
	for _, c := range cases {
		fmt.Fprintln(&wire, c.wire)
	}
	for i, got := range check.LeanRows("godiff-binary-int-data-size", len(cases), strings.NewReader(wire.String())) {
		if got != cases[i].want {
			log.Fatalf("case %d: Lean %q, want %q", i, got, cases[i].want)
		}
	}
	if err := check.TestInPackage("encoding/binary", "TestGoLeanIntSizeReview", dir+"testdata/bridge_test.go", dir+"review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("intDataSize: %d independent Lean rows and native identity controls\n", len(cases))
}
