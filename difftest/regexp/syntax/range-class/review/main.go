package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		Backing, Clean []rune
		Size, I, J     int
		Lo, Hi, Rune   rune
		Want           string
	}{
		{[]rune{1, 2, 9, 10, 101, 102}, []rune{3, 4}, 2, 0, 1, 3, 4, 4,
			"2|true:9,10,1,2,101,102|1,4,9,10:1,4,9,10,101,102:true|1,4,9,10:1,4,9,10,101,102:true|" +
				"0,2,5,1114111:0,2,5,1114111,101,102:true|0,2,5,1114111:0,2:false|true"},
		{[]rune{1, 2, 4, 5}, []rune{0, 1114111}, 2, 0, 1, 3, 3, 1114111,
			"2|true:4,5,1,2|1,2,3,5:1,2,3,5:true|1,2,0,1114111:1,2,0,1114111:true|" +
				"1,2,4,5:1,2,4,5:true|:0,1114111:true|true"},
		{nil, nil, 0, 0, 0, 0, 0, 0,
			"0|skip|0,0::false|::true|0,1114111::false|0,1114111::false|false"},
		{[]rune{7, 11, 7, 9}, nil, 2, 0, 1, 0, 0, 0,
			"2|true:7,9,7,11|7,11,7,9,0,0:7,11,7,9:false|7,11,7,9:7,11,7,9:true|" +
				"7,11,0,1114111:7,11,0,1114111:true|0,1114111::false|false"},
	}
	var input bytes.Buffer
	file, err := os.CreateTemp("", "range-class-review-*")
	must(err)
	defer os.Remove(file.Name())
	encoder := gob.NewEncoder(file)
	for _, c := range cases {
		fmt.Fprint(&input, c.Size, " ", len(c.Backing)/2, " ", c.Lo, " ", c.Hi, " ", c.Rune, " ", c.I, " ", c.J)
		for _, values := range [][]rune{c.Backing, c.Clean} {
			for _, value := range values {
				fmt.Fprint(&input, " ", value)
			}
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_REGEXP_RANGE_CLASS_ROWS", file.Name()))
	must(check.TestInPackage("regexp/syntax", "TestGoLeanRangeClass",
		"difftest/regexp/syntax/range-class/testdata/class_test.go"))
	rows := check.LeanRows("godiff-regexp-range-class", len(cases), &input)
	for i, got := range rows {
		if got != cases[i].Want {
			log.Fatalf("row %d: got %q, want %q", i, got, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row, column int
		wrong       string
	}{
		{0, 0, "4"},
		{3, 1, "false:7,9,7,11"},
		{0, 1, "true:1,2,9,10,101,102"},
		{1, 2, "1,3,4,5:1,3,4,5:true"},
		{0, 3, "1,2,9,10,3,4:1,2,9,10,3,4:true"},
		{1, 4, "0,1114111,4,5:0,1114111,4,5:true"},
		{2, 5, "::true"},
		{0, 6, "false"},
	} {
		if strings.Split(rows[m.row], "|")[m.column] == m.wrong {
			log.Fatalf("row %d column %d accepted mutant", m.row, m.column)
		}
	}
	fmt.Println("Four literal native and Lean cases reject eight wrong outcomes.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
