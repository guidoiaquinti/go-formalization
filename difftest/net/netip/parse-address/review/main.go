package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func octets(s string) string {
	parts := make([]string, len(s))
	for i, b := range []byte(s) {
		parts[i] = fmt.Sprint(b)
	}
	return strings.Join(parts, ",")
}

func main() {
	zero := "v6:0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0:"
	emptyV6 := "error:each colon-separated field must have at least one digit::"
	emptyDispatch := "error:unable to parse IP::"
	zoneFirst := "%lab"
	zoneV6 := "error:each colon-separated field must have at least one digit:" + octets(zoneFirst) + ":"
	zoneDispatch := "error:missing IPv6 address:" + octets(zoneFirst) + ":"
	mapped := "v6:0,0,0,0,0,0,0,0,0,0,255,255,192,0,2,42:108,97,98"
	cases := []struct {
		input string
		want  [4]string
	}{
		{"::", [4]string{zero, zero, zero, zero}},
		{zoneFirst, [4]string{zoneV6, zoneDispatch, "panic:" + zoneDispatch, zoneDispatch}},
		{"", [4]string{emptyV6, emptyDispatch, "panic:" + emptyDispatch, "clear"}},
		{"::ffff:192.0.2.42%lab", [4]string{mapped, mapped, mapped, mapped}},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "netip-address-review-*")
	must(err)
	defer os.Remove(truth.Name())
	encoder := json.NewEncoder(truth)
	for _, c := range cases {
		fmt.Fprint(&input, "0")
		for _, b := range []byte(c.input) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(struct {
			Input []byte
			Want  string
		}{[]byte(c.input), strings.Join(c.want[:], "|")}))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_NETIP_PARSE_ADDRESS", truth.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanParseAddress",
		"difftest/net/netip/parse-address/testdata/address_test.go"))
	rows := check.LeanRows("godiff-netip-parse-address", len(cases), &input)
	for i, got := range rows {
		if got != strings.Join(cases[i].want[:], "|") {
			log.Fatalf("row %d: got %q, want %q", i, got, cases[i].want)
		}
	}
	mutations := []struct {
		row, column int
		value       string
	}{
		{0, 0, "v6:0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1:"},
		{1, 1, zoneV6},
		{2, 2, emptyDispatch},
		{2, 3, emptyDispatch},
		{3, 1, "v4:192,0,2,42"},
	}
	for _, m := range mutations {
		if strings.Split(rows[m.row], "|")[m.column] == m.value {
			log.Fatalf("row %d column %d accepted %q", m.row, m.column, m.value)
		}
	}
	fmt.Printf("%d pinned Go and Lean rows passed and rejected %d wrong outcomes\n", len(cases), len(mutations))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
