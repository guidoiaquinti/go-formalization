package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

//go:embed testdata/live_walk_test.go
var failureFixture []byte

func main() {
	cases := []struct{ input, want string }{
		{"0 1 0 0 0 0 0 0 0", "ok|root:d:7"},
		{"0 1 0 0 0 0 0 0 1", "skipDir|root:d:7"},
		{"0 1 0 0 0 0 0 0 2", "skipAll|root:d:7"},
		{"0 1 0 0 0 0 0 0 3", "9|root:d:7"},
		{"0 2 0 0 0 0 0 0 1", "ok|root:d,root/a:nil:8,root/d:d,root/d/b:f,root/d/c:f,root/z:f"},
		{"0 2 0 0 0 0 0 0 3", "9|root:d,root/a:nil:8"},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&rows, c.input)
	}
	for i, got := range check.LeanRows("godiff-path-filepath-walk-private", len(cases), &rows) {
		if got != cases[i].want {
			log.Fatalf("model row %d: got %q want %q", i, got, cases[i].want)
		}
	}
	f, err := os.CreateTemp("", "golean-review-live-walk-*_test.go")
	must(err)
	defer os.Remove(f.Name())
	_, err = f.Write(failureFixture)
	must(err)
	must(f.Close())
	const bridge = "difftest/path/filepath/walk-private/testdata/bridge_test.go"
	must(check.TestInPackage("path/filepath", "TestGoLeanPrivateWalkLiveFailure", bridge, f.Name()))
	fmt.Printf("private walk passed %d fixed model rows and eight live error scenarios\n", len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
