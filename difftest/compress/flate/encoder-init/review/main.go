package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	width        int
	size         int64
	native       bool
	want, mutant string
}

func main() {
	w := strconv.IntSize
	cases := []fixture{
		{w, -1, true, "bounds|65535 65535", "0 1 true|65535 65535"},
		{w, 1, true, "1 1 true|65535 65535", "1 2 true|65535 65535"},
		{w, 255, true, "255 256 true|65535 65535", "255 255 true|65535 65535"},
		{w, 256, true, "256 256 true|65535 65535", "256 512 true|65535 65535"},
		{w, 257, true, "257 512 true|65535 65535", "257 512 false|65535 65535"},
		{w, 513, true, "513 1024 true|65535 65535", "513 1024 true|65535 0"},
		{64, 0, false, "0 4294967296 true|65535 65535", "0 1 true|65535 65535"},
		{32, 0, false, "bounds|65535 65535", "0 1 true|65535 65535"},
		{32, 1073741825, false, "bounds|65535 65535", "1073741825 2147483648 true|65535 65535"},
		{64, 4294967297, false, "bounds|65535 65535", "4294967297 8589934592 true|65535 65535"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d\n", c.width, c.size)
	}
	out, err := os.CreateTemp("", "golean-flate-encoder-review-")
	must(err)
	defer os.Remove(out.Name())
	count, rejected := 0, 0
	for i, got := range check.LeanRows("godiff-flate-encoder-init", len(cases), &input) {
		c := cases[i]
		if got != c.want {
			log.Fatalf("row %d: Lean %q want %q", i, got, c.want)
		}
		if got != c.mutant {
			rejected++
		}
		if c.native {
			_, err := fmt.Fprintf(out, "%d %s\n", c.size, c.want)
			must(err)
			count++
		}
	}
	if rejected != len(cases) {
		log.Fatalf("rejected %d of %d mutants", rejected, len(cases))
	}
	must(out.Close())
	must(os.Setenv("GOLEAN_FLATE_INIT_ROWS", out.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanEncoderInit", "difftest/compress/flate/encoder-init/testdata/init_test.go"))
	fmt.Printf("Encoder init: %d Lean rows, %d native rows, %d altered outcomes rejected\n", len(cases), count, rejected)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
