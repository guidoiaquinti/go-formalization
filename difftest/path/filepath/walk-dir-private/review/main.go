package main

import (
	"log"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct{ input, want string }{
		{"0 0 0 0 1 0 0 0 0", "ok|root:d,root/a:f,root/d:d,root/z:f"},
		{"0 0 0 0 0 1 0 0 0", "ok|root:d,root/a:f,root/d:d,root/d/b:f,root/z:f"},
		{"0 0 0 1 0 0 0 0 0", "ok|root:d,root/a:f"},
		{"0 0 0 0 0 2 0 0 0", "skipAll|root:d,root/a:f,root/d:d,root/d/b:f"},
		{"0 1 0 0 0 0 0 0 0", "ok|root:d,root:d:7,root/a:f,root/d:d,root/d/b:f,root/d/c:f,root/z:f"},
		{"0 1 0 0 0 0 0 0 1", "ok|root:d,root:d:7"},
		{"0 1 0 0 0 0 0 0 2", "skipAll|root:d,root:d:7"},
		{"0 1 1 0 0 0 0 0 3", "ok|root:d"},
		{"1 0 0 0 0 1 0 0 0", "ok|root:d,root\\a:f,root\\d:d,root\\d\\b:f,root\\z:f"},
	}
	var input strings.Builder
	for _, c := range cases {
		input.WriteString(c.input + "\n")
	}
	for i, got := range check.LeanRows("godiff-path-filepath-walk-dir-private", len(cases), strings.NewReader(input.String())) {
		if got != cases[i].want {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, got, cases[i].want)
		}
	}
	const data = "difftest/path/filepath/walk-dir-private/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanWalkDirReview", data+"testdata/bridge_test.go", data+"review/testdata/review_test.go"); err != nil {
		log.Fatal(err)
	}
}
