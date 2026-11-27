package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Tag, Repeat   int
	Count, Offset int64
	Fill          byte
	Old, Tail     []byte
	Want          string
}

func main() {
	fixed := []fixture{
		{0, 65534, 0, 0, 'a', nil, []byte("µ"), "tooLong:none|nil:none"},
		{0, 65533, 0, 0, 'a', nil, []byte("µ"), "nil:none|nil:none"},
		{2, 65533, 0, math.MinInt64, 255, []byte{0, 128}, []byte("µ"),
			"ok:0:" + strings.Repeat("255,", 65533) + "194,181|ok:-9223372036854775808:0,128"},
		{2, 65534, -1, 0, 'a', []byte{0, 255}, []byte("µ"), "tooLong:-1:0,255|afterWrite:-1:0,255"},
		{2, 0, 0, math.MaxInt64, 0, []byte{120}, []byte{0, 255}, "ok:0:0,255|ok:9223372036854775807:120"},
		{2, 0, 8, -7, 0, []byte{255}, nil, "ok:8:|afterWrite:8:255"},
		{1, 0, 0, 3, 0, []byte{128}, []byte{255, 0}, "ok:none:255,0|nil:none:128"},
	}
	f, err := os.CreateTemp("", "golean-zip-writer-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range fixed {
		fmt.Fprint(&input, c.Tag, " ", c.Count, " ", c.Offset, " ", c.Repeat, " ", c.Fill, " ", len(c.Old))
		for _, raw := range [][]byte{c.Old, c.Tail} {
			for _, cell := range raw {
				fmt.Fprint(&input, " ", cell)
			}
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-zip-writer-config", len(fixed), &input) {
		if got != fixed[i].Want {
			log.Fatalf("fixture %d: Lean disagrees with literal outcome", i)
		}
	}
	must(os.Setenv("GOLEAN_ZIP_WRITER_CONFIG", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanZipWriterConfig", "difftest/archive/zip/writer-config/testdata/config_test.go"))
	must(check.TestInPackage("archive/zip", "TestGoLeanWriterConfigReview", "difftest/archive/zip/writer-config/review/testdata/config_review_test.go"))
	fmt.Printf("%d exact Go/Lean rows and native mutation traces passed\n", len(fixed))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
