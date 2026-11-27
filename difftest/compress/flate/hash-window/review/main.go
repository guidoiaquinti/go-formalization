package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Size, DstSize int
	Source        []byte
	Dst           []uint32
	Want          string
}

func main() {
	cases := []fixture{
		{0, 1, []byte{1, 2, 3, 4}, []uint32{91, 92}, "31867|false|91,92"},
		{3, 1, []byte{1, 2, 3, 4}, []uint32{91}, "31867|false|91"},
		{3, 1, []byte{1, 2, 3}, []uint32{91}, "bounds|false|91"},
		{6, 1, []byte{1, 2, 3, 4, 5, 6}, []uint32{91, 92, 93, 94}, "31867|true|31867,92,93,94"},
		{5, 2, []byte{1, 2, 3, 4, 5, 6}, []uint32{91, 92, 93}, "31867|false|31867,32064,93"},
		{4, 0, []byte{1, 2, 3, 4}, []uint32{91}, "31867|true|91"},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-flate-hash-window-")
	must(err)
	defer os.Remove(truth.Name())
	encoder := json.NewEncoder(truth)
	for _, c := range cases {
		fmt.Fprint(&input, c.Size, " ", len(c.Source), " ", c.DstSize)
		for _, b := range c.Source {
			fmt.Fprint(&input, " ", b)
		}
		for _, w := range c.Dst {
			fmt.Fprint(&input, " ", w)
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_FLATE_HASH_WINDOW_ROWS", truth.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanHashWindow", "difftest/compress/flate/hash-window/testdata/window_test.go"))

	rows := check.LeanRows("godiff-flate-hash-window", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row    int
		output string
	}{
		{0, "bounds|false|91,92"},               // bound hash4 by length
		{2, "31867|false|91"},                   // ignore short capacity
		{1, "31867|false|31867"},                // count capacity as a visible window
		{3, "31867|true|91,92,93,94"},           // precheck all destination writes
		{3, "31867|false|31867,32064,32260,94"}, // write through destination capacity
		{4, "100185|false|100185,32064,93"},     // load the first word in big endian order
		{5, "31867|false|91"},                   // suppress first write bounds failure
	} {
		if rows[m.row] == m.output {
			log.Fatalf("row %d accepted mutant %q", m.row, m.output)
		}
	}
	log.Print("six independent pinned Go and Lean rows passed and rejected seven wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
