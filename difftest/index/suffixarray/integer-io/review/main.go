package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Length, Chunk, Reply int
	Value                int
	Backing              []byte
	Integers             []int64
	Want                 string
}

func main() {
	padding := func(n int) string { return strings.Repeat(",99", n) }
	cases := []fixture{
		{1, 3, 1, -1, bytes.Repeat([]byte{99}, 12), nil,
			"other/0/1" + padding(11) + "/10,12:1" + padding(9) + "|eof/-50/99" + padding(11) + "/10|other/0/20" + padding(11) + "/10,12:20" + padding(9)},
		{10, 3, 0, 0, bytes.Repeat([]byte{1}, 10), nil,
			"ok/0/0,1,1,1,1,1,1,1,1,1/10,10:0,1,1,1,1,1,1,1,1,1|eof/-1/1,1,1,1,1,1,1,1,1,1/10|ok/0/20,1,1,1,1,1,1,1,1,1/10,10:20,1,1,1,1,1,1,1,1,1"},
		{19, 1, 0, 0, bytes.Repeat([]byte{99}, 20), []int64{1},
			"ok/0/0" + padding(19) + "/10,20:0" + padding(9) + "|eof/-50/99" + padding(19) + "/10|ok/0/20" + padding(19) + "/10,20:20" + padding(9)},
	}
	var input bytes.Buffer
	file, err := os.CreateTemp("", "suffixarray-integer-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d %d 0 %d 0 0", strconv.IntSize, c.Length, len(c.Backing), c.Value, c.Chunk, c.Reply)
		for _, b := range c.Backing {
			fmt.Fprint(&input, " ", b)
		}
		for _, x := range c.Integers {
			fmt.Fprint(&input, " ", x)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(file.Close())
	rows := check.LeanRows("godiff-suffixarray-integer-io", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d: got %q want %q", i, row, cases[i].Want)
		}
	}
	must(os.Setenv("GOLEAN_SUFFIXARRAY_INTEGER_IO_CASES", file.Name()))
	must(check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayIntegerIO",
		"difftest/index/suffixarray/integer-io/testdata/integer_io_test.go"))
	mutants := []string{"other/0/1" + padding(11) + "/1,12:1", "eof/0/1,1,1,1,1,1,1,1,1,1/10",
		"ok/1/22" + padding(9) + ",1" + padding(9) + "/11,20:22" + padding(9) + ",1"}
	parts := []int{0, 1, 2}
	for i, part := range parts {
		if strings.Split(rows[i], "|")[part] == mutants[i] {
			log.Fatalf("row %d accepted a broken framing result", i)
		}
	}
	fmt.Println("Three native and Lean witnesses distinguish truncated writes, early read-error returns, and packing without full-word room.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
