package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Family        int
	Hi, Lo        uint64
	Zone, Backing []byte
	Length        int
	Want          string
}

func wire(text string) string {
	fields := make([]string, len(text))
	for i := range fields {
		fields[i] = strconv.Itoa(int(text[i]))
	}
	return strings.Join(fields, ",")
}

func observed(reuse, nilSlice bool, output, old string) string {
	return fmt.Sprintf("%t:%t:%s:%s", reuse, nilSlice, wire(output), wire(old))
}

func result(compact, expanded, marshal, appendOutput, privateOutput string) string {
	return strings.Join([]string{wire(compact), wire(expanded), wire(marshal), appendOutput, appendOutput, privateOutput}, "|")
}

func main() {
	const v4 = 0xffffff000102
	const mapped = 0xffff01020304
	v4out := observed(false, false, "!255.0.1.2", "!2")
	v6out := observed(false, false, "!2001::1:0:0:2:3", "!2001Z")
	mappedOut := observed(false, false, "!::ffff:1.2.3.4%z\xff", "!ZZZZZZ")
	invalidAppend := observed(true, true, "", "")
	cases := []fixture{
		{0, 0, 0, nil, nil, 0, result("invalid IP", "invalid IP", "", invalidAppend, "skip")},
		{4, 0, v4, nil, []byte("!Z"), 1, result("255.0.1.2", "255.0.1.2", "255.0.1.2", v4out, v4out)},
		{6, 0, mapped, []byte("z\xff"), []byte("!ZZZZZZ"), 1, result("::ffff:1.2.3.4%z\xff",
			"0000:0000:0000:0000:0000:ffff:0102:0304%z\xff", "::ffff:1.2.3.4%z\xff", mappedOut, mappedOut)},
		{6, 0x2001000000000001, 0x0000000000020003, nil, []byte("!ZZZZZ"), 1, result("2001::1:0:0:2:3",
			"2001:0000:0000:0001:0000:0000:0002:0003", "2001::1:0:0:2:3", v6out, v6out)},
	}
	file, err := os.CreateTemp("", "review-netip-addr-text-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	var protocol bytes.Buffer
	for _, c := range cases {
		missing := 0
		if c.Backing == nil {
			missing = 1
		}
		fmt.Fprint(&protocol, c.Family, " ", c.Hi, " ", c.Lo, " ", len(c.Zone), " ", c.Length, " ", missing)
		for _, part := range [][]byte{c.Zone, c.Backing} {
			for _, b := range part {
				fmt.Fprint(&protocol, " ", b)
			}
		}
		fmt.Fprintln(&protocol)
		must(encoder.Encode(c))
	}
	must(file.Close())
	rows := check.LeanRows("godiff-netip-addr-text", len(cases), &protocol)
	for i, row := range rows {
		if row != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_NETIP_ADDR_TEXT_CASES", file.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanAddrText", "difftest/net/netip/addr-text/testdata/addr_text_test.go"))
	reject := func(row, part int, mutated string) {
		if strings.Split(rows[row], "|")[part] == mutated {
			panic("mutation survived")
		}
	}
	reject(0, 2, wire("invalid IP"))
	reject(1, 3, observed(false, false, "!255.0.1.2", "!Z"))
	reject(2, 5, observed(false, false, "!::ffff:1.2.3.4%z\xff", "!::ffff"))
	reject(3, 0, wire("2001:0:0:1::2:3"))
	reject(3, 5, observed(false, false, "!2001::1:0:0:2:3", "!2001:"))
	reject(2, 1, wire("::ffff:1.2.3.4%z\xff"))
	fmt.Println("rejected invalid text leakage, atomic IPv4 digits, split mapped marker, last-run compaction, split double colon and compact expanded form")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
