package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Count             int64
	Written, Existing int
	Failed            bool
	Data              []byte
	Want              string
}

func main() {
	cases := []fixture{
		{math.MaxInt64, 2, 4096, true, []byte{255, 0, 8},
			"-9223372036854775807:2:true:1|1:true:1:9223372036854775807|1:true:0:9223372036854775807|4096:1:true:0:0|false:0"},
		{-1, 3, 4095, true, []byte{3, 2, 1},
			"2:3:true:1|0:true:1:-1|0:true:0:-1|4096:0:false:0:0|false:0"},
		{7, 0, -1, true, nil,
			"7:0:true:1|0:false:0:7|0:false:0:7|4096:0:false:0:0|false:0"},
		{7, 4, 8192, false, []byte{1, 0, 255, 1},
			"11:4:false:1|0:false:1:7|0:false:0:7|8192:1:true:0:0|false:0"},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-zip-writer-io-")
	must(err)
	defer os.Remove(truth.Name())
	enc := gob.NewEncoder(truth)
	for _, c := range cases {
		failure := 0
		if c.Failed {
			failure = 1
		}
		fmt.Fprint(&input, c.Count, " ", c.Written, " ", failure, " ", c.Existing)
		for _, b := range c.Data {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_ZIP_WRITER_IO_ROWS", truth.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanZipWriterIO", "difftest/archive/zip/writer-io/testdata/writer_test.go"))
	rows := check.LeanRows("godiff-zip-writer-io", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row, col int
		value    string
	}{
		{0, 0, "9223372036854775807:2:true:1"}, // skip counting the writes that returned an error
		{2, 0, "7:0:false:0"},                  // skip delegation on empty input
		{0, 1, "0:true:1:9223372036854775807"}, // discard unwritten bytes after failing flush
		{0, 2, "1:true:1:9223372036854775807"}, // retry underlying call after sticky failure
		{1, 2, "0:false:0:-1"},                 // drop sticky error when no bytes remain
		{3, 1, "0:false:1:11"},                 // count flushed bytes a second time
		{0, 3, "4096:0:false:0:0"},             // allocate even at the threshold
		{1, 3, "4095:1:true:0:0"},              // reuse an undersized buffer
		{0, 4, "true:1"},                       // forward the wrapped closer's call and error
	} {
		if strings.Split(rows[m.row], "|")[m.col] == m.value {
			log.Fatalf("row %d column %d accepted mutant %q", m.row, m.col, m.value)
		}
	}
	log.Print("four independent pinned Go and Lean rows passed and rejected nine wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
