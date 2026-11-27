package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	text   string
	lo, hi int
	want   string
}

func main() {
	v4 := "192.0.2.1"
	overflow := "192.0.2.256"
	subrange := "x1.2.qy"
	bracketed := "[2001:db8::1]:443"
	short := "[]:80"
	badBracket := "[x:443"
	fault := func(reason, input, suffix string) string {
		return reason + ":" + wire(input) + ":" + wire(suffix)
	}
	unexpected := func(input, suffix string) string {
		return fault("unexpected character", input, suffix)
	}
	row := func(cells, part, whole, split string) string {
		return strings.Join([]string{cells, part, whole, split}, "|")
	}
	tooShort := func(input string) string { return fault("IPv4 address too short", input, "") }
	cases := []fixture{
		{v4, 0, len(v4), row("192,0,2,1", "ok", "ok:192,0,2,1", "not an ip:port")},
		{overflow, 0, len(overflow), row("192,0,2,93", fault("IPv4 field has value >255", overflow, ""),
			fault("IPv4 field has value >255", overflow, ""), "not an ip:port")},
		{subrange, 1, 6, row("1,2,92,93", unexpected(subrange, "q"), unexpected(subrange, subrange), "not an ip:port")},
		{bracketed, 0, 0, row("90,91,92,93", tooShort(bracketed), unexpected(bracketed, bracketed),
			wire("2001:db8::1")+":"+wire("443")+":true")},
		{short, 0, 0, row("90,91,92,93", tooShort(short), unexpected(short, short), ":"+wire("80")+":true")},
		{badBracket, 0, 0, row("90,91,92,93", tooShort(badBracket), unexpected(badBracket, badBracket), "missing ]")},
		{":", 0, 0, row("90,91,92,93", tooShort(":"), unexpected(":", ":"), "no IP")},
		{"host:", 0, 0, row("90,91,92,93", tooShort("host:"), unexpected("host:", "host:"), "no port")},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-netip-parse-fields-")
	must(err)
	defer os.Remove(truth.Name())
	for _, c := range cases {
		fmt.Fprint(&input, c.lo, " ", c.hi, " 90 91 92 93")
		for _, b := range []byte(c.text) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		_, err := fmt.Fprintf(truth, "%d %d 90 91 92 93 %q\t%s\n", c.lo, c.hi, c.text, c.want)
		must(err)
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_NETIP_PARSE_FIELDS", truth.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanParseFields", "difftest/net/netip/parse-fields/testdata/parse_fields_test.go"))
	rows := check.LeanRows("godiff-netip-parse-fields", len(cases), &input)
	for i, got := range rows {
		if got != cases[i].want {
			log.Fatalf("row %d got %q, want %q", i, got, cases[i].want)
		}
	}
	for _, m := range []struct {
		row, col int
		value    string
	}{
		{1, 0, "90,91,92,93"},              // discard completed field writes on error
		{2, 1, unexpected(subrange, "qy")}, // take diagnostic suffix outside the selected range
		{2, 1, unexpected("1.2.q", "q")},   // report subrange as full input
		{1, 2, "ok:192,0,2,0"},             // turn an out-of-range octet into a success
		{3, 3, "missing ]"},                // split at the first colon
		{4, 3, "no IP"},                    // reject bracketed empty host in the shape-only splitter
		{6, 3, "no port"},                  // check port before host
	} {
		if strings.Split(rows[m.row], "|")[m.col] == m.value {
			log.Fatalf("row %d column %d accepted mutant %q", m.row, m.col, m.value)
		}
	}
	log.Print("eight pinned Go and Lean fixtures passed and rejected seven wrong observations")
}

func wire(raw string) string {
	fields := make([]string, len(raw))
	for i, b := range []byte(raw) {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
