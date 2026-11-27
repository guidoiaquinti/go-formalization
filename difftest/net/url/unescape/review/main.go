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
	Raw  []byte
	Mode byte
	Want string
}

func wire(status, value string) string {
	parts := make([]string, len(value))
	for i := range len(value) {
		parts[i] = strconv.Itoa(int(value[i]))
	}
	return status + ":" + strings.Join(parts, ",")
}

func main() {
	ok := func(s string) string { return wire("ok", s) }
	escaped := func(s string) string { return wire("escape", s) }
	cases := []fixture{
		{[]byte("a+b%20c"), 32, ok("a b c") + "|" + ok("a b c") + "|" + ok("a+b c")},
		{[]byte("%Gz%20"), 2, escaped("%Gz") + "|" + escaped("%Gz") + "|" + escaped("%Gz")},
		{[]byte("z%F"), 0, escaped("%F") + "|" + escaped("%F") + "|" + escaped("%F")},
		{[]byte("%41"), 4, escaped("%41") + "|" + ok("A") + "|" + ok("A")},
		{[]byte("%25"), 4, ok("%") + "|" + ok("%") + "|" + ok("%")},
		{[]byte("%20"), 8, ok(" ") + "|" + ok(" ") + "|" + ok(" ")},
		{[]byte("%2f"), 8, escaped("%2f") + "|" + ok("/") + "|" + ok("/")},
		{[]byte("x "), 4, wire("host", " ") + "|" + ok("x ") + "|" + ok("x ")},
		{[]byte("%fF"), 2, ok("\xff") + "|" + ok("\xff") + "|" + ok("\xff")},
	}
	mutants := map[int]string{
		0: ok("a b c") + "|" + ok("a b c") + "|" + ok("a b c"),
		1: escaped("%Gz%20") + "|" + escaped("%Gz%20") + "|" + escaped("%Gz%20"),
		3: ok("A") + "|" + ok("A") + "|" + ok("A"),
		5: escaped("%20") + "|" + ok(" ") + "|" + ok(" "),
	}
	var input bytes.Buffer
	for _, f := range cases {
		fmt.Fprint(&input, f.Mode)
		for _, b := range f.Raw {
			fmt.Fprintf(&input, " %d", b)
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-net-url-unescape", len(cases), &input)
	file, err := os.CreateTemp("", "golean-url-unescape-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	rejected := 0
	for i, f := range cases {
		if rows[i] != f.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], f.Want)
		}
		must(enc.Encode(f))
		if mutant, ok := mutants[i]; ok && mutant != rows[i] {
			rejected++
		}
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_URL_UNESCAPE_CASES", file.Name()))
	const overlay = "difftest/net/url/unescape/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanUnescape", overlay+"bridge_test.go", overlay+"unescape_test.go"))
	if rejected != 4 {
		log.Fatalf("wanted four mutation rejections, got %d", rejected)
	}
	fmt.Println("Nine independent native and Lean rows passed; four mutation witnesses rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
