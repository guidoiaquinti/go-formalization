package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Word         uint64
	Base, Length int
	Cells        []byte
	Want         string
}

func wire(text string) string {
	fields := make([]string, len(text))
	for i, b := range []byte(text) {
		fields[i] = strconv.Itoa(int(b))
	}
	return "slice:" + strings.Join(fields, ",")
}

func whole(group ...string) string {
	result := strings.Join(group, "|")
	return result + "|" + result
}

func appendCase(out, old string, reused bool) string {
	return fmt.Sprintf("false/%s/%s/%t", wire(out), old, reused)
}

func main() {
	wideItoa := "4294967295"
	if strconv.IntSize == 32 {
		wideItoa = "-1"
	}
	unsignedBinary := strings.Repeat("1", 60) + "0001"
	minimumHex := "8000000000000000"
	minimumDecimal := "-9223372036854775808"
	if strconv.IntSize == 32 {
		minimumDecimal = "0"
	}
	cases := []fixture{
		{math.MaxUint64, 1, 1, []byte("!Z"), whole("panic", "panic", wire("-1"),
			"true/nil/"+wire("!Z")+"/false", "true/nil/"+wire("!-")+"/false")},
		{math.MaxUint64 - 14, 2, 1, []byte("!Z"), whole(wire(unsignedBinary), wire("-1111"), wire("-15"),
			appendCase("!"+unsignedBinary, wire("!Z"), false), appendCase("!-1111", wire("!-"), false))},
		{35, 36, 1, []byte("!ZZ"), whole(wire("z"), wire("z"), wire("35"),
			appendCase("!z", wire("!zZ"), true), appendCase("!z", wire("!zZ"), true))},
		{1 << 63, 16, 0, nil, whole(wire(minimumHex), wire("-"+minimumHex), wire(minimumDecimal),
			appendCase(minimumHex, "nil", false), appendCase("-"+minimumHex, "nil", false))},
		{math.MaxUint32, 16, 0, []byte{}, whole(wire("ffffffff"), wire("ffffffff"), wire(wideItoa),
			appendCase("ffffffff", wire(""), false), appendCase("ffffffff", wire(""), false))},
	}

	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-strconv-integer-format-")
	must(err)
	defer os.Remove(truth.Name())
	enc := json.NewEncoder(truth)
	for _, c := range cases {
		present := 1
		if c.Cells == nil {
			present = 0
		}
		fmt.Fprint(&input, strconv.IntSize, " ", c.Word, " ", c.Base, " ", c.Length, " ", present)
		for _, b := range c.Cells {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_STRCONV_FORMAT_CASES", truth.Name()))
	must(check.TestInPackage("strconv", "TestGoLeanIntegerFormat",
		"difftest/strconv/integer-format/testdata/format_test.go"))
	rows := check.LeanRows("godiff-strconv-integer-format", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d = %q, want %q", i, row, cases[i].Want)
		}
	}
	for i := range 2 {
		reject := func(row, col int, mutant string) {
			if strings.Split(rows[row], "|")[i*5+col] == mutant {
				log.Fatalf("row %d column %d accepted %q", row, col, mutant)
			}
		}
		reject(2, 0, wire("Z"))
		reject(3, 1, wire(minimumHex))
		reject(2, 2, wire("23")) // render using hexadecimal
		reject(2, 3, appendCase("!z", wire("!ZZ"), false))
		reject(0, 4, "true/nil/"+wire("!Z")+"/false")         // validate base before sign
		reject(1, 4, appendCase("!-1111", wire("!Z"), false)) // append sign and digits at once
	}
	log.Print("five pinned Go and Lean rows passed and rejected twelve wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
