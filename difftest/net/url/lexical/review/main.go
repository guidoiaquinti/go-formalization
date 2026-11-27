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
	Present bool
	Raw     []byte
	Want    string
}

func wire(raw string) string {
	fields := make([]string, len(raw))
	for i := range fields {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return fmt.Sprintf("%d:%s", len(raw), strings.Join(fields, ","))
}
func makeFixture(present bool, raw, host, port, name, suffix string, valid, missing, ctl bool) fixture {
	h, p, abs := "panic", "panic", "panic"
	if present {
		h, p, abs = wire(host), wire(port), strconv.FormatBool(raw != "")
	}
	scheme := wire(name) + "/" + wire(suffix)
	if missing {
		scheme = "missing"
	}
	return fixture{present, []byte(raw), strings.Join([]string{strconv.FormatBool(valid), wire(host) + "/" + wire(port), h, p, abs, scheme, strconv.FormatBool(ctl)}, "|")}
}
func main() {
	digits := "000184467440737095516160"
	cases := []fixture{
		makeFixture(true, "", "", "", "", "", true, false, false),
		makeFixture(true, ":", "", "", "", "", true, true, false),
		makeFixture(false, ":"+digits, "", digits, "", "", true, true, false),
		makeFixture(true, "[::1]:080", "::1", "080", "", "[::1]:080", false, false, false),
		makeFixture(true, "[::1]:bad", "[::1]:bad", "", "", "[::1]:bad", false, false, false),
		makeFixture(true, "[a]:1:2", "[a]:1", "2", "", "[a]:1:2", false, false, false),
		makeFixture(true, "A+1:\xff", "A+1:\xff", "", "A+1", "\xff", false, false, false),
		makeFixture(true, "0A:name", "0A:name", "", "", "0A:name", false, false, false),
		makeFixture(true, "\x00", "\x00", "", "", "\x00", false, false, true),
		makeFixture(true, "\x7f", "\x7f", "", "", "\x7f", false, false, true),
		makeFixture(true, "\x80\xff", "\x80\xff", "", "", "\x80\xff", false, false, false),
	}
	mutants := []struct {
		row, field  int
		alternative string
	}{
		{1, 0, "false"}, {1, 5, "0:/0:"}, {2, 0, "false"}, {2, 2, "0:"},
		{3, 3, wire("80")}, {4, 2, wire("::1")}, {5, 1, wire("[a]") + "/" + wire("1:2")},
		{6, 5, wire("a+1") + "/" + wire("\xff")}, {7, 5, wire("0A") + "/" + wire("name")},
		{8, 4, "false"}, {9, 6, "false"}, {10, 6, "true"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		tag := 0
		if c.Present {
			tag = 1
		}
		fmt.Fprint(&input, tag)
		for _, v := range c.Raw {
			fmt.Fprintf(&input, " %d", v)
		}
		fmt.Fprintln(&input)
	}
	rejected := 0
	for i, got := range check.LeanRows("godiff-net-url-lexical", len(cases), &input) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, cases[i].Want)
		}
		fields := strings.Split(got, "|")
		for _, m := range mutants {
			if m.row == i && fields[m.field] != m.alternative {
				rejected++
			}
		}
	}
	if rejected != len(mutants) {
		log.Fatalf("rejected %d of %d variants", rejected, len(mutants))
	}
	f, err := os.CreateTemp("", "golean-url-lexical-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_URL_LEXICAL_CASES", f.Name()))
	const base = "difftest/net/url/lexical/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanURLLexical", base+"bridge_test.go", base+"lexical_test.go"))
	fmt.Printf("URL lexical review: %d paired fixed observations, %d altered outcomes rejected\n", len(cases), rejected)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
