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
	Input, Message, Suffix []byte
	Want                   string
}

func wire(s string) string {
	fields := make([]string, len(s))
	for i, b := range []byte(s) {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func main() {
	makeFixture := func(input, message, suffix, quoted, clause string) fixture {
		return fixture{[]byte(input), []byte(message), []byte(suffix),
			wire("ParseAddr("+quoted+"): "+message+clause) + "|" +
				wire("netip.ParsePrefix("+quoted+"): "+message)}
	}
	cases := []fixture{
		makeFixture("", "", "", `""`, ""),
		makeFixture("192.0.2.1", "raw\x00message\xff", "", `"192.0.2.1"`, ""),
		makeFixture("x\xff�", "bad byte", "\xff�", `"x\xff�"`, ` (at "\xff�")`),
		makeFixture("\u00a0\u2028\U0010ffff😀", "mark", "\t", `"\u00a0\u2028\U0010ffff😀"`, ` (at "\t")`),
		makeFixture("\"\\\n", "escape", "\"\\", `"\"\\\n"`, ` (at "\"\\")`),
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-netip-error-render-")
	must(err)
	defer os.Remove(truth.Name())
	enc := json.NewEncoder(truth)
	for _, c := range cases {
		fmt.Fprint(&input, len(c.Input), " ", len(c.Message))
		for _, field := range [][]byte{c.Input, c.Message, c.Suffix} {
			for _, b := range field {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_NETIP_RENDER_ERRORS", truth.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanRenderErrors", "difftest/net/netip/render-errors/testdata/render_test.go"))
	rows := check.LeanRows("godiff-netip-render-errors", len(cases), &input)
	for i, got := range rows {
		if got != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, got, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row, col int
		wrong    string
	}{
		{0, 0, `ParseAddr(""):  (at "")`},
		{1, 0, `ParseAddr("192.0.2.1"): raw\x00message\xff`},
		{2, 0, `ParseAddr("x��"): bad byte (at "��")`},
		{2, 1, `netip.ParsePrefix("x\xff�"): bad byte (at "\xff�")`},
		{3, 0, `ParseAddr("\u00a0\u2028\u10ffff😀"): mark (at "\t")`},
		{3, 1, `netip.ParsePrefix("\u00a0\u2028\U0010ffff\U0001f600"): mark`},
		{4, 0, `ParseAddr("\"\\\n"): escape`}, // omit the quoted suffix
	} {
		if strings.Split(rows[m.row], "|")[m.col] == wire(m.wrong) {
			log.Fatalf("row %d column %d accepted mutant %q", m.row, m.col, m.wrong)
		}
	}
	log.Print("five independent pinned Go and Lean cases passed and rejected seven wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
