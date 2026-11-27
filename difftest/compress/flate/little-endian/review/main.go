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
	Tag, Length  int
	Index, Value uint64
	Cells        []byte
	Want         string
}

func byteWire(raw []byte) string {
	parts := make([]string, len(raw))
	for i, b := range raw {
		parts[i] = strconv.Itoa(int(b))
	}
	return strings.Join(parts, ",")
}

func main() {
	cells := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	short := append(make([]byte, 124), cells[:8]...)
	long := append(make([]byte, 252), cells[:8]...)
	const value = 0x1122334455667788
	shortStored := append([]byte{136, 119, 102, 85, 68, 51, 34, 17}, short[8:]...)
	longStored := append([]byte{136, 119, 102, 85, 68, 51, 34, 17}, long[8:]...)
	digits := "|67305985|578437695752307201|"
	cases := []fixture{
		{1, 0, 0, value, cells, "bounds" + digits + "false|" + byteWire(cells)},
		{1, 7, 0, value, cells, "1" + digits + "false|" + byteWire(cells)},
		{1, 8, 0, value, cells, "1" + digits + "true|136,119,102,85,68,51,34,17,9,10"},
		{0, 125, 124, value, short, "1|bounds|bounds|true|" + byteWire(shortStored)},
		{1, 125, 124, value, short, "1" + digits + "true|" + byteWire(shortStored)},
		{1, 253, 252, value, long, "1|bounds|bounds|true|" + byteWire(longStored)},
		{3, 253, 252, value, long, "1" + digits + "true|" + byteWire(longStored)},
		{0, 8, 255, value, make([]byte, 8), "bounds|bounds|bounds|true|136,119,102,85,68,51,34,17"},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-flate-little-endian-")
	must(err)
	defer os.Remove(truth.Name())
	enc := json.NewEncoder(truth)
	for _, c := range cases {
		width, signed, start := 8, 0, int64(c.Index)
		switch c.Tag {
		case 0:
			signed, start = 1, int64(int8(c.Index))
		case 3:
			width = 16
		}
		fmt.Fprint(&input, width, " ", signed, " ", start, " ", c.Length, " ", c.Value)
		for _, b := range c.Cells {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_FLATE_ENDIAN_ROWS", truth.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanFlateEndian",
		"difftest/compress/flate/little-endian/testdata/little_endian_test.go"))
	rows := check.LeanRows("godiff-flate-little-endian", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	reject := func(row, field int, mutant string) {
		if strings.Split(rows[row], "|")[field] == mutant {
			log.Fatalf("row %d column %d accepted %q", row, field, mutant)
		}
	}
	reject(0, 0, "1")      // index load using capacity
	reject(0, 1, "bounds") // slice load using length
	reject(0, 2, "bounds")
	reject(1, 4, "136,119,102,85,68,51,34,8,9,10") // partial store then panic
	reject(2, 4, "136,119,102,85,68,51,34,17,0,0") // zero the spare backing
	reject(3, 1, "67305985")                       // widen signed endpoint
	reject(5, 2, "578437695752307201")             // widen unsigned endpoint
	reject(7, 0, "0")                              // coerce negative index to zero
	log.Print("eight independent pinned Go and Lean cases passed and rejected eight wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
