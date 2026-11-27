package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	creator       uint16
	attrs, legacy uint32
	wide          uint64
	present       bool
	name, want    string
}

func main() {
	cases := []fixture{
		{0, 0, 0, 0, false, "", "true true true"},
		{0xff00, 0, 7, math.MaxUint64, true, "old/\xff/", "true true true 255 -1 2147483648 2147483648 true"},
		{0, 1, math.MaxUint32, 0, true, "part/raw\\name", "true true true 114,97,119,92,110,97,109,101 4294967295 292 0 false"},
		{0xff00, 0, 0, 0, true, "", "true true true 46 0 0 0 false"},
		{0xff00, 0, 0, 0, true, "///", "true true true 47 0 2147483648 2147483648 true"},
		{0x0300, 0xa1ff0000, 5, 6, true, "link/", "true true true 108,105,110,107 6 2281701887 2281701376 true"},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-zip-header-info-")
	must(err)
	defer os.Remove(truth.Name())
	for _, c := range cases {
		present := 0
		if c.present {
			present = 1
		}
		fmt.Fprint(&input, c.creator, " ", c.attrs, " ", c.legacy, " ", c.wide, " ", present)
		for _, b := range []byte(c.name) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		_, err := fmt.Fprintf(truth, "%d %d %d %d %t %q\t%s\n",
			c.creator, c.attrs, c.legacy, c.wide, c.present, c.name, c.want)
		must(err)
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_ZIP_HEADER_INFO", truth.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanHeaderInfo", "difftest/archive/zip/header-info/testdata/header_info_test.go"))
	rows := check.LeanRows("godiff-zip-header-info", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].want)
		}
	}
	for _, m := range []struct {
		row, field int
		value      string
	}{
		{0, 0, "false"},       // discard the adapter for a nil header
		{0, 1, "false"},       // turn typed nil into a nil interface
		{0, 2, "false"},       // substitute a different adapter or nonnil error
		{1, 3, "239,191,189"}, // replace invalid UTF-8 in the basename
		{1, 4, "7"},           // ignore the nonzero wrapped wide size
		{2, 4, "-1"},          // sign-extend the 32-bit size
		{2, 5, "0"},           // return the pre-mutation header mode
		{5, 6, "2281701887"},  // retain permissions in Type
		{5, 7, "false"},       // ignore the live trailing-slash directory bit
	} {
		if strings.Fields(rows[m.row])[m.field] == m.value {
			log.Fatalf("row %d column %d accepted mutant %q", m.row, m.field, m.value)
		}
	}
	log.Print("six independent pinned Go and Lean rows passed and rejected nine wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
