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

type fixture struct {
	Mode, Bins int
	Text       []byte
	Prefix     []int
	Want       string
}

func suite(byteMode, intMode string) string {
	return strings.Join([]string{byteMode, byteMode, intMode, intMode}, "|")
}

func main() {
	const suffixes = "3.1,1.1,0.3"
	const byteBucket = "0.2,3.1,5.253,91.1,92.1"
	const wideBucket = "0.3,2.1,3.2,5.1,91.1,92.1"
	const byteFreq = "0.1,2.1,1.1,2.1,0.252,93.1,94.1"
	const wideFreq = "0.2,2.1,0.1,1.1,0.1,2.1,93.1,94.1"
	noCache := suite(suffixes+"/nil/"+byteBucket, suffixes+"/nil/"+wideBucket)
	cached := suite(suffixes+"/"+byteFreq+"/"+byteBucket, suffixes+"/"+wideFreq+"/"+wideBucket)
	cases := []fixture{
		{0, 7, []byte{3, 1, 3, 1, 2}, []int{3, 1}, noCache},
		{1, 7, []byte{3, 1, 3, 1, 2}, []int{3, 1}, cached},
		{2, 7, []byte{3, 1, 3, 1, 2}, []int{3, 1}, cached},
		{0, 511, []byte{255, 0, 255, 1, 255}, []int{1, 3},
			suite("1.1,3.1,0.3/nil/0.1,1.1,2.253,5.1,91.1,92.1",
				"1.1,3.1,0.3/nil/0.1,1.2,2.507,5.1,91.1,92.1")},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-suffix-expand-")
	must(err)
	defer os.Remove(truth.Name())
	enc := gob.NewEncoder(truth)
	for _, c := range cases {
		fmt.Fprint(&input, c.Mode, " ", c.Bins, " ", len(c.Text))
		for _, b := range c.Text {
			fmt.Fprint(&input, " ", b)
		}
		for _, p := range c.Prefix {
			fmt.Fprint(&input, " ", p)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_SUFFIX_EXPAND_ROWS", truth.Name()))
	must(check.TestInPackage("index/suffixarray", "TestGoLeanExpandLms", "difftest/index/suffixarray/expand-lms/testdata/expand_test.go"))
	rows := check.LeanRows("godiff-suffixarray-expand-lms", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row, field int
		first, end int
		value      string
	}{
		{0, 0, 0, 4, "1.1,3.1,0.3"},                 // reverse equal-symbol LMS placements
		{0, 0, 0, 4, "3.1,1.1,-7.3"},                // leave unselected cells untouched
		{1, 1, 0, 2, "0.1,2.1,1.1,2.1,0.252,0.2"},   // clear spare byte-frequency backing
		{1, 1, 2, 4, "0.2,2.1,0.1,1.1,0.1,2.1,0.2"}, // clear spare wide-frequency backing
		{1, 2, 0, 2, "0.2,3.1,5.253,0.2"},           // clear spare byte-bucket backing
		{1, 2, 2, 4, "0.3,2.1,3.2,5.1,0.2"},         // clear spare wide-bucket backing
		{2, 2, 2, 4, byteBucket},                    // use byte-bucket geometry for the wide path
	} {
		fields := strings.Split(rows[m.row], "|")
		for variant := m.first; variant < m.end; variant++ {
			if strings.Split(fields[variant], "/")[m.field] == m.value {
				log.Fatalf("row %d variant %d accepted mutated field %q", m.row, variant, m.value)
			}
		}
	}
	log.Print("four independent pinned Go and Lean rows passed and rejected seven wrong field expectations across the profiles")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
