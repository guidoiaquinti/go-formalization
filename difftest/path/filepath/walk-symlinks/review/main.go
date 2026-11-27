package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

//go:embed testdata/boundary_test.go
var nativeFixture []byte

func main() {
	cases := []struct {
		profile      int
		path, result string
		ok           bool
	}{
		{0, "dir/../chain/file", "dir/file", true},
		{0, "up/file", "dir/file", true},
		{0, "filelink/../dir/file", "notDirectory", false},
		{0, "cycle", "tooManyLinks", false},
		{0, "broken/file", "missing", false},
		{1, "C:/dir/../dir/file", `C:\dir\file`, true},
	}
	var model bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&model, c.profile)
		for i := range len(c.path) {
			fmt.Fprintf(&model, " %d", c.path[i])
		}
		fmt.Fprintln(&model)
	}
	for i, got := range check.LeanRows("godiff-path-filepath-walk-symlinks", len(cases), &model) {
		c := cases[i]
		want := c.result
		if c.ok {
			bytes := make([]string, len(want))
			for j := range len(want) {
				bytes[j] = strconv.Itoa(int(want[j]))
			}
			want = "ok:" + strings.Join(bytes, ",")
		}
		if got != want {
			log.Fatalf("model row %d: got %q want %q", i, got, want)
		}
	}
	f, err := os.CreateTemp("", "golean-review-walk-symlinks-*_test.go")
	must(err)
	defer os.Remove(f.Name())
	_, err = f.Write(nativeFixture)
	must(err)
	must(f.Close())
	const bridge = "difftest/path/filepath/walk-symlinks/testdata/bridge_test.go"
	must(check.TestInPackage("path/filepath", "TestGoLeanWalkSymlinksBoundary", bridge, f.Name()))
	fmt.Printf("walkSymlinks passed %d independent model rows and native link-limit and parent fixtures on Unix\n", len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
