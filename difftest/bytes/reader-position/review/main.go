package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Tag, Previous, Whence, Script, Count int
	Cursor, Offset                       int64
	Raw, Spare                           []byte
}

type fixture struct {
	Input input
	Want  string
}

func main() {
	raw, spare := []byte{11, 22, 33}, []byte{99}
	cases := []fixture{
		{input{}, "nilReader|nilReader"},
		{input{2, 8, 9, 0, 1, 1, 0, raw, spare}, "whence:0:1:-1|shortWrite:1:2:-1:3:22,33"},
		{input{2, 8, 1, 3, 0, math.MaxInt64, 1, raw, spare}, "negative:0:9223372036854775807:-1|ok:0:9223372036854775807:-1:none"},
		{input{2, 8, 2, 1, 1, 1, -2, raw, spare}, "ok:1:1:-1|writerError:1:2:-1:3:22,33"},
		{input{2, 8, 0, 0, 3, 1, 2, raw, spare}, "ok:2:2:-1|invalidCount:0:1:-1:3:22,33"},
		{input{2, 8, 0, 2, 0, 1, -1, raw, spare}, "negative:0:1:-1|writerPanic:0:1:-1:3:22,33"},
		{input{2, 8, 0, 0, -2, 1, 2, raw, spare}, "ok:2:2:-1|shortWrite:-2:-1:-1:3:22,33"},
		{input{2, 8, 0, 3, 0, 1, 0, raw, spare}, "ok:0:0:-1|nilWriter:0:1:-1:none"},
	}
	mutants := []struct {
		row    int
		broken string
	}{
		{1, "whence:0:1:8|shortWrite:1:2:-1:3:22,33"},
		{2, "ok:9223372036854775807:9223372036854775807:-1|ok:0:9223372036854775807:-1:none"},
		{3, "ok:1:1:-1|shortWrite:1:2:-1:3:22,33"},
		{4, "ok:2:2:-1|invalidCount:0:4:-1:3:22,33"},
		{6, "ok:2:2:-1|invalidCount:0:1:-1:3:22,33"},
	}
	var wire bytes.Buffer
	for _, f := range cases {
		c := f.Input
		fmt.Fprintf(&wire, "%d %d %d %d %d %d %d %d", c.Tag, len(c.Raw), c.Cursor, c.Previous, c.Offset, c.Whence, c.Script, c.Count)
		for _, values := range [][]byte{c.Raw, c.Spare} {
			for _, b := range values {
				fmt.Fprintf(&wire, " %d", b)
			}
		}
		fmt.Fprintln(&wire)
	}
	rows := check.LeanRows("godiff-bytes-reader-position", len(cases), &wire)
	for i, f := range cases {
		if rows[i] != f.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], f.Want)
		}
	}
	file, err := os.CreateTemp("", "golean-reader-position-review")
	must(err)
	defer os.Remove(file.Name())
	must(json.NewEncoder(file).Encode(cases))
	must(file.Close())
	must(os.Setenv("GOLEAN_READER_POSITION_CASES", file.Name()))
	const testdata = "difftest/bytes/reader-position/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanReaderPosition", testdata+"bridge_test.go", testdata+"reader_position_test.go"))
	for _, m := range mutants {
		if rows[m.row] == m.broken {
			log.Fatalf("accepted mutant at fixture %d", m.row)
		}
	}
	fmt.Println("Eight independent literal traces matched Go and Lean; five mutant traces rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
