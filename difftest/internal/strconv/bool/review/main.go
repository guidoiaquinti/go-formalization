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
	Bool, Nil   bool
	Tag, Length int
	Backing     []byte
	Raw, Want   string
}

func main() {
	cases := []fixture{
		{true, false, 2, 1, []byte{19, 20, 21, 22, 23}, "False", "false|false|true|19,116,114,117,101|19,116,114,117,101|true|invalid syntax"},
		{false, false, 3, 1, []byte{19, 20, 21, 22, 23}, "TrUe", "false|true|false|19,102,97,108,115,101|19,20,21,22,23|false|invalid base"},
		{true, true, math.MinInt, 0, nil, "1", "true|false|true|116,114,117,101||false|unknown error"},
		{true, false, 4, 1, []byte{255, 20, 21, 22, 23, 24}, " true", "false|true|true|255,116,114,117,101|255,116,114,117,101,24|true|invalid bit size"},
		{false, true, 1, 0, nil, "true\xff", "false|true|false|102,97,108,115,101||false|value out of range"},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-strconv-bool-")
	must(err)
	defer os.Remove(truth.Name())
	enc := gob.NewEncoder(truth)
	for _, c := range cases {
		flag := 0
		if c.Bool {
			flag = 1
		}
		fmt.Fprint(&input, flag, " ", c.Tag, " ", c.Length, " ", len(c.Backing))
		for _, b := range c.Backing {
			fmt.Fprint(&input, " ", b)
		}
		for _, b := range []byte(c.Raw) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_INTERNAL_STRCONV_BOOL_ROWS", truth.Name()))
	must(check.TestInPackage("internal/strconv", "TestGoLeanInternalBool", "difftest/internal/strconv/bool/testdata/bool_test.go"))
	rows := check.LeanRows("godiff-internal-strconv-bool", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row, col int
		value    string
	}{
		{0, 1, "true"},                  // reject admissible title-style false
		{1, 0, "true"},                  // parse arbitrary mixed case
		{1, 1, "false"},                 // omit the syntax error
		{3, 1, "false"},                 // accept leading whitespace
		{3, 2, "True"},                  // format title case
		{1, 4, "19,102,97,108,115"},     // write partially before growth
		{3, 4, "255,116,114,117,101,0"}, // clear the unused tail
		{0, 5, "false"},                 // allocate on exact fit
		{2, 6, "invalid syntax"},        // use a known error for an unknown tag
		{4, 6, "invalid base"},          // mix up known error codes
	} {
		if strings.Split(rows[m.row], "|")[m.col] == m.value {
			log.Fatalf("row %d column %d accepted mutant %q", m.row, m.col, m.value)
		}
	}
	log.Print("five independent pinned Go and Lean rows passed and rejected ten wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
