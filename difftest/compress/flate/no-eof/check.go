package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	want := []string{"0", "2", "2", "3", "4", "5", "6", "7"}
	rows := check.LeanRows("godiff-flate-no-eof", len(want), strings.NewReader("0\n1\n2\n3\n4\n5\n6\n7\n"))
	if !slices.Equal(rows, want) {
		log.Fatalf("Lean returned %v, expected %v", rows, want)
	}
	f, err := os.CreateTemp("", "golean-flate-no-eof")
	must(err)
	defer os.Remove(f.Name())
	_, err = fmt.Fprintln(f, strings.Join(want, "\n"))
	must(err)
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_NO_EOF", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanNoEOF", "difftest/compress/flate/no-eof/testdata/noeof_test.go"))
	fmt.Println("Exact EOF replacement and unchanged identity of all seven other inputs verified.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
